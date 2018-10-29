// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types.proto

package sbac

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Opcode int32

const (
	Opcode_UNKNOWN         Opcode = 0
	Opcode_ADD_TRANSACTION Opcode = 1
	Opcode_QUERY_OBJECT    Opcode = 2
	Opcode_CREATE_OBJECT   Opcode = 3
	Opcode_DELETE_OBJECT   Opcode = 4
	Opcode_STATES          Opcode = 5
	Opcode_SBAC            Opcode = 6
	Opcode_CREATE_OBJECTS  Opcode = 7
)

var Opcode_name = map[int32]string{
	0: "UNKNOWN",
	1: "ADD_TRANSACTION",
	2: "QUERY_OBJECT",
	3: "CREATE_OBJECT",
	4: "DELETE_OBJECT",
	5: "STATES",
	6: "SBAC",
	7: "CREATE_OBJECTS",
}
var Opcode_value = map[string]int32{
	"UNKNOWN":         0,
	"ADD_TRANSACTION": 1,
	"QUERY_OBJECT":    2,
	"CREATE_OBJECT":   3,
	"DELETE_OBJECT":   4,
	"STATES":          5,
	"SBAC":            6,
	"CREATE_OBJECTS":  7,
}

func (x Opcode) String() string {
	return proto.EnumName(Opcode_name, int32(x))
}
func (Opcode) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{0} }

type ObjectStatus int32

const (
	ObjectStatus_ACTIVE   ObjectStatus = 0
	ObjectStatus_INACTIVE ObjectStatus = 1
	ObjectStatus_LOCKED   ObjectStatus = 2
)

var ObjectStatus_name = map[int32]string{
	0: "ACTIVE",
	1: "INACTIVE",
	2: "LOCKED",
}
var ObjectStatus_value = map[string]int32{
	"ACTIVE":   0,
	"INACTIVE": 1,
	"LOCKED":   2,
}

func (x ObjectStatus) String() string {
	return proto.EnumName(ObjectStatus_name, int32(x))
}
func (ObjectStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{1} }

type SBACOpcode int32

const (
	SBACOpcode_CONSENSUS1       SBACOpcode = 0
	SBACOpcode_CONSENSUS2       SBACOpcode = 1
	SBACOpcode_PHASE1           SBACOpcode = 2
	SBACOpcode_PHASE2           SBACOpcode = 3
	SBACOpcode_COMMIT           SBACOpcode = 4
	SBACOpcode_CONSENSUS_COMMIT SBACOpcode = 5
)

var SBACOpcode_name = map[int32]string{
	0: "CONSENSUS1",
	1: "CONSENSUS2",
	2: "PHASE1",
	3: "PHASE2",
	4: "COMMIT",
	5: "CONSENSUS_COMMIT",
}
var SBACOpcode_value = map[string]int32{
	"CONSENSUS1":       0,
	"CONSENSUS2":       1,
	"PHASE1":           2,
	"PHASE2":           3,
	"COMMIT":           4,
	"CONSENSUS_COMMIT": 5,
}

func (x SBACOpcode) String() string {
	return proto.EnumName(SBACOpcode_name, int32(x))
}
func (SBACOpcode) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{2} }

type SBACDecision int32

const (
	SBACDecision_ACCEPT SBACDecision = 0
	SBACDecision_REJECT SBACDecision = 1
)

var SBACDecision_name = map[int32]string{
	0: "ACCEPT",
	1: "REJECT",
}
var SBACDecision_value = map[string]int32{
	"ACCEPT": 0,
	"REJECT": 1,
}

func (x SBACDecision) String() string {
	return proto.EnumName(SBACDecision_name, int32(x))
}
func (SBACDecision) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{3} }

type ConsensusOp int32

const (
	ConsensusOp_Consensus1 ConsensusOp = 0
	ConsensusOp_Consensus2 ConsensusOp = 1
	ConsensusOp_Consensus3 ConsensusOp = 2
)

var ConsensusOp_name = map[int32]string{
	0: "Consensus1",
	1: "Consensus2",
	2: "Consensus3",
}
var ConsensusOp_value = map[string]int32{
	"Consensus1": 0,
	"Consensus2": 1,
	"Consensus3": 2,
}

func (x ConsensusOp) String() string {
	return proto.EnumName(ConsensusOp_name, int32(x))
}
func (ConsensusOp) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{4} }

type SBACOp int32

const (
	SBACOp_Phase1      SBACOp = 0
	SBACOp_Phase2      SBACOp = 1
	SBACOp_PhaseCommit SBACOp = 2
)

var SBACOp_name = map[int32]string{
	0: "Phase1",
	1: "Phase2",
	2: "PhaseCommit",
}
var SBACOp_value = map[string]int32{
	"Phase1":      0,
	"Phase2":      1,
	"PhaseCommit": 2,
}

func (x SBACOp) String() string {
	return proto.EnumName(SBACOp_name, int32(x))
}
func (SBACOp) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{5} }

type Transaction struct {
	Traces []*Trace `protobuf:"bytes,1,rep,name=traces" json:"traces,omitempty"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{0} }

func (m *Transaction) GetTraces() []*Trace {
	if m != nil {
		return m.Traces
	}
	return nil
}

type Strings struct {
	Strs []string `protobuf:"bytes,1,rep,name=strs" json:"strs,omitempty"`
}

func (m *Strings) Reset()                    { *m = Strings{} }
func (m *Strings) String() string            { return proto.CompactTextString(m) }
func (*Strings) ProtoMessage()               {}
func (*Strings) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{1} }

func (m *Strings) GetStrs() []string {
	if m != nil {
		return m.Strs
	}
	return nil
}

type Trace struct {
	ContractID               string     `protobuf:"bytes,1,opt,name=contractID,proto3" json:"contractID,omitempty"`
	Procedure                string     `protobuf:"bytes,2,opt,name=procedure,proto3" json:"procedure,omitempty"`
	InputObjectVersionIDs    [][]byte   `protobuf:"bytes,3,rep,name=inputObjectVersionIDs" json:"inputObjectVersionIDs,omitempty"`
	InputReferenceVersionIDs [][]byte   `protobuf:"bytes,4,rep,name=inputReferenceVersionIDs" json:"inputReferenceVersionIDs,omitempty"`
	OutputObjects            [][]byte   `protobuf:"bytes,5,rep,name=outputObjects" json:"outputObjects,omitempty"`
	Parameters               [][]byte   `protobuf:"bytes,6,rep,name=parameters" json:"parameters,omitempty"`
	Returns                  [][]byte   `protobuf:"bytes,7,rep,name=returns" json:"returns,omitempty"`
	Labels                   []*Strings `protobuf:"bytes,8,rep,name=labels" json:"labels,omitempty"`
	Dependencies             []*Trace   `protobuf:"bytes,9,rep,name=dependencies" json:"dependencies,omitempty"`
	InputObjects             [][]byte   `protobuf:"bytes,10,rep,name=inputObjects" json:"inputObjects,omitempty"`
	InputReferences          [][]byte   `protobuf:"bytes,11,rep,name=inputReferences" json:"inputReferences,omitempty"`
}

func (m *Trace) Reset()                    { *m = Trace{} }
func (m *Trace) String() string            { return proto.CompactTextString(m) }
func (*Trace) ProtoMessage()               {}
func (*Trace) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{2} }

func (m *Trace) GetContractID() string {
	if m != nil {
		return m.ContractID
	}
	return ""
}

func (m *Trace) GetProcedure() string {
	if m != nil {
		return m.Procedure
	}
	return ""
}

func (m *Trace) GetInputObjectVersionIDs() [][]byte {
	if m != nil {
		return m.InputObjectVersionIDs
	}
	return nil
}

func (m *Trace) GetInputReferenceVersionIDs() [][]byte {
	if m != nil {
		return m.InputReferenceVersionIDs
	}
	return nil
}

func (m *Trace) GetOutputObjects() [][]byte {
	if m != nil {
		return m.OutputObjects
	}
	return nil
}

func (m *Trace) GetParameters() [][]byte {
	if m != nil {
		return m.Parameters
	}
	return nil
}

func (m *Trace) GetReturns() [][]byte {
	if m != nil {
		return m.Returns
	}
	return nil
}

func (m *Trace) GetLabels() []*Strings {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Trace) GetDependencies() []*Trace {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

func (m *Trace) GetInputObjects() [][]byte {
	if m != nil {
		return m.InputObjects
	}
	return nil
}

func (m *Trace) GetInputReferences() [][]byte {
	if m != nil {
		return m.InputReferences
	}
	return nil
}

type Object struct {
	VersionID []byte       `protobuf:"bytes,1,opt,name=versionID,proto3" json:"versionID,omitempty"`
	Value     []byte       `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Status    ObjectStatus `protobuf:"varint,3,opt,name=status,proto3,enum=sbac.ObjectStatus" json:"status,omitempty"`
	Labels    []string     `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty"`
}

func (m *Object) Reset()                    { *m = Object{} }
func (m *Object) String() string            { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()               {}
func (*Object) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{3} }

func (m *Object) GetVersionID() []byte {
	if m != nil {
		return m.VersionID
	}
	return nil
}

func (m *Object) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Object) GetStatus() ObjectStatus {
	if m != nil {
		return m.Status
	}
	return ObjectStatus_ACTIVE
}

func (m *Object) GetLabels() []string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type AddTransactionRequest struct {
	Tx        *Transaction      `protobuf:"bytes,1,opt,name=tx" json:"tx,omitempty"`
	Evidences map[uint64][]byte `protobuf:"bytes,2,rep,name=evidences" json:"evidences,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *AddTransactionRequest) Reset()                    { *m = AddTransactionRequest{} }
func (m *AddTransactionRequest) String() string            { return proto.CompactTextString(m) }
func (*AddTransactionRequest) ProtoMessage()               {}
func (*AddTransactionRequest) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{4} }

func (m *AddTransactionRequest) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *AddTransactionRequest) GetEvidences() map[uint64][]byte {
	if m != nil {
		return m.Evidences
	}
	return nil
}

type ObjectTraceIDPair struct {
	TraceID []byte    `protobuf:"bytes,1,opt,name=traceID,proto3" json:"traceID,omitempty"`
	Objects []*Object `protobuf:"bytes,2,rep,name=objects" json:"objects,omitempty"`
}

func (m *ObjectTraceIDPair) Reset()                    { *m = ObjectTraceIDPair{} }
func (m *ObjectTraceIDPair) String() string            { return proto.CompactTextString(m) }
func (*ObjectTraceIDPair) ProtoMessage()               {}
func (*ObjectTraceIDPair) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{5} }

func (m *ObjectTraceIDPair) GetTraceID() []byte {
	if m != nil {
		return m.TraceID
	}
	return nil
}

func (m *ObjectTraceIDPair) GetObjects() []*Object {
	if m != nil {
		return m.Objects
	}
	return nil
}

type ObjectList struct {
	List []*Object `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *ObjectList) Reset()                    { *m = ObjectList{} }
func (m *ObjectList) String() string            { return proto.CompactTextString(m) }
func (*ObjectList) ProtoMessage()               {}
func (*ObjectList) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{6} }

func (m *ObjectList) GetList() []*Object {
	if m != nil {
		return m.List
	}
	return nil
}

type AddTransactionResponse struct {
	Objects map[string]*ObjectList `protobuf:"bytes,1,rep,name=objects" json:"objects,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *AddTransactionResponse) Reset()                    { *m = AddTransactionResponse{} }
func (m *AddTransactionResponse) String() string            { return proto.CompactTextString(m) }
func (*AddTransactionResponse) ProtoMessage()               {}
func (*AddTransactionResponse) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{7} }

func (m *AddTransactionResponse) GetObjects() map[string]*ObjectList {
	if m != nil {
		return m.Objects
	}
	return nil
}

type StateReport struct {
	HashID          uint32          `protobuf:"varint,1,opt,name=hashID,proto3" json:"hashID,omitempty"`
	CommitDecisions map[uint64]bool `protobuf:"bytes,2,rep,name=commitDecisions" json:"commitDecisions,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Phase1Decisions map[uint64]bool `protobuf:"bytes,3,rep,name=phase1Decisions" json:"phase1Decisions,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Phase2Decisions map[uint64]bool `protobuf:"bytes,4,rep,name=phase2Decisions" json:"phase2Decisions,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	State           string          `protobuf:"bytes,5,opt,name=state,proto3" json:"state,omitempty"`
	PendingEvents   int32           `protobuf:"varint,6,opt,name=pendingEvents,proto3" json:"pendingEvents,omitempty"`
}

func (m *StateReport) Reset()                    { *m = StateReport{} }
func (m *StateReport) String() string            { return proto.CompactTextString(m) }
func (*StateReport) ProtoMessage()               {}
func (*StateReport) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{8} }

func (m *StateReport) GetHashID() uint32 {
	if m != nil {
		return m.HashID
	}
	return 0
}

func (m *StateReport) GetCommitDecisions() map[uint64]bool {
	if m != nil {
		return m.CommitDecisions
	}
	return nil
}

func (m *StateReport) GetPhase1Decisions() map[uint64]bool {
	if m != nil {
		return m.Phase1Decisions
	}
	return nil
}

func (m *StateReport) GetPhase2Decisions() map[uint64]bool {
	if m != nil {
		return m.Phase2Decisions
	}
	return nil
}

func (m *StateReport) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *StateReport) GetPendingEvents() int32 {
	if m != nil {
		return m.PendingEvents
	}
	return 0
}

type StatesReportResponse struct {
	States        []*StateReport `protobuf:"bytes,1,rep,name=states" json:"states,omitempty"`
	EventsInQueue int32          `protobuf:"varint,2,opt,name=eventsInQueue,proto3" json:"eventsInQueue,omitempty"`
}

func (m *StatesReportResponse) Reset()                    { *m = StatesReportResponse{} }
func (m *StatesReportResponse) String() string            { return proto.CompactTextString(m) }
func (*StatesReportResponse) ProtoMessage()               {}
func (*StatesReportResponse) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{9} }

func (m *StatesReportResponse) GetStates() []*StateReport {
	if m != nil {
		return m.States
	}
	return nil
}

func (m *StatesReportResponse) GetEventsInQueue() int32 {
	if m != nil {
		return m.EventsInQueue
	}
	return 0
}

type QueryObjectRequest struct {
	VersionID []byte `protobuf:"bytes,1,opt,name=versionID,proto3" json:"versionID,omitempty"`
}

func (m *QueryObjectRequest) Reset()                    { *m = QueryObjectRequest{} }
func (m *QueryObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*QueryObjectRequest) ProtoMessage()               {}
func (*QueryObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{10} }

func (m *QueryObjectRequest) GetVersionID() []byte {
	if m != nil {
		return m.VersionID
	}
	return nil
}

type QueryObjectResponse struct {
	Object *Object `protobuf:"bytes,1,opt,name=object" json:"object,omitempty"`
	Error  string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *QueryObjectResponse) Reset()                    { *m = QueryObjectResponse{} }
func (m *QueryObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*QueryObjectResponse) ProtoMessage()               {}
func (*QueryObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{11} }

func (m *QueryObjectResponse) GetObject() *Object {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *QueryObjectResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type DeleteObjectRequest struct {
	VersionID []byte `protobuf:"bytes,1,opt,name=versionID,proto3" json:"versionID,omitempty"`
}

func (m *DeleteObjectRequest) Reset()                    { *m = DeleteObjectRequest{} }
func (m *DeleteObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteObjectRequest) ProtoMessage()               {}
func (*DeleteObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{12} }

func (m *DeleteObjectRequest) GetVersionID() []byte {
	if m != nil {
		return m.VersionID
	}
	return nil
}

type DeleteObjectResponse struct {
	Object *Object `protobuf:"bytes,1,opt,name=object" json:"object,omitempty"`
	Error  string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *DeleteObjectResponse) Reset()                    { *m = DeleteObjectResponse{} }
func (m *DeleteObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteObjectResponse) ProtoMessage()               {}
func (*DeleteObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{13} }

func (m *DeleteObjectResponse) GetObject() *Object {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *DeleteObjectResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type NewObjectRequest struct {
	Object []byte `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
}

func (m *NewObjectRequest) Reset()                    { *m = NewObjectRequest{} }
func (m *NewObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*NewObjectRequest) ProtoMessage()               {}
func (*NewObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{14} }

func (m *NewObjectRequest) GetObject() []byte {
	if m != nil {
		return m.Object
	}
	return nil
}

type NewObjectResponse struct {
	ID    []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *NewObjectResponse) Reset()                    { *m = NewObjectResponse{} }
func (m *NewObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*NewObjectResponse) ProtoMessage()               {}
func (*NewObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{15} }

func (m *NewObjectResponse) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *NewObjectResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type CreateObjectsRequest struct {
	Objects [][]byte `protobuf:"bytes,1,rep,name=objects" json:"objects,omitempty"`
}

func (m *CreateObjectsRequest) Reset()                    { *m = CreateObjectsRequest{} }
func (m *CreateObjectsRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateObjectsRequest) ProtoMessage()               {}
func (*CreateObjectsRequest) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{16} }

func (m *CreateObjectsRequest) GetObjects() [][]byte {
	if m != nil {
		return m.Objects
	}
	return nil
}

type CreateObjectsResponse struct {
	IDs   [][]byte `protobuf:"bytes,1,rep,name=ids" json:"ids,omitempty"`
	Error string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *CreateObjectsResponse) Reset()                    { *m = CreateObjectsResponse{} }
func (m *CreateObjectsResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateObjectsResponse) ProtoMessage()               {}
func (*CreateObjectsResponse) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{17} }

func (m *CreateObjectsResponse) GetIDs() [][]byte {
	if m != nil {
		return m.IDs
	}
	return nil
}

func (m *CreateObjectsResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type SBACMessage struct {
	Op       SBACOpcode       `protobuf:"varint,1,opt,name=op,proto3,enum=sbac.SBACOpcode" json:"op,omitempty"`
	Decision SBACDecision     `protobuf:"varint,2,opt,name=decision,proto3,enum=sbac.SBACDecision" json:"decision,omitempty"`
	Tx       *SBACTransaction `protobuf:"bytes,3,opt,name=tx" json:"tx,omitempty"`
	PeerID   uint64           `protobuf:"varint,4,opt,name=peerId,proto3" json:"peerId,omitempty"`
}

func (m *SBACMessage) Reset()                    { *m = SBACMessage{} }
func (m *SBACMessage) String() string            { return proto.CompactTextString(m) }
func (*SBACMessage) ProtoMessage()               {}
func (*SBACMessage) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{18} }

func (m *SBACMessage) GetOp() SBACOpcode {
	if m != nil {
		return m.Op
	}
	return SBACOpcode_CONSENSUS1
}

func (m *SBACMessage) GetDecision() SBACDecision {
	if m != nil {
		return m.Decision
	}
	return SBACDecision_ACCEPT
}

func (m *SBACMessage) GetTx() *SBACTransaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *SBACMessage) GetPeerID() uint64 {
	if m != nil {
		return m.PeerID
	}
	return 0
}

type SBACMessageAck struct {
	LastID uint64 `protobuf:"varint,1,opt,name=lastId,proto3" json:"lastId,omitempty"`
}

func (m *SBACMessageAck) Reset()                    { *m = SBACMessageAck{} }
func (m *SBACMessageAck) String() string            { return proto.CompactTextString(m) }
func (*SBACMessageAck) ProtoMessage()               {}
func (*SBACMessageAck) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{19} }

func (m *SBACMessageAck) GetLastID() uint64 {
	if m != nil {
		return m.LastID
	}
	return 0
}

type SBACTransaction struct {
	ID        []byte            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Tx        *Transaction      `protobuf:"bytes,2,opt,name=tx" json:"tx,omitempty"`
	Evidences map[uint64][]byte `protobuf:"bytes,3,rep,name=evidences" json:"evidences,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Op        SBACOpcode        `protobuf:"varint,4,opt,name=op,proto3,enum=sbac.SBACOpcode" json:"op,omitempty"`
	Signature []byte            `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *SBACTransaction) Reset()                    { *m = SBACTransaction{} }
func (m *SBACTransaction) String() string            { return proto.CompactTextString(m) }
func (*SBACTransaction) ProtoMessage()               {}
func (*SBACTransaction) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{20} }

func (m *SBACTransaction) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *SBACTransaction) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *SBACTransaction) GetEvidences() map[uint64][]byte {
	if m != nil {
		return m.Evidences
	}
	return nil
}

func (m *SBACTransaction) GetOp() SBACOpcode {
	if m != nil {
		return m.Op
	}
	return SBACOpcode_CONSENSUS1
}

func (m *SBACTransaction) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type ConsensusTransaction struct {
	TxID      []byte            `protobuf:"bytes,1,opt,name=txId,proto3" json:"txId,omitempty"`
	Tx        *Transaction      `protobuf:"bytes,2,opt,name=tx" json:"tx,omitempty"`
	Evidences map[uint64][]byte `protobuf:"bytes,3,rep,name=evidences" json:"evidences,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Op        ConsensusOp       `protobuf:"varint,4,opt,name=op,proto3,enum=sbac.ConsensusOp" json:"op,omitempty"`
	Initiator uint64            `protobuf:"varint,5,opt,name=initiator,proto3" json:"initiator,omitempty"`
}

func (m *ConsensusTransaction) Reset()                    { *m = ConsensusTransaction{} }
func (m *ConsensusTransaction) String() string            { return proto.CompactTextString(m) }
func (*ConsensusTransaction) ProtoMessage()               {}
func (*ConsensusTransaction) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{21} }

func (m *ConsensusTransaction) GetTxID() []byte {
	if m != nil {
		return m.TxID
	}
	return nil
}

func (m *ConsensusTransaction) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *ConsensusTransaction) GetEvidences() map[uint64][]byte {
	if m != nil {
		return m.Evidences
	}
	return nil
}

func (m *ConsensusTransaction) GetOp() ConsensusOp {
	if m != nil {
		return m.Op
	}
	return ConsensusOp_Consensus1
}

func (m *ConsensusTransaction) GetInitiator() uint64 {
	if m != nil {
		return m.Initiator
	}
	return 0
}

type SBACMessage2 struct {
	Op        SBACOp            `protobuf:"varint,1,opt,name=op,proto3,enum=sbac.SBACOp" json:"op,omitempty"`
	Decision  SBACDecision      `protobuf:"varint,2,opt,name=decision,proto3,enum=sbac.SBACDecision" json:"decision,omitempty"`
	TxID      []byte            `protobuf:"bytes,3,opt,name=txId,proto3" json:"txId,omitempty"`
	Tx        *Transaction      `protobuf:"bytes,4,opt,name=tx" json:"tx,omitempty"`
	Evidences map[uint64][]byte `protobuf:"bytes,5,rep,name=evidences" json:"evidences,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Signature []byte            `protobuf:"bytes,6,opt,name=signature,proto3" json:"signature,omitempty"`
	PeerID    uint64            `protobuf:"varint,7,opt,name=peerId,proto3" json:"peerId,omitempty"`
}

func (m *SBACMessage2) Reset()                    { *m = SBACMessage2{} }
func (m *SBACMessage2) String() string            { return proto.CompactTextString(m) }
func (*SBACMessage2) ProtoMessage()               {}
func (*SBACMessage2) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{22} }

func (m *SBACMessage2) GetOp() SBACOp {
	if m != nil {
		return m.Op
	}
	return SBACOp_Phase1
}

func (m *SBACMessage2) GetDecision() SBACDecision {
	if m != nil {
		return m.Decision
	}
	return SBACDecision_ACCEPT
}

func (m *SBACMessage2) GetTxID() []byte {
	if m != nil {
		return m.TxID
	}
	return nil
}

func (m *SBACMessage2) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *SBACMessage2) GetEvidences() map[uint64][]byte {
	if m != nil {
		return m.Evidences
	}
	return nil
}

func (m *SBACMessage2) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *SBACMessage2) GetPeerID() uint64 {
	if m != nil {
		return m.PeerID
	}
	return 0
}

func init() {
	proto.RegisterType((*Transaction)(nil), "sbac.Transaction")
	proto.RegisterType((*Strings)(nil), "sbac.strings")
	proto.RegisterType((*Trace)(nil), "sbac.Trace")
	proto.RegisterType((*Object)(nil), "sbac.Object")
	proto.RegisterType((*AddTransactionRequest)(nil), "sbac.AddTransactionRequest")
	proto.RegisterType((*ObjectTraceIDPair)(nil), "sbac.ObjectTraceIDPair")
	proto.RegisterType((*ObjectList)(nil), "sbac.ObjectList")
	proto.RegisterType((*AddTransactionResponse)(nil), "sbac.AddTransactionResponse")
	proto.RegisterType((*StateReport)(nil), "sbac.StateReport")
	proto.RegisterType((*StatesReportResponse)(nil), "sbac.StatesReportResponse")
	proto.RegisterType((*QueryObjectRequest)(nil), "sbac.QueryObjectRequest")
	proto.RegisterType((*QueryObjectResponse)(nil), "sbac.QueryObjectResponse")
	proto.RegisterType((*DeleteObjectRequest)(nil), "sbac.DeleteObjectRequest")
	proto.RegisterType((*DeleteObjectResponse)(nil), "sbac.DeleteObjectResponse")
	proto.RegisterType((*NewObjectRequest)(nil), "sbac.NewObjectRequest")
	proto.RegisterType((*NewObjectResponse)(nil), "sbac.NewObjectResponse")
	proto.RegisterType((*CreateObjectsRequest)(nil), "sbac.CreateObjectsRequest")
	proto.RegisterType((*CreateObjectsResponse)(nil), "sbac.CreateObjectsResponse")
	proto.RegisterType((*SBACMessage)(nil), "sbac.SBACMessage")
	proto.RegisterType((*SBACMessageAck)(nil), "sbac.SBACMessageAck")
	proto.RegisterType((*SBACTransaction)(nil), "sbac.SBACTransaction")
	proto.RegisterType((*ConsensusTransaction)(nil), "sbac.ConsensusTransaction")
	proto.RegisterType((*SBACMessage2)(nil), "sbac.SBACMessage2")
	proto.RegisterEnum("sbac.Opcode", Opcode_name, Opcode_value)
	proto.RegisterEnum("sbac.ObjectStatus", ObjectStatus_name, ObjectStatus_value)
	proto.RegisterEnum("sbac.SBACOpcode", SBACOpcode_name, SBACOpcode_value)
	proto.RegisterEnum("sbac.SBACDecision", SBACDecision_name, SBACDecision_value)
	proto.RegisterEnum("sbac.ConsensusOp", ConsensusOp_name, ConsensusOp_value)
	proto.RegisterEnum("sbac.SBACOp", SBACOp_name, SBACOp_value)
}

func init() { proto.RegisterFile("sbac/types.proto", fileDescriptorTypes) }

var fileDescriptorTypes = []byte{
	// 1478 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x57, 0xdd, 0x72, 0xdb, 0xc4,
	0x17, 0x8f, 0x65, 0xd9, 0x8e, 0x8f, 0x9d, 0x44, 0xd9, 0x26, 0x19, 0x35, 0x93, 0xff, 0xdf, 0xa9,
	0x68, 0x33, 0xa9, 0x19, 0x1c, 0xe2, 0xf6, 0x82, 0x61, 0x60, 0x18, 0x7f, 0x68, 0xa8, 0x69, 0x6a,
	0x27, 0x6b, 0xb7, 0xc0, 0x55, 0x47, 0xb1, 0xb7, 0x8e, 0x68, 0x22, 0x19, 0xed, 0x3a, 0x24, 0x2f,
	0xc0, 0x15, 0x03, 0x8f, 0xc0, 0x05, 0x77, 0xf0, 0x16, 0x3c, 0x88, 0x2f, 0xfc, 0x24, 0xcc, 0x7e,
	0xc8, 0x96, 0x5c, 0xdb, 0x40, 0xc3, 0x9d, 0xce, 0xd7, 0xef, 0x9c, 0xfd, 0xed, 0xd9, 0xb3, 0x2b,
	0x30, 0xe8, 0xb9, 0xd3, 0x3d, 0x62, 0xb7, 0x03, 0x42, 0x4b, 0x83, 0xc0, 0x67, 0x3e, 0xd2, 0xb9,
	0x66, 0xf7, 0xa3, 0xbe, 0xcb, 0x2e, 0x86, 0xe7, 0xa5, 0xae, 0x7f, 0x75, 0xd4, 0xf7, 0xfb, 0xfe,
	0x91, 0x30, 0x9e, 0x0f, 0xdf, 0x08, 0x49, 0x08, 0xe2, 0x4b, 0x06, 0x59, 0x65, 0xc8, 0x75, 0x02,
	0xc7, 0xa3, 0x4e, 0x97, 0xb9, 0xbe, 0x87, 0x3e, 0x80, 0x34, 0x0b, 0x9c, 0x2e, 0xa1, 0x66, 0x62,
	0x3f, 0x79, 0x98, 0x2b, 0xe7, 0x4a, 0x1c, 0xb4, 0xd4, 0xe1, 0x3a, 0xac, 0x4c, 0xd6, 0xff, 0x20,
	0x43, 0x59, 0xe0, 0x7a, 0x7d, 0x8a, 0x10, 0xe8, 0x94, 0x05, 0xd2, 0x3b, 0x8b, 0xc5, 0xb7, 0xf5,
	0xa3, 0x0e, 0x29, 0x11, 0x80, 0x4a, 0x00, 0x5d, 0xdf, 0xe3, 0x51, 0xac, 0x51, 0x37, 0x13, 0xfb,
	0x89, 0xc3, 0x6c, 0x75, 0x7d, 0x3c, 0x2a, 0x40, 0x6d, 0xa2, 0xc5, 0x11, 0x0f, 0xb4, 0x07, 0xd9,
	0x41, 0xe0, 0x77, 0x49, 0x6f, 0x18, 0x10, 0x53, 0xe3, 0xee, 0x78, 0xaa, 0x40, 0x2d, 0xd8, 0x76,
	0xbd, 0xc1, 0x90, 0xb5, 0xce, 0xbf, 0x23, 0x5d, 0xf6, 0x8a, 0x04, 0xd4, 0xf5, 0xbd, 0x46, 0x9d,
	0x9a, 0xc9, 0xfd, 0xe4, 0x61, 0xbe, 0x7a, 0x7f, 0x3c, 0x2a, 0x6c, 0x37, 0xe6, 0x39, 0xe0, 0xf9,
	0x71, 0xe8, 0x1b, 0x30, 0x85, 0x01, 0x93, 0x37, 0x24, 0x20, 0x5e, 0x97, 0x44, 0x30, 0x75, 0x81,
	0xb9, 0x37, 0x1e, 0x15, 0xcc, 0xc6, 0x02, 0x1f, 0xbc, 0x30, 0x1a, 0x3d, 0x84, 0x35, 0x7f, 0xc8,
	0x26, 0x39, 0xa9, 0x99, 0xe2, 0x70, 0x38, 0xae, 0x44, 0xff, 0x07, 0x18, 0x38, 0x81, 0x73, 0x45,
	0x18, 0x09, 0xa8, 0x99, 0x16, 0x2e, 0x11, 0x0d, 0x32, 0x21, 0x13, 0x10, 0x36, 0x0c, 0x3c, 0x6a,
	0x66, 0x84, 0x31, 0x14, 0xd1, 0x23, 0x48, 0x5f, 0x3a, 0xe7, 0xe4, 0x92, 0x9a, 0xab, 0x62, 0x9b,
	0xd6, 0xe4, 0x36, 0xa9, 0x5d, 0xc1, 0xca, 0x88, 0x8e, 0x20, 0xdf, 0x23, 0x03, 0xe2, 0xf5, 0x88,
	0xd7, 0x75, 0x09, 0x35, 0xb3, 0xef, 0xee, 0x69, 0xcc, 0x01, 0x59, 0x90, 0x8f, 0x50, 0x45, 0x4d,
	0x10, 0x69, 0x63, 0x3a, 0x74, 0x08, 0x1b, 0xf1, 0x75, 0x53, 0x33, 0x27, 0xdc, 0x66, 0xd5, 0xd6,
	0x2f, 0x09, 0x48, 0xcb, 0x28, 0xf4, 0x21, 0x64, 0xaf, 0x43, 0x7a, 0x44, 0x23, 0xe4, 0xab, 0x6b,
	0xe3, 0x51, 0x21, 0x3b, 0xe1, 0x0c, 0x4f, 0xed, 0x68, 0x0b, 0x52, 0xd7, 0xce, 0xe5, 0x50, 0xb6,
	0x40, 0x1e, 0x4b, 0x01, 0x15, 0x21, 0x4d, 0x99, 0xc3, 0x86, 0x7c, 0xbf, 0x13, 0x87, 0xeb, 0x65,
	0x24, 0x97, 0x21, 0x13, 0xb4, 0x85, 0x05, 0x2b, 0x0f, 0xb4, 0x33, 0xe1, 0x47, 0x17, 0x8d, 0xa9,
	0x24, 0xeb, 0xcf, 0x04, 0x6c, 0x57, 0x7a, 0xbd, 0x48, 0xc7, 0x63, 0xf2, 0xfd, 0x90, 0x50, 0x86,
	0x1e, 0x80, 0xc6, 0x6e, 0x44, 0x65, 0xb9, 0xf2, 0xe6, 0x84, 0xa0, 0x89, 0x97, 0xc6, 0x6e, 0xd0,
	0x33, 0xc8, 0x92, 0x6b, 0xb7, 0x27, 0x97, 0xac, 0x09, 0x2a, 0x8b, 0xd2, 0x73, 0x2e, 0x64, 0xc9,
	0x0e, 0x9d, 0x6d, 0x8f, 0x05, 0xb7, 0x78, 0x1a, 0xbc, 0xfb, 0x19, 0xac, 0xc7, 0x8d, 0xc8, 0x80,
	0xe4, 0x5b, 0x72, 0x2b, 0xf2, 0xeb, 0x98, 0x7f, 0xce, 0x27, 0xe1, 0x53, 0xed, 0x93, 0x84, 0x75,
	0x0e, 0x9b, 0x72, 0xd1, 0x62, 0x07, 0x1b, 0xf5, 0x53, 0xc7, 0x0d, 0xd0, 0x23, 0xc8, 0x30, 0x29,
	0x2a, 0x7a, 0x73, 0xe3, 0x51, 0x21, 0xa3, 0x3c, 0x70, 0x68, 0x43, 0x07, 0x90, 0xf1, 0xd5, 0xde,
	0xca, 0x15, 0xe4, 0xa3, 0x2c, 0xe2, 0xd0, 0x68, 0x95, 0x00, 0xa4, 0xea, 0xc4, 0xa5, 0x0c, 0xed,
	0x83, 0x7e, 0xe9, 0x52, 0xa6, 0x66, 0x42, 0x3c, 0x44, 0x58, 0xac, 0x3f, 0x12, 0xb0, 0x33, 0xcb,
	0x02, 0x1d, 0xf8, 0x1e, 0x25, 0xa8, 0x36, 0x4d, 0x29, 0xe3, 0x1f, 0xcf, 0x27, 0x4d, 0xba, 0x2b,
	0x58, 0xc5, 0x59, 0x18, 0xb9, 0x7b, 0x02, 0xf9, 0xa8, 0x21, 0xca, 0x57, 0x56, 0xf2, 0x75, 0x10,
	0xe5, 0x2b, 0x57, 0x36, 0xa2, 0x45, 0xf2, 0x45, 0x44, 0x19, 0xfc, 0x4d, 0x87, 0x1c, 0xef, 0x18,
	0x82, 0xc9, 0xc0, 0x0f, 0x18, 0x6f, 0x97, 0x0b, 0x87, 0x5e, 0x28, 0xee, 0xd6, 0xb0, 0x92, 0xd0,
	0x29, 0x6c, 0x74, 0xfd, 0xab, 0x2b, 0x97, 0xd5, 0x49, 0xd7, 0xe5, 0xcd, 0x19, 0xb2, 0x76, 0x20,
	0xd1, 0x23, 0x18, 0xa5, 0x5a, 0xdc, 0x51, 0xd6, 0x3f, 0x1b, 0xce, 0x11, 0x07, 0x17, 0x0e, 0x25,
	0xc7, 0x53, 0xc4, 0xe4, 0x22, 0xc4, 0xd3, 0xb8, 0xa3, 0x42, 0x9c, 0x09, 0x9f, 0x20, 0x96, 0xa7,
	0x88, 0xfa, 0x52, 0xc4, 0xf2, 0x5c, 0xc4, 0xa9, 0x96, 0x77, 0x1e, 0x3f, 0x46, 0xc4, 0x4c, 0x09,
	0x76, 0xa5, 0xc0, 0x47, 0x1a, 0x1f, 0x14, 0xae, 0xd7, 0xb7, 0xaf, 0x89, 0xc7, 0xf8, 0xbc, 0x4a,
	0x1c, 0xa6, 0x70, 0x5c, 0xb9, 0x5b, 0x85, 0xad, 0x79, 0x44, 0xfc, 0x5d, 0x7f, 0xaf, 0x46, 0x76,
	0x87, 0x63, 0xcc, 0x5b, 0xfa, 0x7b, 0x61, 0x94, 0xdf, 0x1f, 0xc3, 0xea, 0xc3, 0x96, 0x20, 0x8f,
	0x4a, 0xf6, 0x26, 0x0d, 0xfd, 0x58, 0x0e, 0xa2, 0xc9, 0x1d, 0xb9, 0xf9, 0x0e, 0xd1, 0x58, 0x39,
	0x70, 0xd2, 0x88, 0x20, 0xa6, 0xe1, 0x9d, 0x0d, 0x89, 0x4a, 0x92, 0xc2, 0x71, 0xa5, 0x55, 0x01,
	0x74, 0x36, 0x24, 0xc1, 0xad, 0x3a, 0x51, 0x6a, 0x22, 0xfd, 0x9b, 0x91, 0x69, 0x9d, 0xc1, 0xbd,
	0x18, 0x84, 0x2a, 0xf5, 0x21, 0xa4, 0xe5, 0x09, 0x52, 0x93, 0x2d, 0x7e, 0x74, 0x95, 0x8d, 0x53,
	0x40, 0x82, 0xc0, 0x0f, 0xd4, 0x95, 0x2b, 0x05, 0xab, 0x0a, 0xf7, 0xea, 0xe4, 0x92, 0x30, 0x72,
	0x87, 0xb2, 0x30, 0x6c, 0xc5, 0x31, 0xfe, 0x83, 0xba, 0x8a, 0x60, 0x34, 0xc9, 0x0f, 0xf1, 0xa2,
	0x76, 0x62, 0x78, 0xf9, 0x10, 0xc1, 0xaa, 0xc0, 0x66, 0xc4, 0x57, 0x25, 0xdf, 0x01, 0xcd, 0xed,
	0xa9, 0xd2, 0xd3, 0xe3, 0x51, 0x41, 0x6b, 0xd4, 0xb1, 0xe6, 0xf6, 0x16, 0xa4, 0xfb, 0x18, 0xb6,
	0x6a, 0x01, 0x71, 0xc2, 0x25, 0xd0, 0x30, 0xa5, 0x19, 0x1f, 0x6b, 0xf9, 0xe9, 0xec, 0x7c, 0x06,
	0xdb, 0x33, 0x11, 0x2a, 0xf1, 0x7d, 0x48, 0xba, 0x3d, 0xe5, 0x5e, 0xcd, 0x8c, 0x47, 0x85, 0x24,
	0x7f, 0x45, 0x70, 0xdd, 0x82, 0xdc, 0xbf, 0x27, 0x20, 0xd7, 0xae, 0x56, 0x6a, 0x2f, 0x08, 0xa5,
	0x4e, 0x9f, 0xa0, 0x7d, 0xd0, 0xfc, 0x81, 0xa8, 0x7c, 0x3d, 0x1c, 0x70, 0xdc, 0xdc, 0x1a, 0x74,
	0xfd, 0x1e, 0xc1, 0x9a, 0x3f, 0x40, 0x25, 0x58, 0xed, 0xa9, 0x8e, 0x17, 0x50, 0x93, 0x6b, 0x92,
	0xfb, 0x85, 0x67, 0x01, 0x4f, 0x7c, 0xd0, 0x23, 0x71, 0xed, 0x25, 0xc5, 0x26, 0x6c, 0x4f, 0x3d,
	0x67, 0xaf, 0x3e, 0x0b, 0xd2, 0x03, 0x42, 0x82, 0x46, 0xcf, 0xd4, 0xf9, 0xc9, 0xa9, 0xc2, 0x78,
	0x54, 0x48, 0x9f, 0x72, 0x4d, 0x1d, 0x2b, 0x8b, 0xf5, 0x14, 0xd6, 0x23, 0xb5, 0x56, 0xba, 0x6f,
	0x79, 0xd4, 0xa5, 0x43, 0x59, 0x43, 0x92, 0xad, 0xa2, 0x4e, 0xb8, 0xa6, 0x8e, 0x95, 0xc5, 0xfa,
	0x59, 0x83, 0x8d, 0x99, 0x8c, 0x0b, 0x37, 0x48, 0xde, 0xd1, 0xda, 0xb2, 0x3b, 0xba, 0x1a, 0xbd,
	0xa3, 0xe5, 0x64, 0x7d, 0x38, 0x77, 0x59, 0x8b, 0x6f, 0x67, 0xc5, 0xb2, 0xbe, 0x84, 0xe5, 0x3d,
	0xc8, 0x52, 0xb7, 0xef, 0x39, 0x8c, 0xbf, 0x53, 0x53, 0xa2, 0xe3, 0xa6, 0x8a, 0x3b, 0xde, 0xee,
	0xbf, 0x6a, 0x7c, 0x84, 0x7a, 0x94, 0x78, 0x74, 0x48, 0xa3, 0xac, 0xec, 0x81, 0xce, 0x6e, 0x1a,
	0x21, 0x2f, 0xab, 0xe3, 0x51, 0x41, 0xef, 0xdc, 0x34, 0xea, 0x58, 0x68, 0xff, 0x09, 0x37, 0x5f,
	0xbe, 0xcb, 0x8d, 0xba, 0x8a, 0xe7, 0xe5, 0x5b, 0x42, 0xd0, 0x83, 0x08, 0x41, 0x9b, 0x33, 0x08,
	0xad, 0x41, 0xc8, 0x90, 0xeb, 0xb9, 0xcc, 0x75, 0x98, 0x1f, 0x08, 0x86, 0x74, 0x3c, 0x55, 0xdc,
	0x91, 0xa1, 0x91, 0x06, 0xf9, 0x48, 0xa7, 0x95, 0xd1, 0x5e, 0xe4, 0x58, 0xe4, 0xa3, 0x1b, 0xf6,
	0x5e, 0x47, 0x22, 0xe4, 0x39, 0xb9, 0x84, 0x67, 0x7d, 0x19, 0xcf, 0x5f, 0x44, 0x79, 0x4e, 0x09,
	0x9e, 0x1f, 0x4c, 0x33, 0x86, 0x55, 0x2f, 0xe1, 0x37, 0xd6, 0x5e, 0xe9, 0x99, 0xf6, 0x8a, 0x9c,
	0xc5, 0xcc, 0xa2, 0xb3, 0x78, 0x37, 0x82, 0x8b, 0x3f, 0xf1, 0x77, 0xbb, 0xe8, 0x76, 0x94, 0x83,
	0xcc, 0xcb, 0xe6, 0xf3, 0x66, 0xeb, 0xeb, 0xa6, 0xb1, 0x82, 0xee, 0xc1, 0x46, 0xa5, 0x5e, 0x7f,
	0xdd, 0xc1, 0x95, 0x66, 0xbb, 0x52, 0xeb, 0x34, 0x5a, 0x4d, 0x23, 0x81, 0x0c, 0xc8, 0x9f, 0xbd,
	0xb4, 0xf1, 0xb7, 0xaf, 0x5b, 0xd5, 0xaf, 0xec, 0x5a, 0xc7, 0xd0, 0xd0, 0x26, 0xac, 0xd5, 0xb0,
	0x5d, 0xe9, 0xd8, 0xa1, 0x2a, 0xc9, 0x55, 0x75, 0xfb, 0xc4, 0x9e, 0xaa, 0x74, 0x04, 0x90, 0x6e,
	0x77, 0x2a, 0x1d, 0xbb, 0x6d, 0xa4, 0xd0, 0x2a, 0xe8, 0x9c, 0x1a, 0x23, 0x8d, 0x10, 0xac, 0xc7,
	0x62, 0xdb, 0x46, 0xa6, 0xf8, 0x34, 0x7c, 0xfb, 0xc9, 0x47, 0x3e, 0x8f, 0xe4, 0xd9, 0x5f, 0xd9,
	0xc6, 0x0a, 0xca, 0xc3, 0x6a, 0xa3, 0xa9, 0xa4, 0x04, 0xb7, 0x9c, 0xb4, 0x6a, 0xcf, 0xed, 0xba,
	0xa1, 0x15, 0xdf, 0x00, 0x4c, 0x4f, 0x2d, 0x5a, 0x07, 0xa8, 0xb5, 0x9a, 0x6d, 0xbb, 0xd9, 0x7e,
	0xd9, 0x3e, 0x36, 0x56, 0x62, 0x72, 0x59, 0x46, 0x9e, 0x3e, 0xab, 0xb4, 0xed, 0x63, 0x43, 0x9b,
	0x7c, 0x97, 0x8d, 0x24, 0xff, 0xae, 0xb5, 0x5e, 0xbc, 0x68, 0xf0, 0x8a, 0xb7, 0xc0, 0x98, 0xc4,
	0xbc, 0x56, 0xda, 0x54, 0xf1, 0x40, 0x36, 0x63, 0xd8, 0x48, 0xb2, 0xba, 0x9a, 0x7d, 0xda, 0x31,
	0x56, 0xf8, 0x37, 0xb6, 0xc5, 0x7a, 0x13, 0xc5, 0xcf, 0x21, 0x17, 0x39, 0x24, 0xa2, 0x80, 0x50,
	0x0c, 0x0b, 0x0a, 0x65, 0x5e, 0x50, 0x54, 0x7e, 0x62, 0x68, 0xc5, 0x63, 0x48, 0xcb, 0xe5, 0x88,
	0xf2, 0xc4, 0xf3, 0x48, 0x26, 0x90, 0xcf, 0x1c, 0x23, 0x81, 0x36, 0x20, 0x27, 0xbe, 0xe5, 0xfb,
	0xcb, 0xd0, 0xce, 0xd3, 0xe2, 0x0f, 0xff, 0xc9, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x55, 0xf5,
	0x4f, 0x29, 0x2a, 0x10, 0x00, 0x00,
}
