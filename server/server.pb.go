// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: server.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// AppendEntry
//Invoked by leader to replicate log entries (§5.3); also used as heartbeat (§5.2).
type EntriesArguments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//leader’s term
	Term uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	//so follower can redirect clients
	LeaderID string `protobuf:"bytes,2,opt,name=leaderID,proto3" json:"leaderID,omitempty"`
	//index of log entry immediately preceding new ones
	PrevLogIndex uint64 `protobuf:"varint,3,opt,name=prevLogIndex,proto3" json:"prevLogIndex,omitempty"`
	//term of prevLogIndex entry
	PrevLogTerm uint64 `protobuf:"varint,4,opt,name=prevLogTerm,proto3" json:"prevLogTerm,omitempty"`
	//Entries is a partment of Log
	Entries []*Entry `protobuf:"bytes,5,rep,name=Entries,proto3" json:"Entries,omitempty"`
	//leader’s commitIndex
	LeaderCommit uint64 `protobuf:"varint,6,opt,name=LeaderCommit,proto3" json:"LeaderCommit,omitempty"`
}

func (x *EntriesArguments) Reset() {
	*x = EntriesArguments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntriesArguments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntriesArguments) ProtoMessage() {}

func (x *EntriesArguments) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntriesArguments.ProtoReflect.Descriptor instead.
func (*EntriesArguments) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{0}
}

func (x *EntriesArguments) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *EntriesArguments) GetLeaderID() string {
	if x != nil {
		return x.LeaderID
	}
	return ""
}

func (x *EntriesArguments) GetPrevLogIndex() uint64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *EntriesArguments) GetPrevLogTerm() uint64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *EntriesArguments) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *EntriesArguments) GetLeaderCommit() uint64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

type EntriesResults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//currentTerm, for leader to update itself
	Term uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	//true if follower contained entry matching prevLogIndex and prevLogTerm
	Success bool `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *EntriesResults) Reset() {
	*x = EntriesResults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntriesResults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntriesResults) ProtoMessage() {}

func (x *EntriesResults) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntriesResults.ProtoReflect.Descriptor instead.
func (*EntriesResults) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{1}
}

func (x *EntriesResults) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *EntriesResults) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type VoteArguments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//candidate’s term
	Term uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	//candidate requesting vote
	CandidateID string `protobuf:"bytes,2,opt,name=candidateID,proto3" json:"candidateID,omitempty"`
	//index of candidate’s last log entry (§5.4)
	LastLogIndex uint64 `protobuf:"varint,3,opt,name=lastLogIndex,proto3" json:"lastLogIndex,omitempty"`
	//term of candidate’s last log entry (§5.4)
	LastLogTerm uint64 `protobuf:"varint,4,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
}

func (x *VoteArguments) Reset() {
	*x = VoteArguments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteArguments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteArguments) ProtoMessage() {}

func (x *VoteArguments) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteArguments.ProtoReflect.Descriptor instead.
func (*VoteArguments) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{2}
}

func (x *VoteArguments) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteArguments) GetCandidateID() string {
	if x != nil {
		return x.CandidateID
	}
	return ""
}

func (x *VoteArguments) GetLastLogIndex() uint64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *VoteArguments) GetLastLogTerm() uint64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type VoteResults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//currentTerm, for candidate to update itself
	Term uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	//true means candidate received vote
	VoteGranted bool `protobuf:"varint,2,opt,name=voteGranted,proto3" json:"voteGranted,omitempty"`
}

func (x *VoteResults) Reset() {
	*x = VoteResults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResults) ProtoMessage() {}

func (x *VoteResults) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResults.ProtoReflect.Descriptor instead.
func (*VoteResults) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{3}
}

func (x *VoteResults) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteResults) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term uint64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{4}
}

func (x *Entry) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Entry) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_server_proto protoreflect.FileDescriptor

var file_server_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0xd6, 0x01, 0x0a, 0x10, 0x65, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d,
	0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c,
	0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65,
	0x72, 0x6d, 0x12, 0x28, 0x0a, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x0c,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0c, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x22, 0x3e, 0x0a, 0x0e, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x8b, 0x01, 0x0a, 0x0d, 0x76, 0x6f, 0x74, 0x65, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x6e,
	0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74,
	0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x43,
	0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72,
	0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e,
	0x74, 0x65, 0x64, 0x22, 0x2f, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d,
	0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x32, 0x8e, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12,
	0x45, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x19, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x65, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x17, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x16, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x76, 0x6f, 0x74, 0x65, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x14, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x72, 0x70, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_proto_rawDescOnce sync.Once
	file_server_proto_rawDescData = file_server_proto_rawDesc
)

func file_server_proto_rawDescGZIP() []byte {
	file_server_proto_rawDescOnce.Do(func() {
		file_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_proto_rawDescData)
	})
	return file_server_proto_rawDescData
}

var file_server_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_server_proto_goTypes = []interface{}{
	(*EntriesArguments)(nil), // 0: service.entriesArguments
	(*EntriesResults)(nil),   // 1: service.entriesResults
	(*VoteArguments)(nil),    // 2: service.voteArguments
	(*VoteResults)(nil),      // 3: service.voteResults
	(*Entry)(nil),            // 4: service.Entry
}
var file_server_proto_depIdxs = []int32{
	4, // 0: service.entriesArguments.Entries:type_name -> service.Entry
	0, // 1: service.Server.AppendEntries:input_type -> service.entriesArguments
	2, // 2: service.Server.RequestVote:input_type -> service.voteArguments
	1, // 3: service.Server.AppendEntries:output_type -> service.entriesResults
	3, // 4: service.Server.RequestVote:output_type -> service.voteResults
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_server_proto_init() }
func file_server_proto_init() {
	if File_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntriesArguments); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntriesResults); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteArguments); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResults); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_proto_goTypes,
		DependencyIndexes: file_server_proto_depIdxs,
		MessageInfos:      file_server_proto_msgTypes,
	}.Build()
	File_server_proto = out.File
	file_server_proto_rawDesc = nil
	file_server_proto_goTypes = nil
	file_server_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerClient interface {
	AppendEntries(ctx context.Context, in *EntriesArguments, opts ...grpc.CallOption) (*EntriesResults, error)
	RequestVote(ctx context.Context, in *VoteArguments, opts ...grpc.CallOption) (*VoteResults, error)
}

type serverClient struct {
	cc grpc.ClientConnInterface
}

func NewServerClient(cc grpc.ClientConnInterface) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) AppendEntries(ctx context.Context, in *EntriesArguments, opts ...grpc.CallOption) (*EntriesResults, error) {
	out := new(EntriesResults)
	err := c.cc.Invoke(ctx, "/service.Server/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) RequestVote(ctx context.Context, in *VoteArguments, opts ...grpc.CallOption) (*VoteResults, error) {
	out := new(VoteResults)
	err := c.cc.Invoke(ctx, "/service.Server/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
type ServerServer interface {
	AppendEntries(context.Context, *EntriesArguments) (*EntriesResults, error)
	RequestVote(context.Context, *VoteArguments) (*VoteResults, error)
}

// UnimplementedServerServer can be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (*UnimplementedServerServer) AppendEntries(context.Context, *EntriesArguments) (*EntriesResults, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (*UnimplementedServerServer) RequestVote(context.Context, *VoteArguments) (*VoteResults, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}

func RegisterServerServer(s *grpc.Server, srv ServerServer) {
	s.RegisterService(&_Server_serviceDesc, srv)
}

func _Server_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntriesArguments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Server/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).AppendEntries(ctx, req.(*EntriesArguments))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteArguments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Server/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).RequestVote(ctx, req.(*VoteArguments))
	}
	return interceptor(ctx, in, info, handler)
}

var _Server_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Server_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Server_RequestVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
