package node // import "chainspace.io/prototype/node"

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base32"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"chainspace.io/prototype/config"
	"chainspace.io/prototype/crypto/signature"
	"chainspace.io/prototype/freeport"
	"chainspace.io/prototype/log"
	"chainspace.io/prototype/network"
	"chainspace.io/prototype/service"
	"chainspace.io/prototype/service/broadcast"
	"chainspace.io/prototype/service/transactor"
	"github.com/gogo/protobuf/proto"
	"github.com/lucas-clemente/quic-go"
	"github.com/tav/golly/process"
)

const (
	dirPerms = 0700
)

// Error values.
var (
	errDirectoryMissing = errors.New("node: config is missing a value for Directory")
	errInvalidSignature = errors.New("node: invalid signature on hello")
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

type usedNonce struct {
	nonce []byte
	ts    time.Time
}

// Config represents the configuration of an individual node within a Chainspace network.
type Config struct {
	Directory   string
	Keys        *config.Keys
	Network     *config.Network
	NetworkName string
	NodeID      uint64
	Node        *config.Node
}

// Server represents a running Chainspace node.
type Server struct {
	broadcaster     *broadcast.Service
	cancel          context.CancelFunc
	ctx             context.Context
	id              uint64
	initialBackoff  time.Duration
	key             signature.KeyPair
	keys            map[uint64]signature.PublicKey
	readTimeout     time.Duration
	maxBackoff      time.Duration
	mu              sync.RWMutex // protects nonceMap
	nonceExpiration time.Duration
	nonceMap        map[uint64][]usedNonce
	payloadLimit    int
	sharder         service.Handler
	top             *network.Topology
	transactor      service.Handler
	writeTimeout    time.Duration
}

func (s *Server) handleConnection(conn quic.Session) {
	for {
		stream, err := conn.AcceptStream()
		if err != nil {
			log.Errorf("Unable to open new stream from session: %s", err)
			return
		}
		go s.handleStream(stream)
	}
}

func (s *Server) handleStream(stream quic.Stream) {
	defer stream.Close()
	c := network.NewConn(stream)
	hello, err := c.ReadHello(s.payloadLimit, s.readTimeout)
	if err != nil {
		log.Errorf("Unable to read hello message from stream: %s", err)
		return
	}
	var (
		svc    service.Handler
		peerID uint64
	)
	switch hello.Type {
	case service.CONNECTION_BROADCAST:
		svc = s.broadcaster
		peerID, err = s.verifyPeerID(hello)
		if err != nil {
			log.Errorf("Unable to verify peer ID from the hello message: %s", err)
			return
		}
	case service.CONNECTION_TRANSACTOR:
		svc = s.transactor
		log.Infof("new transactor hello nesssage")
	default:
		log.Errorf("Unknown connection type: %#v", hello.Type)
		return
	}
	ctx := stream.Context()
	for {
		msg, err := c.ReadMessage(s.payloadLimit, s.readTimeout)
		if err != nil {
			log.Errorf("Could not decode message from an incoming stream: %s", err)
			return
		}
		resp, err := svc.Handle(ctx, peerID, msg)
		if err != nil {
			log.Errorf("Received error response from the %s service: %s", svc.Name(), err)
			return
		}
		if resp != nil {
			if err = c.WritePayload(resp, s.payloadLimit, s.writeTimeout); err != nil {
				log.Errorf("Unable to write response to peer %d: %s", peerID, err)
				return
			}
		}
	}

}

func (s *Server) listen(l quic.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			s.cancel()
			log.Fatalf("node: could not accept new connections: %s", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) maintainBroadcast(peerID uint64) {
	var (
		blocks []*broadcast.SignedBlock
		round  uint64
	)
	backoff := s.initialBackoff
	msg := &service.Message{Opcode: uint32(broadcast.OP_LISTING)}
	listing := &broadcast.Listing{}
	retry := false
	for {
		if retry {
			backoff *= 2
			if backoff > s.maxBackoff {
				backoff = s.maxBackoff
			}
			time.Sleep(backoff)
			retry = false
		}
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		conn, err := s.top.Dial(s.ctx, peerID)
		if err == nil {
			backoff = s.initialBackoff
		} else {
			log.Errorf("Couldn't dial node %d: %s", peerID, err)
			retry = true
			continue
		}
		hello, err := service.BroadcastHello(s.id, peerID, s.key)
		if err != nil {
			log.Errorf("Couldn't create Hello payload for broadcast: %s", err)
			retry = true
			continue
		}
		if err = conn.WritePayload(hello, s.payloadLimit, s.writeTimeout); err != nil {
			log.Errorf("Couldn't send Hello to node %d: %s", peerID, err)
			retry = true
			continue
		}
		for {
			if len(blocks) == 0 {
				blocks, round = s.broadcaster.GetUnsent(peerID)
				if blocks == nil {
					return
				}
			}
			listing.Blocks = blocks
			msg.Payload, err = proto.Marshal(listing)
			if err != nil {
				log.Fatalf("Could not encode listing for broadcast: %s", err)
			}
			if err = conn.WritePayload(msg, s.payloadLimit, s.writeTimeout); err != nil {
				log.Errorf("Could not write listing to node %d for broadcast: %s", peerID, err)
				retry = true
				break
			}
			s.broadcaster.MarkSent(peerID, round)
			blocks = nil
		}
	}
}

func (s *Server) pruneNonceMap() {
	exp := s.nonceExpiration
	tick := time.NewTicker(exp)
	var xs []usedNonce
	for {
		select {
		case <-tick.C:
			now := time.Now()
			nmap := map[uint64][]usedNonce{}
			s.mu.Lock()
			for nodeID, nonces := range s.nonceMap {
				xs = nil
				for _, used := range nonces {
					diff := now.Sub(used.ts)
					if diff < 0 {
						diff = -diff
					}
					if diff < exp {
						xs = append(xs, used)
					}
				}
				if xs != nil {
					nmap[nodeID] = xs
				}
			}
			s.nonceMap = nmap
			s.mu.Unlock()
		case <-s.ctx.Done():
			tick.Stop()
			return
		}
	}
}

func (s *Server) verifyPeerID(hello *service.Hello) (uint64, error) {
	if len(hello.Payload) == 0 {
		return 0, nil
	}
	info := &service.HelloInfo{}
	if err := proto.Unmarshal(hello.Payload, info); err != nil {
		return 0, err
	}
	if info.Server != s.id {
		return 0, fmt.Errorf("node: mismatched server field in hello payload: expected %d, got %d", s.id, info.Server)
	}
	key, exists := s.keys[info.Client]
	if !exists {
		return 0, fmt.Errorf("node: could not find public signing key for node %d", info.Client)
	}
	if !key.Verify(hello.Payload, hello.Signature) {
		return 0, fmt.Errorf("node: could not validate signature claiming to be from node %d", info.Client)
	}
	diff := time.Now().Sub(info.Timestamp)
	if diff < 0 {
		diff = -diff
	}
	if diff > s.nonceExpiration {
		return 0, fmt.Errorf("node: timestamp in client hello is outside of the max clock skew range: %d", diff)
	}
	s.mu.RLock()
	nonces, exists := s.nonceMap[info.Client]
	if exists {
		for _, used := range nonces {
			if bytes.Equal(used.nonce, info.Nonce) {
				s.mu.RUnlock()
				return 0, fmt.Errorf("node: received duplicate nonce from %d", info.Client)
			}
		}
	}
	s.mu.RUnlock()
	return info.Client, nil
}

// AddTransaction adds the provided raw transaction data to the node's
// blockchain.
// func (s *Server) AddTransaction(txdata *TransactionData) {
// 	s.transactor.AddTransaction(txdata)
// }

// AddTransactionData is a temporary method for testing. It adds the given
// transaction to the node's current block.
func (s *Server) AddTransactionData(txdata *broadcast.TransactionData) {
	s.broadcaster.AddTransaction(txdata)
}

// Shutdown closes all underlying resources associated with the node.
func (s *Server) Shutdown() {
	s.cancel()
}

// Run initialises a node with the given config.
func Run(cfg *Config) (*Server, error) {
	var err error

	// Initialise the runtime directory for the node.
	dir := cfg.Directory
	if dir == "" {
		return nil, errDirectoryMissing
	}
	if !filepath.IsAbs(dir) {
		if dir, err = filepath.Abs(dir); err != nil {
			return nil, err
		}
	}
	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		log.Infof("Creating %s", dir)
		if err = os.MkdirAll(dir, dirPerms); err != nil {
			return nil, err
		}
	}

	if err := process.Init(dir, "chainspace"); err != nil {
		return nil, err
	}

	if err := log.InitFileLogger(filepath.Join(dir, "server.log"), log.DebugLevel); err != nil {
		log.Fatal(err)
	}

	// Initialise the topology.
	top, err := network.New(cfg.NetworkName, cfg.Network)
	if err != nil {
		return nil, err
	}

	// Bootstrap using a contacts JSON file.
	if cfg.Node.Bootstrap.File != "" {
		if err = top.BootstrapFile(cfg.Node.Bootstrap.File); err != nil {
			return nil, err
		}
	}

	// Bootstrap using mDNS.
	if cfg.Node.Bootstrap.MDNS {
		top.BootstrapMDNS()
	}

	// Bootstrap using a static map of addresses.
	if cfg.Node.Bootstrap.Static != nil {
		if err = top.BootstrapStatic(cfg.Node.Bootstrap.Static); err != nil {
			return nil, err
		}
	}

	// Bootstrap using a URL endpoint.
	if cfg.Node.Bootstrap.URL != "" {
		if err = top.BootstrapURL(cfg.Node.Bootstrap.URL); err != nil {
			return nil, err
		}
	}

	// Get a port to listen on.
	port, err := freeport.UDP("")
	if err != nil {
		return nil, err
	}

	// Initialise the TLS cert.
	cert, err := tls.X509KeyPair([]byte(cfg.Keys.TransportCert.Public), []byte(cfg.Keys.TransportCert.Private))
	if err != nil {
		return nil, fmt.Errorf("node: could not load the X.509 transport.cert from config: %s", err)
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Initialise signing keys.
	var key signature.KeyPair
	switch cfg.Keys.SigningKey.Type {
	case "ed25519":
		pubkey, err := b32.DecodeString(cfg.Keys.SigningKey.Public)
		if err != nil {
			return nil, fmt.Errorf("node: could not decode the base32-encoded public signing.key from config: %s", err)
		}
		privkey, err := b32.DecodeString(cfg.Keys.SigningKey.Private)
		if err != nil {
			return nil, fmt.Errorf("node: could not decode the base32-encoded private signing.key from config: %s", err)
		}
		key, err = signature.LoadKeyPair(signature.Ed25519, append(pubkey, privkey...))
		if err != nil {
			return nil, fmt.Errorf("node: unable to load the signing.key from config: %s", err)
		}
	default:
		return nil, fmt.Errorf("node: unknown type of signing.key found in config: %q", cfg.Keys.SigningKey.Type)
	}

	keys := top.SeedPublicKeys()
	shardID := top.ShardForNode(cfg.NodeID)
	nodes := top.NodesInShard(shardID)
	if len(nodes) == 0 {
		return nil, fmt.Errorf("node: got no peers for the node %d", cfg.NodeID)
	}

	idx := 0
	peers := make([]uint64, len(nodes)-1)
	for _, peer := range nodes {
		if peer != cfg.NodeID {
			peers[idx] = peer
			idx++
		}
	}

	blockLimit, err := cfg.Network.Consensus.BlockLimit.Int()
	if err != nil {
		return nil, err
	}

	maxPayload, err := cfg.Network.MaxPayload.Int()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	bcfg := &broadcast.Config{
		BlockLimit:    blockLimit,
		Directory:     dir,
		Key:           key,
		Keys:          keys,
		NodeID:        cfg.NodeID,
		Peers:         peers,
		RoundInterval: cfg.Network.Consensus.RoundInterval,
	}

	broadcaster, err := broadcast.New(ctx, bcfg, top)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("node: unable to instantiate the broadcast service: %s", err)
	}

	tcfg := &transactor.Config{
		Broadcaster: broadcaster,
		Checkers: []transactor.Checker{
			&transactor.DummyCheckerOK{}, &transactor.DummyCheckerKO{}},
		Directory:  dir,
		NodeID:     cfg.NodeID,
		Top:        top,
		SigningKey: cfg.Keys.SigningKey,
		MaxPayload: int(cfg.Network.MaxPayload),
	}
	txtor, err := transactor.New(tcfg)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("node: unable to instanciate the transactor service: %v", err)
	}

	node := &Server{
		broadcaster:     broadcaster,
		cancel:          cancel,
		ctx:             ctx,
		id:              cfg.NodeID,
		initialBackoff:  cfg.Node.Broadcast.InitialBackoff,
		key:             key,
		keys:            keys,
		maxBackoff:      cfg.Node.Broadcast.MaxBackoff,
		nonceExpiration: cfg.Network.Consensus.NonceExpiration,
		nonceMap:        map[uint64][]usedNonce{},
		payloadLimit:    maxPayload,
		readTimeout:     cfg.Node.Connections.ReadTimeout,
		top:             top,
		transactor:      txtor,
		writeTimeout:    cfg.Node.Connections.WriteTimeout,
	}

	// Maintain persistent connections to all other nodes in our shard.
	for _, peer := range peers {
		go node.maintainBroadcast(peer)
	}

	// Start listening on the given port.
	l, err := quic.ListenAddr(fmt.Sprintf(":%d", port), tlsConf, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("node: unable to listen on port %d: %s", port, err)
	}

	go node.pruneNonceMap()
	go node.listen(l)

	// Announce our location to the world.
	for _, meth := range cfg.Node.Announce {
		if meth == "mdns" {
			if err = announceMDNS(cfg.Network.ID, cfg.NodeID, port); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(meth, "https://") || strings.HasPrefix(meth, "http://") {
			// Announce to the given URL endpoint.
		} else {
			return nil, fmt.Errorf("node: announcement endpoint %q does not start with http:// or https://", meth)
		}
	}

	log.Infof("Node %d of the %s network is now running on port %d", cfg.NodeID, cfg.NetworkName, port)
	log.Infof("Runtime directory: %s", dir)
	return node, nil

}
