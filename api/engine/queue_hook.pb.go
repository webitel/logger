// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: queue_hook.proto

package engine

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DeleteQueueHookRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       uint32                 `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Id            uint32                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteQueueHookRequest) Reset() {
	*x = DeleteQueueHookRequest{}
	mi := &file_queue_hook_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteQueueHookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteQueueHookRequest) ProtoMessage() {}

func (x *DeleteQueueHookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteQueueHookRequest.ProtoReflect.Descriptor instead.
func (*DeleteQueueHookRequest) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteQueueHookRequest) GetQueueId() uint32 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *DeleteQueueHookRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type PatchQueueHookRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Fields        []string               `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	QueueId       uint32                 `protobuf:"varint,2,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Id            uint32                 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	Schema        *Lookup                `protobuf:"bytes,4,opt,name=schema,proto3" json:"schema,omitempty"`
	Event         string                 `protobuf:"bytes,5,opt,name=event,proto3" json:"event,omitempty"`
	Enabled       bool                   `protobuf:"varint,6,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Properties    []string               `protobuf:"bytes,7,rep,name=properties,proto3" json:"properties,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchQueueHookRequest) Reset() {
	*x = PatchQueueHookRequest{}
	mi := &file_queue_hook_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchQueueHookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchQueueHookRequest) ProtoMessage() {}

func (x *PatchQueueHookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchQueueHookRequest.ProtoReflect.Descriptor instead.
func (*PatchQueueHookRequest) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{1}
}

func (x *PatchQueueHookRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *PatchQueueHookRequest) GetQueueId() uint32 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *PatchQueueHookRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PatchQueueHookRequest) GetSchema() *Lookup {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *PatchQueueHookRequest) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *PatchQueueHookRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *PatchQueueHookRequest) GetProperties() []string {
	if x != nil {
		return x.Properties
	}
	return nil
}

type UpdateQueueHookRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       uint32                 `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Id            uint32                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Schema        *Lookup                `protobuf:"bytes,3,opt,name=schema,proto3" json:"schema,omitempty"`
	Event         string                 `protobuf:"bytes,4,opt,name=event,proto3" json:"event,omitempty"`
	Enabled       bool                   `protobuf:"varint,5,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Properties    []string               `protobuf:"bytes,6,rep,name=properties,proto3" json:"properties,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateQueueHookRequest) Reset() {
	*x = UpdateQueueHookRequest{}
	mi := &file_queue_hook_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateQueueHookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateQueueHookRequest) ProtoMessage() {}

func (x *UpdateQueueHookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateQueueHookRequest.ProtoReflect.Descriptor instead.
func (*UpdateQueueHookRequest) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateQueueHookRequest) GetQueueId() uint32 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *UpdateQueueHookRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateQueueHookRequest) GetSchema() *Lookup {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *UpdateQueueHookRequest) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *UpdateQueueHookRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *UpdateQueueHookRequest) GetProperties() []string {
	if x != nil {
		return x.Properties
	}
	return nil
}

type SearchQueueHookRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       uint32                 `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Q             string                 `protobuf:"bytes,4,opt,name=q,proto3" json:"q,omitempty"`
	Sort          string                 `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields        []string               `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            []uint32               `protobuf:"varint,7,rep,packed,name=id,proto3" json:"id,omitempty"`
	SchemaId      []uint32               `protobuf:"varint,8,rep,packed,name=schema_id,json=schemaId,proto3" json:"schema_id,omitempty"`
	Event         []string               `protobuf:"bytes,9,rep,name=event,proto3" json:"event,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchQueueHookRequest) Reset() {
	*x = SearchQueueHookRequest{}
	mi := &file_queue_hook_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchQueueHookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchQueueHookRequest) ProtoMessage() {}

func (x *SearchQueueHookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchQueueHookRequest.ProtoReflect.Descriptor instead.
func (*SearchQueueHookRequest) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{3}
}

func (x *SearchQueueHookRequest) GetQueueId() uint32 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *SearchQueueHookRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchQueueHookRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchQueueHookRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchQueueHookRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchQueueHookRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchQueueHookRequest) GetId() []uint32 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SearchQueueHookRequest) GetSchemaId() []uint32 {
	if x != nil {
		return x.SchemaId
	}
	return nil
}

func (x *SearchQueueHookRequest) GetEvent() []string {
	if x != nil {
		return x.Event
	}
	return nil
}

type CreateQueueHookRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       uint32                 `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Schema        *Lookup                `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	Event         string                 `protobuf:"bytes,3,opt,name=event,proto3" json:"event,omitempty"`
	Enabled       bool                   `protobuf:"varint,4,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Properties    []string               `protobuf:"bytes,5,rep,name=properties,proto3" json:"properties,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateQueueHookRequest) Reset() {
	*x = CreateQueueHookRequest{}
	mi := &file_queue_hook_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateQueueHookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateQueueHookRequest) ProtoMessage() {}

func (x *CreateQueueHookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateQueueHookRequest.ProtoReflect.Descriptor instead.
func (*CreateQueueHookRequest) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{4}
}

func (x *CreateQueueHookRequest) GetQueueId() uint32 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *CreateQueueHookRequest) GetSchema() *Lookup {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *CreateQueueHookRequest) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *CreateQueueHookRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *CreateQueueHookRequest) GetProperties() []string {
	if x != nil {
		return x.Properties
	}
	return nil
}

type ReadQueueHookRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       uint32                 `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Id            uint32                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadQueueHookRequest) Reset() {
	*x = ReadQueueHookRequest{}
	mi := &file_queue_hook_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadQueueHookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadQueueHookRequest) ProtoMessage() {}

func (x *ReadQueueHookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadQueueHookRequest.ProtoReflect.Descriptor instead.
func (*ReadQueueHookRequest) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{5}
}

func (x *ReadQueueHookRequest) GetQueueId() uint32 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *ReadQueueHookRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type QueueHook struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Schema        *Lookup                `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	Event         string                 `protobuf:"bytes,3,opt,name=event,proto3" json:"event,omitempty"`
	Enabled       bool                   `protobuf:"varint,4,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Properties    []string               `protobuf:"bytes,5,rep,name=properties,proto3" json:"properties,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueueHook) Reset() {
	*x = QueueHook{}
	mi := &file_queue_hook_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueHook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueHook) ProtoMessage() {}

func (x *QueueHook) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueHook.ProtoReflect.Descriptor instead.
func (*QueueHook) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{6}
}

func (x *QueueHook) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *QueueHook) GetSchema() *Lookup {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *QueueHook) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *QueueHook) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *QueueHook) GetProperties() []string {
	if x != nil {
		return x.Properties
	}
	return nil
}

type ListQueueHook struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*QueueHook           `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListQueueHook) Reset() {
	*x = ListQueueHook{}
	mi := &file_queue_hook_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListQueueHook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueueHook) ProtoMessage() {}

func (x *ListQueueHook) ProtoReflect() protoreflect.Message {
	mi := &file_queue_hook_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueueHook.ProtoReflect.Descriptor instead.
func (*ListQueueHook) Descriptor() ([]byte, []int) {
	return file_queue_hook_proto_rawDescGZIP(), []int{7}
}

func (x *ListQueueHook) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListQueueHook) GetItems() []*QueueHook {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_queue_hook_proto protoreflect.FileDescriptor

const file_queue_hook_proto_rawDesc = "" +
	"\n" +
	"\x10queue_hook.proto\x12\x06engine\x1a\x1cgoogle/api/annotations.proto\x1a\vconst.proto\"C\n" +
	"\x16DeleteQueueHookRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\rR\aqueueId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\rR\x02id\"\xd2\x01\n" +
	"\x15PatchQueueHookRequest\x12\x16\n" +
	"\x06fields\x18\x01 \x03(\tR\x06fields\x12\x19\n" +
	"\bqueue_id\x18\x02 \x01(\rR\aqueueId\x12\x0e\n" +
	"\x02id\x18\x03 \x01(\rR\x02id\x12&\n" +
	"\x06schema\x18\x04 \x01(\v2\x0e.engine.LookupR\x06schema\x12\x14\n" +
	"\x05event\x18\x05 \x01(\tR\x05event\x12\x18\n" +
	"\aenabled\x18\x06 \x01(\bR\aenabled\x12\x1e\n" +
	"\n" +
	"properties\x18\a \x03(\tR\n" +
	"properties\"\xbb\x01\n" +
	"\x16UpdateQueueHookRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\rR\aqueueId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\rR\x02id\x12&\n" +
	"\x06schema\x18\x03 \x01(\v2\x0e.engine.LookupR\x06schema\x12\x14\n" +
	"\x05event\x18\x04 \x01(\tR\x05event\x12\x18\n" +
	"\aenabled\x18\x05 \x01(\bR\aenabled\x12\x1e\n" +
	"\n" +
	"properties\x18\x06 \x03(\tR\n" +
	"properties\"\xd8\x01\n" +
	"\x16SearchQueueHookRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\rR\aqueueId\x12\x12\n" +
	"\x04page\x18\x02 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x03 \x01(\x05R\x04size\x12\f\n" +
	"\x01q\x18\x04 \x01(\tR\x01q\x12\x12\n" +
	"\x04sort\x18\x05 \x01(\tR\x04sort\x12\x16\n" +
	"\x06fields\x18\x06 \x03(\tR\x06fields\x12\x0e\n" +
	"\x02id\x18\a \x03(\rR\x02id\x12\x1b\n" +
	"\tschema_id\x18\b \x03(\rR\bschemaId\x12\x14\n" +
	"\x05event\x18\t \x03(\tR\x05event\"\xab\x01\n" +
	"\x16CreateQueueHookRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\rR\aqueueId\x12&\n" +
	"\x06schema\x18\x02 \x01(\v2\x0e.engine.LookupR\x06schema\x12\x14\n" +
	"\x05event\x18\x03 \x01(\tR\x05event\x12\x18\n" +
	"\aenabled\x18\x04 \x01(\bR\aenabled\x12\x1e\n" +
	"\n" +
	"properties\x18\x05 \x03(\tR\n" +
	"properties\"A\n" +
	"\x14ReadQueueHookRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\rR\aqueueId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\rR\x02id\"\x93\x01\n" +
	"\tQueueHook\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\rR\x02id\x12&\n" +
	"\x06schema\x18\x02 \x01(\v2\x0e.engine.LookupR\x06schema\x12\x14\n" +
	"\x05event\x18\x03 \x01(\tR\x05event\x12\x18\n" +
	"\aenabled\x18\x04 \x01(\bR\aenabled\x12\x1e\n" +
	"\n" +
	"properties\x18\x05 \x03(\tR\n" +
	"properties\"L\n" +
	"\rListQueueHook\x12\x12\n" +
	"\x04next\x18\x01 \x01(\bR\x04next\x12'\n" +
	"\x05items\x18\x02 \x03(\v2\x11.engine.QueueHookR\x05items2\xe5\x05\n" +
	"\x10QueueHookService\x12u\n" +
	"\x0fCreateQueueHook\x12\x1e.engine.CreateQueueHookRequest\x1a\x11.engine.QueueHook\"/\x82\xd3\xe4\x93\x02):\x01*\"$/call_center/queues/{queue_id}/hooks\x12v\n" +
	"\x0fSearchQueueHook\x12\x1e.engine.SearchQueueHookRequest\x1a\x15.engine.ListQueueHook\",\x82\xd3\xe4\x93\x02&\x12$/call_center/queues/{queue_id}/hooks\x12s\n" +
	"\rReadQueueHook\x12\x1c.engine.ReadQueueHookRequest\x1a\x11.engine.QueueHook\"1\x82\xd3\xe4\x93\x02+\x12)/call_center/queues/{queue_id}/hooks/{id}\x12z\n" +
	"\x0fUpdateQueueHook\x12\x1e.engine.UpdateQueueHookRequest\x1a\x11.engine.QueueHook\"4\x82\xd3\xe4\x93\x02.:\x01*\x1a)/call_center/queues/{queue_id}/hooks/{id}\x12x\n" +
	"\x0ePatchQueueHook\x12\x1d.engine.PatchQueueHookRequest\x1a\x11.engine.QueueHook\"4\x82\xd3\xe4\x93\x02.:\x01*2)/call_center/queues/{queue_id}/hooks/{id}\x12w\n" +
	"\x0fDeleteQueueHook\x12\x1e.engine.DeleteQueueHookRequest\x1a\x11.engine.QueueHook\"1\x82\xd3\xe4\x93\x02+*)/call_center/queues/{queue_id}/hooks/{id}B\"Z github.com/webitel/protos/engineb\x06proto3"

var (
	file_queue_hook_proto_rawDescOnce sync.Once
	file_queue_hook_proto_rawDescData []byte
)

func file_queue_hook_proto_rawDescGZIP() []byte {
	file_queue_hook_proto_rawDescOnce.Do(func() {
		file_queue_hook_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_queue_hook_proto_rawDesc), len(file_queue_hook_proto_rawDesc)))
	})
	return file_queue_hook_proto_rawDescData
}

var file_queue_hook_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_queue_hook_proto_goTypes = []any{
	(*DeleteQueueHookRequest)(nil), // 0: engine.DeleteQueueHookRequest
	(*PatchQueueHookRequest)(nil),  // 1: engine.PatchQueueHookRequest
	(*UpdateQueueHookRequest)(nil), // 2: engine.UpdateQueueHookRequest
	(*SearchQueueHookRequest)(nil), // 3: engine.SearchQueueHookRequest
	(*CreateQueueHookRequest)(nil), // 4: engine.CreateQueueHookRequest
	(*ReadQueueHookRequest)(nil),   // 5: engine.ReadQueueHookRequest
	(*QueueHook)(nil),              // 6: engine.QueueHook
	(*ListQueueHook)(nil),          // 7: engine.ListQueueHook
	(*Lookup)(nil),                 // 8: engine.Lookup
}
var file_queue_hook_proto_depIdxs = []int32{
	8,  // 0: engine.PatchQueueHookRequest.schema:type_name -> engine.Lookup
	8,  // 1: engine.UpdateQueueHookRequest.schema:type_name -> engine.Lookup
	8,  // 2: engine.CreateQueueHookRequest.schema:type_name -> engine.Lookup
	8,  // 3: engine.QueueHook.schema:type_name -> engine.Lookup
	6,  // 4: engine.ListQueueHook.items:type_name -> engine.QueueHook
	4,  // 5: engine.QueueHookService.CreateQueueHook:input_type -> engine.CreateQueueHookRequest
	3,  // 6: engine.QueueHookService.SearchQueueHook:input_type -> engine.SearchQueueHookRequest
	5,  // 7: engine.QueueHookService.ReadQueueHook:input_type -> engine.ReadQueueHookRequest
	2,  // 8: engine.QueueHookService.UpdateQueueHook:input_type -> engine.UpdateQueueHookRequest
	1,  // 9: engine.QueueHookService.PatchQueueHook:input_type -> engine.PatchQueueHookRequest
	0,  // 10: engine.QueueHookService.DeleteQueueHook:input_type -> engine.DeleteQueueHookRequest
	6,  // 11: engine.QueueHookService.CreateQueueHook:output_type -> engine.QueueHook
	7,  // 12: engine.QueueHookService.SearchQueueHook:output_type -> engine.ListQueueHook
	6,  // 13: engine.QueueHookService.ReadQueueHook:output_type -> engine.QueueHook
	6,  // 14: engine.QueueHookService.UpdateQueueHook:output_type -> engine.QueueHook
	6,  // 15: engine.QueueHookService.PatchQueueHook:output_type -> engine.QueueHook
	6,  // 16: engine.QueueHookService.DeleteQueueHook:output_type -> engine.QueueHook
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_queue_hook_proto_init() }
func file_queue_hook_proto_init() {
	if File_queue_hook_proto != nil {
		return
	}
	file_const_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_queue_hook_proto_rawDesc), len(file_queue_hook_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_queue_hook_proto_goTypes,
		DependencyIndexes: file_queue_hook_proto_depIdxs,
		MessageInfos:      file_queue_hook_proto_msgTypes,
	}.Build()
	File_queue_hook_proto = out.File
	file_queue_hook_proto_goTypes = nil
	file_queue_hook_proto_depIdxs = nil
}
