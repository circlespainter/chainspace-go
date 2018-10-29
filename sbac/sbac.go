package sbac // import "chainspace.io/prototype/sbac"

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"path"
	"sync"
	"time"

	"chainspace.io/prototype/broadcast"
	"chainspace.io/prototype/config"
	"chainspace.io/prototype/internal/combihash"
	"chainspace.io/prototype/internal/conns"
	"chainspace.io/prototype/internal/crypto/signature"
	"chainspace.io/prototype/internal/log"
	"chainspace.io/prototype/internal/log/fld"
	"chainspace.io/prototype/kv"
	"chainspace.io/prototype/network"
	"chainspace.io/prototype/pubsub"
	"chainspace.io/prototype/service"
	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"
)

const (
	badgerStorePath = "/sbac/"
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

type Config struct {
	Broadcaster *broadcast.Service
	Directory   string
	KVStore     *kv.Service
	NodeID      uint64
	Top         *network.Topology
	SigningKey  *config.Key
	Pubsub      *pubsub.Server
	ShardSize   uint64
	ShardCount  uint64
	MaxPayload  int
	Key         signature.KeyPair
}

type Service struct {
	broadcaster *broadcast.Service
	conns       *conns.Pool
	kvstore     *kv.Service
	nodeID      uint64
	pe          *pendingEvents
	privkey     signature.PrivateKey
	ps          *pubsub.Server
	store       *badger.DB
	shardCount  uint64
	shardID     uint64
	shardSize   uint64
	table       *StateTable
	top         *network.Topology
	txstates    map[string]*StateMachine
	txstatesmu  sync.Mutex
}

func (s *Service) handleDeliver(round uint64, blocks []*broadcast.SignedData) {
	for _, signed := range blocks {
		block, err := signed.Block()
		if err != nil {
			log.Fatal("Unable to decode delivered block", fld.Round(round), fld.Err(err))
		}
		it := block.Iter()
		for it.Valid() {
			it.Next()
			// TODO(): do stuff with the fee?
			tx := ConsensusTransaction{}
			err := proto.Unmarshal(it.TxData, &tx)
			if err != nil {
				log.Error("Unable to unmarshal transaction data", fld.Err(err))
				continue
			}
			e := NewConsensusEvent(&tx)
			s.checkEvent(e)
			s.pe.OnEvent(e)
		}
	}
}

// checkEvent check the event Op and create a new StateMachine if
// the OpCode is Consensus1 or ConsensusCommit
func (s *Service) checkEvent(e *ConsensusEvent) {
	if e.data.Op == ConsensusOp_Consensus1 {
		txbytes, _ := proto.Marshal(e.data.Tx)
		detail := DetailTx{ID: e.data.TxID, RawTx: txbytes, Tx: e.data.Tx}
		s.addStateMachine(&detail, StateWaitingForConsensus1)
	}
}

func (s *Service) Handle(peerID uint64, m *service.Message) (*service.Message, error) {
	ctx := context.TODO()
	switch Opcode(m.Opcode) {
	case Opcode_ADD_TRANSACTION:
		return s.addTransaction(ctx, m.Payload, m.ID)
	case Opcode_QUERY_OBJECT:
		return s.queryObject(ctx, m.Payload, m.ID)
	case Opcode_DELETE_OBJECT:
		return s.deleteObject(ctx, m.Payload, m.ID)
	case Opcode_CREATE_OBJECT:
		return s.createObject(ctx, m.Payload, m.ID)
	case Opcode_STATES:
		return s.handleStates(ctx, m.Payload, m.ID)
	case Opcode_SBAC:
		return s.handleSBAC(ctx, m.Payload, peerID, m.ID)
	case Opcode_CREATE_OBJECTS:
		return s.createObjects(ctx, m.Payload, m.ID)
	default:
		log.Error("sbac: unknown message opcode", log.Int32("opcode", m.Opcode), fld.PeerID(peerID), log.Int("len", len(m.Payload)))
		return nil, fmt.Errorf("sbac: unknown message opcode: %v", m.Opcode)
	}
}

func (s *Service) consumeEvents(e EventExt) bool {
	// check if statemachine is finished
	ok, err := TxnFinished(s.store, e.TxID())
	if err != nil {
		log.Error("error calling TxnFinished", fld.Err(err))
		// do nothing
		return true
	}
	if ok {
		log.Error("event for a finished transaction, skipping it",
			fld.TxID(ID(e.TxID())), log.Uint64("peer.id", e.PeerID()))
		// txn already finished. move on
		return true
	}

	sm, ok := s.getSateMachine(e.TxID())
	if ok {
		if log.AtDebug() {
			log.Debug("sending new event to statemachine", fld.TxID(ID(e.TxID())), fld.PeerID(e.PeerID()), fld.PeerShard(s.top.ShardForNode(e.PeerID())))
		}
		sm.OnEvent(e)
		return true
	}
	if log.AtDebug() {
		log.Debug("statemachine not ready", fld.TxID(ID(e.TxID())))
	}
	return false
}

func (s *Service) handleSBAC(
	ctx context.Context, payload []byte, peerID uint64, msgID uint64,
) (*service.Message, error) {
	req := &SBACMessage2{}
	err := proto.Unmarshal(payload, req)
	if err != nil {
		log.Error("sbac: sbac unmarshaling error", fld.Err(err))
		return nil, fmt.Errorf("sbac: sbac unmarshaling error: %v", err)
	}
	// if we received a COMMIT opcode, the statemachine may not exists
	// lets check and create it here.
	if req.Op == SBACOp_PhaseCommit {
		// txdetails := NewTxDetails(req.Tx.ID, []byte{}, req.Tx.Tx, req.Tx.Evidences)
		// _ = s.getOrCreateStateMachine(txdetails, StateWaitingForCommit)
	}
	// e := &Event{msg: req, peerID: peerID}
	req.PeerID = peerID
	e := NewSBACEvent(req)
	s.pe.OnEvent(e)
	res := SBACMessageAck{LastID: msgID}
	payloadres, err := proto.Marshal(&res)
	if err != nil {
		log.Error("unable to marshall handleSBAC response", fld.Err(err))
		return nil, err
	}
	return &service.Message{ID: msgID, Opcode: int32(Opcode_SBAC), Payload: payloadres}, nil
}

func (s *Service) verifySignatures(txID []byte, evidences map[uint64][]byte) bool {
	ok := true
	keys := s.top.SeedPublicKeys()
	if len(evidences) <= 0 {
		log.Errorf("missing signatures")
		return false
	}
	for nodeID, sig := range evidences {
		key := keys[nodeID]
		if nodeID == s.nodeID {
			if !key.Verify(txID, sig) {
				log.Errorf("invalid signature")
				if log.AtDebug() {
					log.Debug("invalid signature", fld.PeerID(nodeID))
				}
				ok = false
			}
		}
	}
	return ok
}

func (s *Service) addStateMachine(detail *DetailTx, initialState State) *StateMachine {
	s.txstatesmu.Lock()
	cfg := &StateMachineConfig{
		Consensus1Action: s.onConsensusEvent,
		Consensus2Action: s.onConsensusEvent,
		Consensus3Action: s.onConsensusEvent,
		Table:            s.table,
		Detail:           detail,
		InitialState:     initialState,
	}
	sm := NewStateMachine(cfg)
	s.txstates[string(detail.ID)] = sm
	s.txstatesmu.Unlock()
	return sm
}

func (s *Service) getSateMachine(txID []byte) (*StateMachine, bool) {
	s.txstatesmu.Lock()
	sm, ok := s.txstates[string(txID)]
	s.txstatesmu.Unlock()
	return sm, ok
}

func (s *Service) getOrCreateStateMachine(
	detail *DetailTx, initialState State) *StateMachine {
	s.txstatesmu.Lock()
	defer s.txstatesmu.Unlock()
	sm, ok := s.txstates[string(detail.ID)]
	if ok {
		return sm
	}
	cfg := &StateMachineConfig{
		Consensus1Action: s.onConsensusEvent,
		Consensus2Action: s.onConsensusEvent,
		Consensus3Action: s.onConsensusEvent,
		Table:            s.table,
		Detail:           detail,
		InitialState:     initialState,
	}
	sm = NewStateMachine(cfg)
	s.txstates[string(detail.ID)] = sm
	return sm
}

func (s *Service) gcStateMachines() {
	for {
		time.Sleep(1 * time.Second)
		s.txstatesmu.Lock()
		for k, v := range s.txstates {
			if v.State() == StateAborted || v.State() == StateSucceeded {
				if log.AtDebug() {
					log.Debug("removing statemachine", log.String("finale_state", v.State().String()), fld.TxID(ID([]byte(k))))
				}
				// v.Close()
				delete(s.txstates, k)
			}
		}
		s.txstatesmu.Unlock()
	}
}

func (s *Service) AddTransaction(
	ctx context.Context, tx *Transaction, evidences map[uint64][]byte,
) ([]*Object, error) {
	ids, err := MakeIDs(tx)
	if err != nil {
		log.Error("unable to create IDs", fld.Err(err))
		return nil, err
	}

	txbytes, _ := proto.Marshal(tx)
	if !s.verifySignatures(txbytes, evidences) {
		log.Error("invalid evidences from nodes")
		return nil, errors.New("sbac: invalid evidences from nodes")
	}
	log.Info("sbac: all evidence verified with success")

	objects := []*Object{}
	for _, v := range ids.TraceObjectPairs {
		objects = append(objects, v.OutputObjects...)
	}

	detail := DetailTx{ID: ids.TxID, RawTx: txbytes, Tx: tx}
	s.addStateMachine(&detail, StateWaitingForConsensus1)

	// broadcast transaction
	consensusTx := &ConsensusTransaction{
		TxID:      ids.TxID,
		Tx:        tx,
		Evidences: evidences,
		Op:        ConsensusOp_Consensus1,
		Initiator: s.nodeID,
	}
	b, err := proto.Marshal(consensusTx)
	if err != nil {
		return nil, fmt.Errorf("sbac: unable to marshal consensus tx: %v", err)
	}
	s.broadcaster.AddTransaction(b, 0)
	log.Info("sbac: transaction added successfully")
	return objects, nil
}

func (s *Service) addTransaction(ctx context.Context, payload []byte, id uint64) (*service.Message, error) {
	req := &AddTransactionRequest{}
	err := proto.Unmarshal(payload, req)
	if err != nil {
		log.Error("sbac: unable to unmarshal AddTransaction", fld.Err(err))
		return nil, fmt.Errorf("sbac: add_transaction unmarshaling error: %v", err)
	}

	ids, err := MakeIDs(req.Tx)
	if err != nil {
		log.Error("unable to create IDs", fld.Err(err))
		return nil, err
	}

	txbytes, _ := proto.Marshal(req.Tx)
	if !s.verifySignatures(txbytes, req.Evidences) {
		log.Error("invalid evidences from nodes")
		return nil, errors.New("sbac: invalid evidences from nodes")
	}
	log.Info("sbac: all evidence verified with success")

	objects := map[string]*ObjectList{}
	for _, v := range ids.TraceObjectPairs {
		trID := base64.StdEncoding.EncodeToString(v.Trace.ID)
		objects[trID] = &ObjectList{v.OutputObjects}
	}
	// txdetails := NewTxDetails(ids.TxID, txbytes, req.Tx, req.Evidences)

	detail := DetailTx{ID: ids.TxID, RawTx: txbytes, Tx: req.Tx}
	s.addStateMachine(&detail, StateWaitingForConsensus1)

	// broadcast transaction
	consensusTx := &ConsensusTransaction{
		TxID:      ids.TxID,
		Tx:        req.Tx,
		Evidences: req.Evidences,
		Op:        ConsensusOp_Consensus1,
		Initiator: s.nodeID,
	}
	b, err := proto.Marshal(consensusTx)
	if err != nil {
		return nil, fmt.Errorf("sbac: unable to marshal consensus tx: %v", err)
	}
	if s.isNodeInitiatingBroadcast(ID(ids.TxID)) {
		s.broadcaster.AddTransaction(b, 0)
	}
	res := &AddTransactionResponse{
		Objects: objects,
	}

	b, err = proto.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("sbac: unable to marshal add_transaction response, %v", err)
	}
	log.Info("sbac: transaction added successfully")
	return &service.Message{
		ID:      id,
		Opcode:  int32(Opcode_ADD_TRANSACTION),
		Payload: b,
	}, nil
}

func (s *Service) QueryObjectByVersionID(versionid []byte) ([]byte, error) {
	objects, err := GetObjects(s.store, [][]byte{versionid})
	if err != nil {
		return nil, err
	} else if len(objects) != 1 {
		return nil, fmt.Errorf("invalid number of objects found, expected %v found %v", 1, len(objects))
	}
	return objects[0].Value, nil
}

func (s *Service) publishObjects(ids *IDs, success bool) {
	for _, topair := range ids.TraceObjectPairs {
		for _, outo := range topair.OutputObjects {
			shard := s.top.ShardForVersionID(outo.GetVersionID())
			if shard == s.shardID {
				s.ps.Publish(outo.VersionID, outo.Labels, success)
			}
		}
	}
}

func (s *Service) saveLabels(ids *IDs) error {
	for _, topair := range ids.TraceObjectPairs {
		for _, outo := range topair.OutputObjects {
			shard := s.top.ShardForVersionID(outo.GetVersionID())
			if shard == s.shardID {
				for _, label := range outo.Labels {
					s.kvstore.Set([]byte(label), outo.VersionID)
				}
			}
		}
	}

	return nil
}

func queryPayload(id uint64, res *QueryObjectResponse) (*service.Message, error) {
	b, _ := proto.Marshal(res)
	return &service.Message{
		ID:      id,
		Opcode:  int32(Opcode_QUERY_OBJECT),
		Payload: b,
	}, nil
}

func (s *Service) queryObject(ctx context.Context, payload []byte, id uint64) (*service.Message, error) {
	req := &QueryObjectRequest{}
	err := proto.Unmarshal(payload, req)
	res := &QueryObjectResponse{}
	if err != nil {
		res.Error = fmt.Errorf("sbac: query_object unmarshaling error: %v", err).Error()
		return queryPayload(id, res)
	}

	if req.VersionID == nil {
		res.Error = fmt.Errorf("sbac: nil versionid").Error()
		return queryPayload(id, res)
	}
	objects, err := GetObjects(s.store, [][]byte{req.VersionID})
	if err != nil {
		res.Error = err.Error()
	} else if len(objects) != 1 {
		res.Error = fmt.Errorf("sbac: invalid number of objects found, expected %v found %v", 1, len(objects)).Error()
	} else {
		res.Object = objects[0]
	}
	return queryPayload(id, res)
}

func (s *Service) handleStates(ctx context.Context, payload []byte, id uint64) (*service.Message, error) {
	sr := []*StateReport{}
	s.txstatesmu.Lock()
	for _, v := range s.txstates {
		_ = v
		// sr = append(sr, v.StateReport())
	}
	s.txstatesmu.Unlock()
	res := &StatesReportResponse{
		States:        sr,
		EventsInQueue: int32(s.pe.Len()),
	}
	b, err := proto.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("sbac: unable to marshal states reports response")
	}
	return &service.Message{
		ID:      id,
		Opcode:  int32(Opcode_STATES),
		Payload: b,
	}, nil
}

func deletePayload(id uint64, res *DeleteObjectResponse) (*service.Message, error) {
	b, _ := proto.Marshal(res)
	return &service.Message{
		ID:      id,
		Opcode:  int32(Opcode_DELETE_OBJECT),
		Payload: b,
	}, nil
}

func (s *Service) deleteObject(ctx context.Context, payload []byte, id uint64) (*service.Message, error) {
	req := &DeleteObjectRequest{}
	res := &DeleteObjectResponse{}
	err := proto.Unmarshal(payload, req)
	if err != nil {
		res.Error = fmt.Errorf("sbac: remove_object unmarshaling error: %v", err).Error()
		return deletePayload(id, res)
	}

	if req.VersionID == nil {
		res.Error = fmt.Errorf("sbac: nil object versionid").Error()
		return deletePayload(id, res)
	}
	err = DeleteObjects(s.store, [][]byte{req.VersionID})
	if err != nil {
		res.Error = err.Error()
	}
	objects, err := GetObjects(s.store, [][]byte{req.VersionID})
	if err != nil {
		res.Error = err.Error()
	} else if len(objects) != 1 {
		res.Error = fmt.Errorf("sbac: invalid number of objects removed, expected %v found %v", 1, len(objects)).Error()
	} else {
		res.Object = objects[0]
	}
	return deletePayload(id, res)
}

func createPayload(id uint64, res *NewObjectResponse) (*service.Message, error) {
	b, _ := proto.Marshal(res)
	return &service.Message{
		ID:      id,
		Opcode:  int32(Opcode_CREATE_OBJECT),
		Payload: b,
	}, nil
}

func (s *Service) createObject(ctx context.Context, payload []byte, id uint64) (*service.Message, error) {
	req := &NewObjectRequest{}
	err := proto.Unmarshal(payload, req)
	res := &NewObjectResponse{}
	if err != nil {
		res.Error = fmt.Errorf("sbac: new_object unmarshaling error: %v", err).Error()
		return createPayload(id, res)
	}

	if req.Object == nil || len(req.Object) <= 0 {
		res.Error = fmt.Errorf("sbac: nil object").Error()
		return createPayload(id, res)
	}
	ch := combihash.New()
	ch.Write([]byte(req.Object))
	versionid := ch.Digest()
	if log.AtDebug() {
		log.Debug("sbac: creating new object", log.String("objet", string(req.Object)), log.Uint32("object.id", ID(versionid)))
	}
	o, err := CreateObject(s.store, versionid, req.Object)
	if err != nil {
		if log.AtDebug() {
			log.Debug("sbac: unable to create object", log.String("objet", string(req.Object)), log.Uint32("object.id", ID(versionid)), fld.Err(err))
		}
		res.Error = err.Error()
	} else {
		res.ID = o.VersionID
	}
	return createPayload(id, res)
}

func createObjectsPayload(id uint64, res *CreateObjectsResponse) (*service.Message, error) {
	b, _ := proto.Marshal(res)
	return &service.Message{
		ID:      id,
		Opcode:  int32(Opcode_CREATE_OBJECT),
		Payload: b,
	}, nil
}

func (s *Service) createObjects(ctx context.Context, payload []byte, id uint64) (*service.Message, error) {
	req := &CreateObjectsRequest{}
	err := proto.Unmarshal(payload, req)
	res := &CreateObjectsResponse{}
	if err != nil {
		res.Error = fmt.Errorf("sbac: new_object unmarshaling error: %v", err).Error()
		return createObjectsPayload(id, res)
	}

	if req.Objects == nil || len(req.Objects) <= 0 {
		res.Error = fmt.Errorf("sbac: nil object").Error()
		return createObjectsPayload(id, res)
	}

	ch := combihash.New()
	out := make([][]byte, 0, len(req.Objects))
	for _, object := range req.Objects {
		ch.Reset()
		ch.Write(object)
		versionid := ch.Digest()
		o, err := CreateObject(s.store, versionid, object)
		if err != nil {
			res.Error = err.Error()
			break
		}
		out = append(out, o.VersionID)
	}
	if len(res.Error) <= 0 {
		res.IDs = out
	}
	return createObjectsPayload(id, res)
}

func (s *Service) Name() string {
	return "sbac"
}

func (s *Service) Stop() error {
	return s.store.Close()
}

func New(cfg *Config) (*Service, error) {
	algorithm, err := signature.AlgorithmFromString(cfg.SigningKey.Type)
	if err != nil {
		return nil, err
	}
	privKeybytes, err := b32.DecodeString(cfg.SigningKey.Private)
	if err != nil {
		return nil, err
	}
	privkey, err := signature.LoadPrivateKey(algorithm, privKeybytes)
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions
	badgerPath := path.Join(cfg.Directory, badgerStorePath)
	opts.Dir = badgerPath
	opts.ValueDir = badgerPath
	store, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	s := &Service{
		broadcaster: cfg.Broadcaster,
		conns:       conns.NewPool(20, cfg.NodeID, cfg.Top, cfg.MaxPayload, cfg.Key, service.CONNECTION_SBAC),
		kvstore:     cfg.KVStore,
		nodeID:      cfg.NodeID,
		privkey:     privkey,
		ps:          cfg.Pubsub,
		top:         cfg.Top,
		txstates:    map[string]*StateMachine{},
		store:       store,
		shardID:     cfg.Top.ShardForNode(cfg.NodeID),
		shardCount:  cfg.ShardCount,
		shardSize:   cfg.ShardSize,
	}
	s.pe = NewPendingEvents(s.consumeEvents)
	s.table = s.makeStateTable()
	s.broadcaster.Register(s.handleDeliver)
	go s.pe.Run()
	// go s.gcStateMachines()
	return s, nil
}

func (s *Service) isNodeInitiatingBroadcast(txID uint32) bool {
	nodesInShard := s.top.NodesInShard(s.shardID)
	n := nodesInShard[txID%(uint32(len(nodesInShard)))]
	log.Debug("consensus will be started", log.Uint64("peer", n))
	return n == s.nodeID
}
