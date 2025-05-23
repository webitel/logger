// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: queue_bucket.proto

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

type DeleteQueueBucketRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	QueueId       int64                  `protobuf:"varint,2,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteQueueBucketRequest) Reset() {
	*x = DeleteQueueBucketRequest{}
	mi := &file_queue_bucket_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteQueueBucketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteQueueBucketRequest) ProtoMessage() {}

func (x *DeleteQueueBucketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteQueueBucketRequest.ProtoReflect.Descriptor instead.
func (*DeleteQueueBucketRequest) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteQueueBucketRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DeleteQueueBucketRequest) GetQueueId() int64 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

type UpdateQueueBucketRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	QueueId       int64                  `protobuf:"varint,2,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Ratio         int32                  `protobuf:"varint,3,opt,name=ratio,proto3" json:"ratio,omitempty"`
	Bucket        *Lookup                `protobuf:"bytes,4,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Disabled      bool                   `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Priority      int32                  `protobuf:"varint,6,opt,name=priority,proto3" json:"priority,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateQueueBucketRequest) Reset() {
	*x = UpdateQueueBucketRequest{}
	mi := &file_queue_bucket_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateQueueBucketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateQueueBucketRequest) ProtoMessage() {}

func (x *UpdateQueueBucketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateQueueBucketRequest.ProtoReflect.Descriptor instead.
func (*UpdateQueueBucketRequest) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateQueueBucketRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateQueueBucketRequest) GetQueueId() int64 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *UpdateQueueBucketRequest) GetRatio() int32 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

func (x *UpdateQueueBucketRequest) GetBucket() *Lookup {
	if x != nil {
		return x.Bucket
	}
	return nil
}

func (x *UpdateQueueBucketRequest) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *UpdateQueueBucketRequest) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type PatchQueueBucketRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	QueueId       int64                  `protobuf:"varint,2,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Ratio         int32                  `protobuf:"varint,3,opt,name=ratio,proto3" json:"ratio,omitempty"`
	Bucket        *Lookup                `protobuf:"bytes,4,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Disabled      bool                   `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Priority      int32                  `protobuf:"varint,6,opt,name=priority,proto3" json:"priority,omitempty"`
	Fields        []string               `protobuf:"bytes,7,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchQueueBucketRequest) Reset() {
	*x = PatchQueueBucketRequest{}
	mi := &file_queue_bucket_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchQueueBucketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchQueueBucketRequest) ProtoMessage() {}

func (x *PatchQueueBucketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchQueueBucketRequest.ProtoReflect.Descriptor instead.
func (*PatchQueueBucketRequest) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{2}
}

func (x *PatchQueueBucketRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PatchQueueBucketRequest) GetQueueId() int64 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *PatchQueueBucketRequest) GetRatio() int32 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

func (x *PatchQueueBucketRequest) GetBucket() *Lookup {
	if x != nil {
		return x.Bucket
	}
	return nil
}

func (x *PatchQueueBucketRequest) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *PatchQueueBucketRequest) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *PatchQueueBucketRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type SearchQueueBucketRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       int64                  `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Q             string                 `protobuf:"bytes,4,opt,name=q,proto3" json:"q,omitempty"`
	Sort          string                 `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields        []string               `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            []uint32               `protobuf:"varint,7,rep,packed,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchQueueBucketRequest) Reset() {
	*x = SearchQueueBucketRequest{}
	mi := &file_queue_bucket_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchQueueBucketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchQueueBucketRequest) ProtoMessage() {}

func (x *SearchQueueBucketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchQueueBucketRequest.ProtoReflect.Descriptor instead.
func (*SearchQueueBucketRequest) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{3}
}

func (x *SearchQueueBucketRequest) GetQueueId() int64 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *SearchQueueBucketRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchQueueBucketRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchQueueBucketRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchQueueBucketRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchQueueBucketRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchQueueBucketRequest) GetId() []uint32 {
	if x != nil {
		return x.Id
	}
	return nil
}

type ListQueueBucket struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*QueueBucket         `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListQueueBucket) Reset() {
	*x = ListQueueBucket{}
	mi := &file_queue_bucket_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListQueueBucket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueueBucket) ProtoMessage() {}

func (x *ListQueueBucket) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueueBucket.ProtoReflect.Descriptor instead.
func (*ListQueueBucket) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{4}
}

func (x *ListQueueBucket) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListQueueBucket) GetItems() []*QueueBucket {
	if x != nil {
		return x.Items
	}
	return nil
}

type ReadQueueBucketRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	QueueId       int64                  `protobuf:"varint,2,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadQueueBucketRequest) Reset() {
	*x = ReadQueueBucketRequest{}
	mi := &file_queue_bucket_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadQueueBucketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadQueueBucketRequest) ProtoMessage() {}

func (x *ReadQueueBucketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadQueueBucketRequest.ProtoReflect.Descriptor instead.
func (*ReadQueueBucketRequest) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{5}
}

func (x *ReadQueueBucketRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ReadQueueBucketRequest) GetQueueId() int64 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

type CreateQueueBucketRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QueueId       int64                  `protobuf:"varint,1,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
	Ratio         int32                  `protobuf:"varint,2,opt,name=ratio,proto3" json:"ratio,omitempty"`
	Bucket        *Lookup                `protobuf:"bytes,3,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Disabled      bool                   `protobuf:"varint,4,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Priority      int32                  `protobuf:"varint,5,opt,name=priority,proto3" json:"priority,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateQueueBucketRequest) Reset() {
	*x = CreateQueueBucketRequest{}
	mi := &file_queue_bucket_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateQueueBucketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateQueueBucketRequest) ProtoMessage() {}

func (x *CreateQueueBucketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateQueueBucketRequest.ProtoReflect.Descriptor instead.
func (*CreateQueueBucketRequest) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{6}
}

func (x *CreateQueueBucketRequest) GetQueueId() int64 {
	if x != nil {
		return x.QueueId
	}
	return 0
}

func (x *CreateQueueBucketRequest) GetRatio() int32 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

func (x *CreateQueueBucketRequest) GetBucket() *Lookup {
	if x != nil {
		return x.Bucket
	}
	return nil
}

func (x *CreateQueueBucketRequest) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *CreateQueueBucketRequest) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type QueueBucket struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Ratio         int32                  `protobuf:"varint,2,opt,name=ratio,proto3" json:"ratio,omitempty"`
	Bucket        *Lookup                `protobuf:"bytes,3,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Disabled      bool                   `protobuf:"varint,4,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Priority      int32                  `protobuf:"varint,5,opt,name=priority,proto3" json:"priority,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueueBucket) Reset() {
	*x = QueueBucket{}
	mi := &file_queue_bucket_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueBucket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueBucket) ProtoMessage() {}

func (x *QueueBucket) ProtoReflect() protoreflect.Message {
	mi := &file_queue_bucket_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueBucket.ProtoReflect.Descriptor instead.
func (*QueueBucket) Descriptor() ([]byte, []int) {
	return file_queue_bucket_proto_rawDescGZIP(), []int{7}
}

func (x *QueueBucket) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *QueueBucket) GetRatio() int32 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

func (x *QueueBucket) GetBucket() *Lookup {
	if x != nil {
		return x.Bucket
	}
	return nil
}

func (x *QueueBucket) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *QueueBucket) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

var File_queue_bucket_proto protoreflect.FileDescriptor

const file_queue_bucket_proto_rawDesc = "" +
	"\n" +
	"\x12queue_bucket.proto\x12\x06engine\x1a\vconst.proto\x1a\x1cgoogle/api/annotations.proto\"E\n" +
	"\x18DeleteQueueBucketRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x19\n" +
	"\bqueue_id\x18\x02 \x01(\x03R\aqueueId\"\xbb\x01\n" +
	"\x18UpdateQueueBucketRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x19\n" +
	"\bqueue_id\x18\x02 \x01(\x03R\aqueueId\x12\x14\n" +
	"\x05ratio\x18\x03 \x01(\x05R\x05ratio\x12&\n" +
	"\x06bucket\x18\x04 \x01(\v2\x0e.engine.LookupR\x06bucket\x12\x1a\n" +
	"\bdisabled\x18\x05 \x01(\bR\bdisabled\x12\x1a\n" +
	"\bpriority\x18\x06 \x01(\x05R\bpriority\"\xd2\x01\n" +
	"\x17PatchQueueBucketRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x19\n" +
	"\bqueue_id\x18\x02 \x01(\x03R\aqueueId\x12\x14\n" +
	"\x05ratio\x18\x03 \x01(\x05R\x05ratio\x12&\n" +
	"\x06bucket\x18\x04 \x01(\v2\x0e.engine.LookupR\x06bucket\x12\x1a\n" +
	"\bdisabled\x18\x05 \x01(\bR\bdisabled\x12\x1a\n" +
	"\bpriority\x18\x06 \x01(\x05R\bpriority\x12\x16\n" +
	"\x06fields\x18\a \x03(\tR\x06fields\"\xa7\x01\n" +
	"\x18SearchQueueBucketRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\x03R\aqueueId\x12\x12\n" +
	"\x04page\x18\x02 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x03 \x01(\x05R\x04size\x12\f\n" +
	"\x01q\x18\x04 \x01(\tR\x01q\x12\x12\n" +
	"\x04sort\x18\x05 \x01(\tR\x04sort\x12\x16\n" +
	"\x06fields\x18\x06 \x03(\tR\x06fields\x12\x0e\n" +
	"\x02id\x18\a \x03(\rR\x02id\"P\n" +
	"\x0fListQueueBucket\x12\x12\n" +
	"\x04next\x18\x01 \x01(\bR\x04next\x12)\n" +
	"\x05items\x18\x02 \x03(\v2\x13.engine.QueueBucketR\x05items\"C\n" +
	"\x16ReadQueueBucketRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x19\n" +
	"\bqueue_id\x18\x02 \x01(\x03R\aqueueId\"\xab\x01\n" +
	"\x18CreateQueueBucketRequest\x12\x19\n" +
	"\bqueue_id\x18\x01 \x01(\x03R\aqueueId\x12\x14\n" +
	"\x05ratio\x18\x02 \x01(\x05R\x05ratio\x12&\n" +
	"\x06bucket\x18\x03 \x01(\v2\x0e.engine.LookupR\x06bucket\x12\x1a\n" +
	"\bdisabled\x18\x04 \x01(\bR\bdisabled\x12\x1a\n" +
	"\bpriority\x18\x05 \x01(\x05R\bpriority\"\x93\x01\n" +
	"\vQueueBucket\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x14\n" +
	"\x05ratio\x18\x02 \x01(\x05R\x05ratio\x12&\n" +
	"\x06bucket\x18\x03 \x01(\v2\x0e.engine.LookupR\x06bucket\x12\x1a\n" +
	"\bdisabled\x18\x04 \x01(\bR\bdisabled\x12\x1a\n" +
	"\bpriority\x18\x05 \x01(\x05R\bpriority2\x99\x06\n" +
	"\x12QueueBucketService\x12}\n" +
	"\x11CreateQueueBucket\x12 .engine.CreateQueueBucketRequest\x1a\x13.engine.QueueBucket\"1\x82\xd3\xe4\x93\x02+:\x01*\"&/call_center/queues/{queue_id}/buckets\x12~\n" +
	"\x11SearchQueueBucket\x12 .engine.SearchQueueBucketRequest\x1a\x17.engine.ListQueueBucket\".\x82\xd3\xe4\x93\x02(\x12&/call_center/queues/{queue_id}/buckets\x12{\n" +
	"\x0fReadQueueBucket\x12\x1e.engine.ReadQueueBucketRequest\x1a\x13.engine.QueueBucket\"3\x82\xd3\xe4\x93\x02-\x12+/call_center/queues/{queue_id}/buckets/{id}\x12\x82\x01\n" +
	"\x11UpdateQueueBucket\x12 .engine.UpdateQueueBucketRequest\x1a\x13.engine.QueueBucket\"6\x82\xd3\xe4\x93\x020:\x01*\x1a+/call_center/queues/{queue_id}/buckets/{id}\x12\x80\x01\n" +
	"\x10PatchQueueBucket\x12\x1f.engine.PatchQueueBucketRequest\x1a\x13.engine.QueueBucket\"6\x82\xd3\xe4\x93\x020:\x01*2+/call_center/queues/{queue_id}/buckets/{id}\x12\x7f\n" +
	"\x11DeleteQueueBucket\x12 .engine.DeleteQueueBucketRequest\x1a\x13.engine.QueueBucket\"3\x82\xd3\xe4\x93\x02-*+/call_center/queues/{queue_id}/buckets/{id}B\"Z github.com/webitel/protos/engineb\x06proto3"

var (
	file_queue_bucket_proto_rawDescOnce sync.Once
	file_queue_bucket_proto_rawDescData []byte
)

func file_queue_bucket_proto_rawDescGZIP() []byte {
	file_queue_bucket_proto_rawDescOnce.Do(func() {
		file_queue_bucket_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_queue_bucket_proto_rawDesc), len(file_queue_bucket_proto_rawDesc)))
	})
	return file_queue_bucket_proto_rawDescData
}

var file_queue_bucket_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_queue_bucket_proto_goTypes = []any{
	(*DeleteQueueBucketRequest)(nil), // 0: engine.DeleteQueueBucketRequest
	(*UpdateQueueBucketRequest)(nil), // 1: engine.UpdateQueueBucketRequest
	(*PatchQueueBucketRequest)(nil),  // 2: engine.PatchQueueBucketRequest
	(*SearchQueueBucketRequest)(nil), // 3: engine.SearchQueueBucketRequest
	(*ListQueueBucket)(nil),          // 4: engine.ListQueueBucket
	(*ReadQueueBucketRequest)(nil),   // 5: engine.ReadQueueBucketRequest
	(*CreateQueueBucketRequest)(nil), // 6: engine.CreateQueueBucketRequest
	(*QueueBucket)(nil),              // 7: engine.QueueBucket
	(*Lookup)(nil),                   // 8: engine.Lookup
}
var file_queue_bucket_proto_depIdxs = []int32{
	8,  // 0: engine.UpdateQueueBucketRequest.bucket:type_name -> engine.Lookup
	8,  // 1: engine.PatchQueueBucketRequest.bucket:type_name -> engine.Lookup
	7,  // 2: engine.ListQueueBucket.items:type_name -> engine.QueueBucket
	8,  // 3: engine.CreateQueueBucketRequest.bucket:type_name -> engine.Lookup
	8,  // 4: engine.QueueBucket.bucket:type_name -> engine.Lookup
	6,  // 5: engine.QueueBucketService.CreateQueueBucket:input_type -> engine.CreateQueueBucketRequest
	3,  // 6: engine.QueueBucketService.SearchQueueBucket:input_type -> engine.SearchQueueBucketRequest
	5,  // 7: engine.QueueBucketService.ReadQueueBucket:input_type -> engine.ReadQueueBucketRequest
	1,  // 8: engine.QueueBucketService.UpdateQueueBucket:input_type -> engine.UpdateQueueBucketRequest
	2,  // 9: engine.QueueBucketService.PatchQueueBucket:input_type -> engine.PatchQueueBucketRequest
	0,  // 10: engine.QueueBucketService.DeleteQueueBucket:input_type -> engine.DeleteQueueBucketRequest
	7,  // 11: engine.QueueBucketService.CreateQueueBucket:output_type -> engine.QueueBucket
	4,  // 12: engine.QueueBucketService.SearchQueueBucket:output_type -> engine.ListQueueBucket
	7,  // 13: engine.QueueBucketService.ReadQueueBucket:output_type -> engine.QueueBucket
	7,  // 14: engine.QueueBucketService.UpdateQueueBucket:output_type -> engine.QueueBucket
	7,  // 15: engine.QueueBucketService.PatchQueueBucket:output_type -> engine.QueueBucket
	7,  // 16: engine.QueueBucketService.DeleteQueueBucket:output_type -> engine.QueueBucket
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_queue_bucket_proto_init() }
func file_queue_bucket_proto_init() {
	if File_queue_bucket_proto != nil {
		return
	}
	file_const_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_queue_bucket_proto_rawDesc), len(file_queue_bucket_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_queue_bucket_proto_goTypes,
		DependencyIndexes: file_queue_bucket_proto_depIdxs,
		MessageInfos:      file_queue_bucket_proto_msgTypes,
	}.Build()
	File_queue_bucket_proto = out.File
	file_queue_bucket_proto_goTypes = nil
	file_queue_bucket_proto_depIdxs = nil
}
