package transactor // import "chainspace.io/prototype/transactor"

import (
	"errors"
	"fmt"
	"hash/fnv"
	"time"

	"chainspace.io/prototype/broadcast"
	"chainspace.io/prototype/log"
	"chainspace.io/prototype/service"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

type State uint8

const (
	StateWaitingForConsensus1 State = iota // waiting for the first consensus to be reached, in order to process the transaction inside the shard.

	StateObjectLocked // check if objects exists and are actives

	StateAcceptPhase1Broadcasted // 2-phase commit phase 1
	StateRejectPhase1Broadcasted
	StateWaitingForPhase1

	StateConsensus2Triggered
	StateWaitingForConsensus2 // second consensus inside the shard, in order to confirm accept of the transaction after the first phase of 2-phase commit. also should kick bad nodes in the future

	StateAcceptPhase2Broadcasted // 2-phase commit phase 2
	StateRejectPhase2Broadcasted
	StateWaitingForPhase2

	StateObjectsDeactivated
	StateObjectsCreated

	StateCommitObjectsBroadcasted // notify other shards to create the objects
	StateWaitingForCommit         // node which need to create object start it state machine with this
	StateConsensusCommitTriggered
	StateWaitingForConsensusCommit

	StateAborted
	StateSucceeded
)

func (s State) String() string {
	switch s {
	case StateWaitingForConsensus1:
		return "StateWaitingOnConsensus1"
	case StateObjectLocked:
		return "StateObjectLocked"
	case StateRejectPhase1Broadcasted:
		return "StateRejectPhase1Broadcasted"
	case StateAcceptPhase1Broadcasted:
		return "StateAcceptPhase1Broadcasted"
	case StateWaitingForPhase1:
		return "StateWaitingForPhase1"
	case StateConsensus2Triggered:
		return "StateConsensus2Triggered"
	case StateWaitingForConsensus2:
		return "StateWaitingForConsensus2"
	case StateAcceptPhase2Broadcasted:
		return "StateAcceptPhase2Broadcasted"
	case StateRejectPhase2Broadcasted:
		return "StateRejectPhase2Broadcasted"
	case StateWaitingForPhase2:
		return "StateWaitingForPhase2"
	case StateObjectsDeactivated:
		return "StateObjectsDeactivated"
	case StateObjectsCreated:
		return "StateObjectsCreated"
	case StateSucceeded:
		return "StateSucceeded"
	case StateCommitObjectsBroadcasted:
		return "StateCommitObjectsBroadcasted"
	case StateConsensusCommitTriggered:
		return "StateConsensusCommitTriggered"
	case StateWaitingForCommit:
		return "StateWaitingForCommit"
	case StateWaitingForConsensusCommit:
		return "StateWaitingForConsensusCommit"
	case StateAborted:
		return "StateAborted"
	default:
		return "error"
	}
}

func (s *Service) makeStateTable() *StateTable {
	// change state are triggered on events
	actionTable := map[State]Action{
		StateWaitingForConsensus1:      s.onWaitingForConsensus1,
		StateWaitingForConsensus2:      s.onWaitingForConsensus2,
		StateWaitingForPhase1:          s.onWaitingForPhase1,
		StateWaitingForPhase2:          s.onWaitingForPhase2,
		StateWaitingForCommit:          s.onWaitingForCommit,
		StateWaitingForConsensusCommit: s.onWaitingForConsensusCommit,
	}

	// change state are triggererd from the previous state change
	transitionTable := map[StateTransition]Transition{
		// first bad path, no initial censensu reach / missing evidence, whatever
		{StateWaitingForConsensus1, StateRejectPhase1Broadcasted}: s.toRejectPhase1Broadcasted,
		{StateRejectPhase1Broadcasted, StateAborted}:              s.toAborted,
		{StateObjectLocked, StateAborted}:                         s.toAborted,

		// objects locked with success, and only 1 shard involved
		{StateObjectLocked, StateObjectsDeactivated}:   s.toObjectDeactivated,
		{StateObjectsDeactivated, StateObjectsCreated}: s.toObjectsCreated,
		{StateObjectsCreated, StateSucceeded}:          s.toSucceeded,

		// object locked with success, but multiple shards involved
		{StateWaitingForConsensus1, StateObjectLocked}:        s.toObjectLocked,
		{StateObjectLocked, StateAcceptPhase1Broadcasted}:     s.toAcceptPhase1Broadcasted,
		{StateAcceptPhase1Broadcasted, StateWaitingForPhase1}: s.toWaitingForPhase1,
		{StateWaitingForPhase1, StateAborted}:                 s.toAborted,
		{StateWaitingForPhase1, StateConsensus2Triggered}:     s.toConsensus2Triggered,
		// consensus 2 / phase 2
		{StateConsensus2Triggered, StateWaitingForConsensus2}:     s.toWaitingForConsensus2,
		{StateWaitingForConsensus2, StateAcceptPhase2Broadcasted}: s.toAcceptPhase2Broadcasted,
		{StateWaitingForConsensus2, StateRejectPhase2Broadcasted}: s.toRejectPhase2Broadcasted,
		{StateRejectPhase2Broadcasted, StateAborted}:              s.toAborted,
		{StateAcceptPhase2Broadcasted, StateWaitingForPhase2}:     s.toWaitingForPhase2,
		{StateWaitingForPhase2, StateAborted}:                     s.toAborted,
		{StateWaitingForPhase2, StateObjectsDeactivated}:          s.toObjectDeactivated,

		{StateObjectsCreated, StateCommitObjectsBroadcasted}: s.toCommitObjectsBroadcasted,
		{StateCommitObjectsBroadcasted, StateAborted}:        s.toAborted,
		{StateCommitObjectsBroadcasted, StateSucceeded}:      s.toSucceeded,

		{StateWaitingForCommit, StateConsensusCommitTriggered}:          s.toConsensusCommitTriggered,
		{StateConsensusCommitTriggered, StateWaitingForConsensusCommit}: s.toWaitingForConsensusCommit,
		{StateWaitingForConsensusCommit, StateObjectsCreated}:           s.toObjectsCreated,
		{StateWaitingForConsensusCommit, StateAborted}:                  s.toAborted,
	}

	return &StateTable{actionTable, transitionTable, s.onEvent}
}

func (s *Service) onEvent(tx *TxDetails, event *Event) error {
	if string(tx.ID) != string(event.msg.Tx.ID) {
		log.Error("invalid transaction sent to state machine", zap.Uint32("expected", tx.HashID), zap.Uint32("got", ID(event.msg.Tx.ID)))
		return errors.New("invalid transaction ID")
	}
	log.Info(SBACOpcode_name[int32(event.msg.Op)]+" decision received", zap.Uint32("tx.id", tx.HashID), zap.String("decision", SBACDecision_name[int32(event.msg.Decision)]), zap.Uint64("peer.id", event.peerID))
	switch event.msg.Op {
	case SBACOpcode_PHASE1:
		tx.Phase1Decisions[event.peerID] = event.msg.Decision
	case SBACOpcode_PHASE2:
		tx.Phase2Decisions[event.peerID] = event.msg.Decision
	case SBACOpcode_CONSENSUS1:
		tx.Consensus1Tx = event.msg.Tx
	case SBACOpcode_CONSENSUS2:
		tx.Consensus2Tx = event.msg.Tx
	case SBACOpcode_COMMIT:
		tx.CommitDecisions[event.peerID] = event.msg.Decision
	case SBACOpcode_CONSENSUS_COMMIT:
		tx.ConsensusCommitTx = event.msg.Tx
	default:
	}
	return nil
}

// shardsInvolvedInTx return a list of IDs of all shards involved in the transaction either by
// holding state for an input object or input reference.
func (s *Service) shardsInvolvedInTx(tx *Transaction) []uint64 {
	uniqids := map[uint64]struct{}{}
	for _, trace := range tx.Traces {
		for _, obj := range trace.InputObjectsKeys {
			uniqids[s.top.ShardForKey(obj)] = struct{}{}
		}
		for _, ref := range trace.InputReferencesKeys {
			uniqids[s.top.ShardForKey(ref)] = struct{}{}
		}
	}
	ids := make([]uint64, 0, len(uniqids))
	for k, _ := range uniqids {
		ids = append(ids, k)
	}
	return ids
}

func twotplusone(shardSize uint64) uint64 {
	return (2*(shardSize/3) + 1)
}

func tplusone(shardSize uint64) uint64 {
	return shardSize/3 + 1
}

type WaitingDecisionResult uint8

const (
	WaitingDecisionResultAbort WaitingDecisionResult = iota
	WaitingDecisionResultPending
	WaitingDecisionResultAccept
)

func (s *Service) onWaitingFor(tx *TxDetails, decisions map[uint64]SBACDecision, phaseName string) WaitingDecisionResult {
	shards := s.shardsInvolvedInTx(tx.Tx)
	var somePending bool
	// for each shards, get the nodes id, and checks if they answered
	vtwotplusone := twotplusone(s.shardSize)
	vtplusone := tplusone(s.shardSize)
	for _, v := range shards {
		nodes := s.top.NodesInShard(v)
		var accepted uint64
		var rejected uint64
		for _, nodeID := range nodes {
			if d, ok := decisions[nodeID]; ok {
				if d == SBACDecision_ACCEPT {
					accepted += 1
					continue
				}
				rejected += 1
			}
		}
		if rejected >= vtplusone {
			log.Info(phaseName+" transaction rejected",
				zap.Uint32("tx.id", tx.HashID),
				zap.Uint64("peer.shard", v),
				zap.Uint64("t+1", vtplusone),
				zap.Uint64("rejected", rejected),
			)
			return WaitingDecisionResultAbort
		}
		if accepted >= vtwotplusone {
			log.Info(phaseName+" transaction accepted",
				zap.Uint32("tx.id", tx.HashID),
				zap.Uint64("peer.shard", v),
				zap.Uint64s("shards_involved", shards),
				zap.Uint64("2t+1", vtwotplusone),
				zap.Uint64("accepted", accepted),
			)
			continue
		}
		somePending = true
	}

	if somePending {
		log.Info(phaseName+" transaction pending, not enough answers from shards", zap.Uint32("tx.id", tx.HashID))
		return WaitingDecisionResultPending
	}

	log.Info(phaseName+" transaction accepted by all shards", zap.Uint32("tx.id", tx.HashID))
	return WaitingDecisionResultAccept
}

func (s *Service) onWaitingForPhase1(tx *TxDetails) (State, error) {
	switch s.onWaitingFor(tx, tx.Phase1Decisions, "phase1") {
	case WaitingDecisionResultAbort:
		return StateAborted, nil
	case WaitingDecisionResultPending:
		return StateWaitingForPhase1, nil
	default:
		return StateConsensus2Triggered, nil
	}
}

func (s *Service) onWaitingForPhase2(tx *TxDetails) (State, error) {
	switch s.onWaitingFor(tx, tx.Phase2Decisions, "phase2") {
	case WaitingDecisionResultAbort:
		return StateAborted, nil
	case WaitingDecisionResultPending:
		return StateWaitingForPhase2, nil
	default:
		return StateObjectsDeactivated, nil
	}
}

func (s *Service) onWaitingForCommit(tx *TxDetails) (State, error) {
	switch s.onWaitingFor(tx, tx.CommitDecisions, "commit") {
	case WaitingDecisionResultAbort:
		return StateAborted, nil
	case WaitingDecisionResultPending:
		return StateWaitingForCommit, nil
	default:
		return StateConsensusCommitTriggered, nil
	}
}

func (s *Service) objectsExists(objs, refs [][]byte) ([]*Object, bool) {
	keys := [][]byte{}
	for _, v := range append(objs, refs...) {
		if s.top.ShardForKey(v) == s.shardID {
			keys = append(keys, v)
		}
	}

	objects, err := GetObjects(s.store, keys)
	if err != nil {
		return nil, false
	}
	return objects, true
}

func (s *Service) onWaitingForConsensus1(tx *TxDetails) (State, error) {
	// TODO: need to check evidences here, as this step should happend after the first consensus round.
	// not sure how we agregate evidence here, if we get new ones from the consensus rounds or the same that where sent with the transaction at the begining.
	if tx.Consensus1Tx == nil {
		return StateWaitingForConsensus1, nil
	}
	if !s.verifySignatures(tx.ID, tx.Consensus1Tx.GetEvidences()) {
		log.Error("consensus1 missing/invalid signatures", zap.Uint32("tx.id", tx.HashID))
		return StateRejectPhase1Broadcasted, nil
	}

	// check that all inputs objects and references part of the state of this node exists.
	for _, trace := range tx.Tx.Traces {
		objects, ok := s.objectsExists(trace.InputObjectsKeys, trace.InputReferencesKeys)
		if !ok {
			log.Error("consensus1 some objects do not exists", zap.Uint32("tx.id", tx.HashID))
			return StateRejectPhase1Broadcasted, nil
		}
		for _, v := range objects {
			if v.Status == ObjectStatus_INACTIVE {
				log.Error("consensus1 some objects are inactive", zap.Uint32("tx.id", tx.HashID))
				return StateRejectPhase1Broadcasted, nil
			}
		}
	}

	log.Info("consensus1 evidences and input objects/references checked successfully", zap.Uint32("tx.id", tx.HashID))
	return StateObjectLocked, nil
}

// TODO(): kick bad node here ? not sure which one are bad
func (s *Service) onWaitingForConsensus2(tx *TxDetails) (State, error) {
	if tx.Consensus2Tx == nil {
		return StateWaitingForConsensus2, nil
	}
	if !s.verifySignatures(tx.ID, tx.Consensus2Tx.GetEvidences()) {
		log.Error("consensus2 missing/invalid signatures", zap.Uint32("tx.id", tx.HashID))
		return StateRejectPhase2Broadcasted, nil
	}

	return StateAcceptPhase2Broadcasted, nil
}

// TODO(): verify signatures here, but need to be passed to first.
func (s *Service) onWaitingForConsensusCommit(tx *TxDetails) (State, error) {
	if tx.ConsensusCommitTx == nil {
		return StateWaitingForConsensusCommit, nil
	}
	if !s.verifySignatures(tx.ID, tx.ConsensusCommitTx.GetEvidences()) {
		log.Error("consensus_commit missing/invalid signatures", zap.Uint32("tx.id", tx.HashID))
		return StateAborted, nil

	}
	return StateObjectsCreated, nil

}

func (s *Service) inputObjectsForShard(shardID uint64, tx *Transaction) (objects [][]byte, allInShard bool) {
	// get all objects part of this current shard state
	allInShard = true
	for _, t := range tx.Traces {
		for _, o := range t.InputObjectsKeys {
			o := o
			if shardID := s.top.ShardForKey(o); shardID == s.shardID {
				objects = append(objects, o)
				continue
			}
			allInShard = false
		}
		for _, ref := range t.InputReferencesKeys {
			ref := ref
			if shardID := s.top.ShardForKey(ref); shardID == s.shardID {
				continue
			}
			allInShard = false
		}
	}
	return
}

// can abort, if impossible to lock object
// then check if one shard or more is involved and return StateAcceptBroadcasted
// or StateObjectSetInactive
func (s *Service) toObjectLocked(tx *TxDetails) (State, error) {
	objects, allInShard := s.inputObjectsForShard(s.shardID, tx.Tx)
	// lock them
	if err := LockObjects(s.store, objects); err != nil {
		log.Error("unable to lock all objects", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		// return nil from here as we can abort as a valid transition
		return StateAborted, nil
	}
	if allInShard {
		return StateObjectsDeactivated, nil
	}
	return StateAcceptPhase1Broadcasted, nil
}

func makeMessage(m *SBACMessage) (*service.Message, error) {
	payload, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &service.Message{
		Opcode:  uint32(Opcode_SBAC),
		Payload: payload,
	}, nil
}

func (s *Service) sendToAllShardInvolved(tx *TxDetails, msg *service.Message) error {
	shards := s.shardsInvolvedInTx(tx.Tx)
	for _, shard := range shards {
		nodes := s.top.NodesInShard(shard)
		for _, node := range nodes {
			// TODO: proper timeout ?
			_, err := s.conns.WriteRequest(node, msg, time.Second)
			if err != nil {
				log.Error("unable to connect to node", zap.Uint32("tx.id", tx.HashID), zap.Uint64("peer.id", node))
				return fmt.Errorf("unable to connect to node(%v): %v", node, err)
			}
		}
	}
	return nil
}

func (s *Service) toRejectPhase1Broadcasted(tx *TxDetails) (State, error) {
	sbacmsg := &SBACMessage{
		Op:       SBACOpcode_PHASE1,
		Decision: SBACDecision_REJECT,
		PeerID:   s.nodeID,
		Tx: &SBACTransaction{
			ID:        tx.ID,
			Tx:        tx.Tx,
			Evidences: tx.CheckersEvidences,
			Op:        SBACOpcode_PHASE1,
		},
	}
	msg, err := makeMessage(sbacmsg)
	if err != nil {
		log.Error("unable to serialize message reject accept transaction", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	err = s.sendToAllShardInvolved(tx, msg)
	if err != nil {
		log.Error("unable to sent reject transaction to all shards", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
	}

	log.Info("reject transaction sent to all shards", zap.Uint32("tx.id", tx.HashID))
	return StateAborted, err
}

func (s *Service) toAcceptPhase1Broadcasted(tx *TxDetails) (State, error) {
	sbacmsg := &SBACMessage{
		Op:       SBACOpcode_PHASE1,
		Decision: SBACDecision_ACCEPT,
		PeerID:   s.nodeID,
		Tx: &SBACTransaction{
			ID:        tx.ID,
			Tx:        tx.Tx,
			Evidences: tx.CheckersEvidences,
			Op:        SBACOpcode_PHASE1,
		},
	}
	msg, err := makeMessage(sbacmsg)
	if err != nil {
		log.Error("unable to serialize message accept transaction accept", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	err = s.sendToAllShardInvolved(tx, msg)
	if err != nil {
		log.Error("unable to sent accept transaction to all shards", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	log.Info("accept transaction sent to all shards", zap.Uint32("tx.id", tx.HashID))
	return StateWaitingForPhase1, nil
}

func (s *Service) toRejectPhase2Broadcasted(tx *TxDetails) (State, error) {
	sbacmsg := &SBACMessage{
		Op:       SBACOpcode_PHASE2,
		Decision: SBACDecision_REJECT,
		PeerID:   s.nodeID,
		Tx: &SBACTransaction{
			ID:        tx.ID,
			Tx:        tx.Tx,
			Evidences: tx.CheckersEvidences,
			Op:        SBACOpcode_PHASE2,
		},
	}
	msg, err := makeMessage(sbacmsg)
	if err != nil {
		log.Error("phase2 reject unable to serialize message", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	err = s.sendToAllShardInvolved(tx, msg)
	if err != nil {
		log.Error("phase2 reject unable to sent reject transaction to all shards", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
	}

	log.Info("phase2 reject transaction sent to all shards", zap.Uint32("tx.id", tx.HashID))
	return StateAborted, err
}

func (s *Service) toAcceptPhase2Broadcasted(tx *TxDetails) (State, error) {
	sbacmsg := &SBACMessage{
		Op:       SBACOpcode_PHASE2,
		Decision: SBACDecision_ACCEPT,
		PeerID:   s.nodeID,
		Tx: &SBACTransaction{
			ID:        tx.ID,
			Tx:        tx.Tx,
			Evidences: tx.CheckersEvidences,
			Op:        SBACOpcode_PHASE2,
		},
	}
	msg, err := makeMessage(sbacmsg)
	if err != nil {
		log.Error("phase2 accept unable to serialize message", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	err = s.sendToAllShardInvolved(tx, msg)
	if err != nil {
		log.Error("phase2 unable to sent accept transaction to all shards", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	log.Info("phase2 accept transaction sent to all shards", zap.Uint32("tx.id", tx.HashID))
	return StateWaitingForPhase2, nil
}

func (s *Service) toWaitingForPhase1(tx *TxDetails) (State, error) {
	return StateWaitingForPhase1, nil
}

func (s *Service) toWaitingForPhase2(tx *TxDetails) (State, error) {
	return StateWaitingForPhase2, nil
}

func (s *Service) toObjectDeactivated(tx *TxDetails) (State, error) {
	objects, _ := s.inputObjectsForShard(s.shardID, tx.Tx)
	// lock them
	if err := DeactivateObjects(s.store, objects); err != nil {
		log.Error("unable to deactivate all objects", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		// return nil from here as we can abort as a valid transition
		return StateAborted, nil
	}

	log.Info("all object deactivated successfully", zap.Uint32("tx.id", tx.HashID))
	return StateObjectsCreated, nil
}

func (s *Service) toObjectsCreated(tx *TxDetails) (State, error) {
	traceIDPairs, err := MakeTraceIDs(tx.Tx.Traces)
	if err != nil {
		return StateAborted, err
	}
	traceObjectPairs, err := MakeTraceObjectPairs(traceIDPairs)
	if err != nil {
		return StateAborted, err
	}
	objects := []*Object{}
	allObjectsInCurrentShard := true
	for _, v := range traceObjectPairs {
		v := v
		for _, o := range v.OutputObjects {
			o := o
			if shardID := s.top.ShardForKey(o.Key); shardID == s.shardID {
				objects = append(objects, o)
				continue
			}
			allObjectsInCurrentShard = false
		}
	}
	err = CreateObjects(s.store, objects)
	if err != nil {
		log.Error("unable to create objects", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}
	log.Info("all objects created successfully", zap.Uint32("tx.id", tx.HashID))
	if allObjectsInCurrentShard {
		return StateSucceeded, nil
	}
	return StateCommitObjectsBroadcasted, nil
}

func (s *Service) toSucceeded(tx *TxDetails) (State, error) {
	tx.Result <- true
	return StateSucceeded, nil
}

func (s *Service) toAborted(tx *TxDetails) (State, error) {
	// unlock any objects maybe related to this transaction.
	objects, _ := s.inputObjectsForShard(s.shardID, tx.Tx)
	err := UnlockObjects(s.store, objects)
	if err != nil {
		log.Error("unable to unlock objects", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
	}
	tx.Result <- false
	return StateAborted, nil
}

func (s *Service) sendToShards(shards []uint64, tx *TxDetails, msg *service.Message) error {
	for _, shard := range shards {
		nodes := s.top.NodesInShard(shard)
		for _, node := range nodes {
			// TODO: proper timeout ?
			_, err := s.conns.WriteRequest(node, msg, time.Second)
			if err != nil {
				log.Error("unable to connect to node", zap.Uint32("tx.id", tx.HashID), zap.Uint64("peer.id", node))
				return fmt.Errorf("unable to connect to node(%v): %v", node, err)
			}
		}
	}
	return nil
}

func (s *Service) toCommitObjectsBroadcasted(tx *TxDetails) (State, error) {
	sbacmsg := &SBACMessage{
		Op:       SBACOpcode_COMMIT,
		Decision: SBACDecision_ACCEPT,
		PeerID:   s.nodeID,
		Tx: &SBACTransaction{
			ID:        tx.ID,
			Tx:        tx.Tx,
			Evidences: tx.CheckersEvidences,
			Op:        SBACOpcode_COMMIT,
		},
	}
	msg, err := makeMessage(sbacmsg)
	if err != nil {
		log.Error("commit accept unable to serialize message", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	// send to all shards which will hold state for new object, and was not involved before
	shards := s.shardsInvolvedInTx(tx.Tx)
	shardsUniq := IDSlice(shards).ToMapUniq()
	ids, err := MakeIDs(tx.Tx)
	if err != nil {
		return StateAborted, err
	}
	shardsInvolvedUniq := map[uint64]struct{}{}
	for _, v := range ids.TraceObjectPairs {
		v := v
		for _, o := range v.OutputObjects {
			shard := s.top.ShardForKey(o.Key)
			if _, ok := shardsUniq[shard]; !ok {
				shardsInvolvedUniq[shard] = struct{}{}
			}

		}
	}
	shardsInvolved := IDSlice{}.FromMapUniq(shardsInvolvedUniq)

	err = s.sendToShards(shardsInvolved, tx, msg)
	if err != nil {
		log.Error("commit unable to sent accept transaction to all shards", zap.Uint32("tx.id", tx.HashID), zap.Error(err))
		return StateAborted, err
	}

	// TODO(): this should be blocking here waiting for the shards to send us back the response so we can return an error

	log.Info("commit accept transaction sent to all shards",
		zap.Uint32("tx.id", tx.HashID), zap.Uint64s("shards_involved", shardsInvolved))
	return StateSucceeded, nil
}

// TODO(): should we make our own evidences ourselves here ?
func (s *Service) toConsensus2Triggered(tx *TxDetails) (State, error) {
	// broadcast transaction
	consensusTx := &SBACTransaction{
		Tx:        tx.Tx,
		ID:        tx.ID,
		Evidences: tx.CheckersEvidences,
		Op:        SBACOpcode_CONSENSUS2,
	}
	b, err := proto.Marshal(consensusTx)
	if err != nil {
		return StateAborted, fmt.Errorf("transactor: unable to marshal consensus tx: %v", err)
	}
	s.broadcaster.AddTransaction(&broadcast.TransactionData{Data: b})
	return StateWaitingForConsensus2, nil
}

func (s *Service) toConsensusCommitTriggered(tx *TxDetails) (State, error) {
	// broadcast transaction
	consensusTx := &SBACTransaction{
		Tx:        tx.Tx,
		ID:        tx.ID,
		Evidences: tx.CheckersEvidences,
		Op:        SBACOpcode_CONSENSUS_COMMIT,
	}
	b, err := proto.Marshal(consensusTx)
	if err != nil {
		return StateAborted, fmt.Errorf("transactor: unable to marshal consensus tx: %v", err)
	}
	s.broadcaster.AddTransaction(&broadcast.TransactionData{Data: b})
	return StateWaitingForConsensusCommit, nil
}

func (s *Service) toWaitingForConsensusCommit(tx *TxDetails) (State, error) {
	return StateWaitingForConsensusCommit, nil
}

func (s *Service) toWaitingForConsensus2(tx *TxDetails) (State, error) {
	return StateWaitingForConsensus2, nil
}

func ID(data []byte) uint32 {
	h := fnv.New32()
	h.Write(data)
	return h.Sum32()
}

type IDSlice []uint64

func (s IDSlice) ToMapUniq() map[uint64]struct{} {
	out := map[uint64]struct{}{}
	for _, v := range s {
		out[v] = struct{}{}
	}
	return out
}

func (s IDSlice) FromMapUniq(m map[uint64]struct{}) []uint64 {
	out := []uint64{}
	for k, _ := range m {
		out = append(out, k)
	}
	return out
}
