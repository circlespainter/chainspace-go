package conns // import "chainspace.io/chainspace-go/internal/conns"

import (
	"sync"
	"time"

	"chainspace.io/chainspace-go/internal/crypto/signature"
	"chainspace.io/chainspace-go/internal/log"
	"chainspace.io/chainspace-go/internal/log/fld"
	"chainspace.io/chainspace-go/network"
	"chainspace.io/chainspace-go/service"
)

type cache struct {
	conns      map[uint64]*MuConn
	cmu        []sync.Mutex
	mu         sync.Mutex
	key        signature.KeyPair
	maxPayload int
	selfID     uint64
	top        network.NetTopology
	connection service.CONNECTION

	pendingAcks   map[AckID]PendingAck
	pendingAcksMu sync.Mutex
}

type AckID struct {
	NodeID    uint64
	RequestID uint64
}

type MuConn struct {
	mu   sync.Mutex
	conn *network.Conn
	die  chan bool
}

type PendingAck struct {
	nodeID  uint64
	msg     *service.Message
	sentAt  time.Time
	timeout time.Duration
	cb      func(uint64, *service.Message)
}

type Cache interface {
	Close()
	WriteRequest(
		nodeID uint64, msg *service.Message, timeout time.Duration, ack bool,
		cb func(uint64, *service.Message)) (uint64, error)
}

func (c *cache) addPendingAck(nodeID uint64, msg *service.Message, timeout time.Duration, id uint64, cb func(uint64, *service.Message)) {
	ack := PendingAck{
		sentAt:  time.Now(),
		nodeID:  nodeID,
		msg:     msg,
		timeout: timeout,
		cb:      cb,
	}
	c.pendingAcksMu.Lock()
	c.pendingAcks[AckID{nodeID, id}] = ack
	c.pendingAcksMu.Unlock()
}

func (c *cache) dial(nodeID uint64) (*MuConn, error) {
	// conn exist
	c.mu.Lock()
	cc, ok := c.conns[nodeID]
	if ok {
		c.mu.Unlock()
		return cc, nil
	}

	defer c.mu.Unlock()
	// log.Error("NEED TO DIAL", fld.NodeID(nodeID))
	// need to dial
	conn, err := c.top.Dial(nodeID, 5*time.Hour)
	if err != nil {
		return nil, err
	}
	err = c.sendHello(nodeID, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	cc = &MuConn{conn: conn, die: make(chan bool)}
	go c.readAckMessage(nodeID, cc.conn, cc.die)
	c.conns[nodeID] = cc
	return cc, nil
}

func (c *cache) processAckMessage(nodeID uint64, msg *service.Message) {
	c.pendingAcksMu.Lock()
	defer c.pendingAcksMu.Unlock()
	if m, ok := c.pendingAcks[AckID{nodeID, msg.ID}]; ok {
		if m.cb != nil {
			m.cb(m.nodeID, msg)
		}
		delete(c.pendingAcks, AckID{nodeID, msg.ID})
	} else {
		log.Error("unknown lastID", log.Uint64("lastid", msg.ID))
		if log.AtDebug() {
			log.Debug("unknown lastID", log.Uint64("lastid", msg.ID))
		}
	}
}

func (c *cache) readAckMessage(nodeID uint64, conn *network.Conn, die chan bool) {
	for {
		select {
		case _ = <-die:
			log.Error("KILLING READ ACK", fld.PeerID(nodeID))
			return
		default:
			msg, err := conn.ReadMessage(int(c.maxPayload), 1*time.Second)
			// if we can read some message, try to process it.
			if err == nil {
				go c.processAckMessage(nodeID, msg)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (c *cache) release(nodeID uint64) {
	c.mu.Lock()
	cc, ok := c.conns[nodeID]
	c.mu.Unlock()
	if ok {
		cc.die <- true
		cc.conn.Close()
		c.mu.Lock()
		delete(c.conns, nodeID)
		c.mu.Unlock()
	}
}

func (c *cache) retryRequests() {
	for {
		redolist := []PendingAck{}
		time.Sleep(3 * time.Second)
		c.pendingAcksMu.Lock()
		for k, v := range c.pendingAcks {
			if time.Since(v.sentAt) >= 3*time.Second {
				redolist = append(redolist, v)
				delete(c.pendingAcks, k)
			}
		}
		c.pendingAcksMu.Unlock()
		for _, v := range redolist {
			c.WriteRequest(v.nodeID, v.msg, v.timeout, true, v.cb)
		}
	}
}

func (c *cache) sendHello(nodeID uint64, conn *network.Conn) error {
	hellomsg, err := service.SignHello(
		c.selfID, nodeID, c.key, c.connection)
	if err != nil {
		return err
	}
	return conn.WritePayload(hellomsg, c.maxPayload, time.Second)
}

func (c *cache) Close() {}

func (c *cache) WriteRequest(
	nodeID uint64, msg *service.Message, timeout time.Duration, ack bool, cb func(uint64, *service.Message)) (uint64, error) {
	c.cmu[nodeID-1].Lock()
	mc, err := c.dial(nodeID)
	if err != nil {
		c.release(nodeID)
		c.cmu[nodeID-1].Unlock()
		time.Sleep(5 * time.Millisecond)
		return c.WriteRequest(nodeID, msg, timeout, ack, cb)
	}
	id, err := mc.conn.WriteRequest(msg, c.maxPayload, timeout)
	if err != nil {
		c.release(nodeID)
		c.cmu[nodeID-1].Unlock()
		return c.WriteRequest(nodeID, msg, timeout, ack, cb)
	}
	c.cmu[nodeID-1].Unlock()
	if ack {
		c.addPendingAck(nodeID, msg, timeout, id, cb)
	}
	return id, nil
}

func NewCache(nodeID uint64, top network.NetTopology, maxPayload int, key signature.KeyPair, connection service.CONNECTION) *cache {
	c := &cache{
		conns:       map[uint64]*MuConn{},
		cmu:         make([]sync.Mutex, top.TotalNodes()),
		maxPayload:  maxPayload,
		selfID:      nodeID,
		top:         top,
		key:         key,
		pendingAcks: map[AckID]PendingAck{},
		connection:  connection,
	}
	go c.retryRequests()

	return c
}
