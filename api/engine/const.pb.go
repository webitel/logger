// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: const.proto

package engine

import (
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

type BoolFilter int32

const (
	BoolFilter_undefined BoolFilter = 0
	BoolFilter_true      BoolFilter = 1
	BoolFilter_false     BoolFilter = 2
)

// Enum value maps for BoolFilter.
var (
	BoolFilter_name = map[int32]string{
		0: "undefined",
		1: "true",
		2: "false",
	}
	BoolFilter_value = map[string]int32{
		"undefined": 0,
		"true":      1,
		"false":     2,
	}
)

func (x BoolFilter) Enum() *BoolFilter {
	p := new(BoolFilter)
	*p = x
	return p
}

func (x BoolFilter) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BoolFilter) Descriptor() protoreflect.EnumDescriptor {
	return file_const_proto_enumTypes[0].Descriptor()
}

func (BoolFilter) Type() protoreflect.EnumType {
	return &file_const_proto_enumTypes[0]
}

func (x BoolFilter) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BoolFilter.Descriptor instead.
func (BoolFilter) EnumDescriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{0}
}

type Lookup struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Lookup) Reset() {
	*x = Lookup{}
	mi := &file_const_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Lookup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Lookup) ProtoMessage() {}

func (x *Lookup) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Lookup.ProtoReflect.Descriptor instead.
func (*Lookup) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{0}
}

func (x *Lookup) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Lookup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type FilterBetween struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	From          int64                  `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
	To            int64                  `protobuf:"varint,2,opt,name=to,proto3" json:"to,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FilterBetween) Reset() {
	*x = FilterBetween{}
	mi := &file_const_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterBetween) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterBetween) ProtoMessage() {}

func (x *FilterBetween) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterBetween.ProtoReflect.Descriptor instead.
func (*FilterBetween) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{1}
}

func (x *FilterBetween) GetFrom() int64 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *FilterBetween) GetTo() int64 {
	if x != nil {
		return x.To
	}
	return 0
}

type ListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DomainId      int64                  `protobuf:"varint,1,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Size          int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Page          int32                  `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	mi := &file_const_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{2}
}

func (x *ListRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *ListRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type DomainRecord struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DomainId      int64                  `protobuf:"varint,1,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	CreatedBy     int64                  `protobuf:"varint,3,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	UpdatedAt     int64                  `protobuf:"varint,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	UpdatedBy     int64                  `protobuf:"varint,5,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DomainRecord) Reset() {
	*x = DomainRecord{}
	mi := &file_const_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DomainRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainRecord) ProtoMessage() {}

func (x *DomainRecord) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainRecord.ProtoReflect.Descriptor instead.
func (*DomainRecord) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{3}
}

func (x *DomainRecord) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *DomainRecord) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *DomainRecord) GetCreatedBy() int64 {
	if x != nil {
		return x.CreatedBy
	}
	return 0
}

func (x *DomainRecord) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *DomainRecord) GetUpdatedBy() int64 {
	if x != nil {
		return x.UpdatedBy
	}
	return 0
}

type ListForItemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DomainId      int64                  `protobuf:"varint,1,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	ItemId        int64                  `protobuf:"varint,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	Size          int32                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Page          int32                  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListForItemRequest) Reset() {
	*x = ListForItemRequest{}
	mi := &file_const_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListForItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListForItemRequest) ProtoMessage() {}

func (x *ListForItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListForItemRequest.ProtoReflect.Descriptor instead.
func (*ListForItemRequest) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{4}
}

func (x *ListForItemRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *ListForItemRequest) GetItemId() int64 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *ListForItemRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListForItemRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type ItemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DomainId      int64                  `protobuf:"varint,1,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Id            int64                  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ItemRequest) Reset() {
	*x = ItemRequest{}
	mi := &file_const_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItemRequest) ProtoMessage() {}

func (x *ItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItemRequest.ProtoReflect.Descriptor instead.
func (*ItemRequest) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{5}
}

func (x *ItemRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *ItemRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file_const_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{6}
}

func (x *Response) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type Tag struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tag) Reset() {
	*x = Tag{}
	mi := &file_const_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{7}
}

func (x *Tag) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ListTags struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*Tag                 `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListTags) Reset() {
	*x = ListTags{}
	mi := &file_const_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTags) ProtoMessage() {}

func (x *ListTags) ProtoReflect() protoreflect.Message {
	mi := &file_const_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTags.ProtoReflect.Descriptor instead.
func (*ListTags) Descriptor() ([]byte, []int) {
	return file_const_proto_rawDescGZIP(), []int{8}
}

func (x *ListTags) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListTags) GetItems() []*Tag {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_const_proto protoreflect.FileDescriptor

const file_const_proto_rawDesc = "" +
	"\n" +
	"\vconst.proto\x12\x06engine\",\n" +
	"\x06Lookup\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\"3\n" +
	"\rFilterBetween\x12\x12\n" +
	"\x04from\x18\x01 \x01(\x03R\x04from\x12\x0e\n" +
	"\x02to\x18\x02 \x01(\x03R\x02to\"R\n" +
	"\vListRequest\x12\x1b\n" +
	"\tdomain_id\x18\x01 \x01(\x03R\bdomainId\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\x12\n" +
	"\x04page\x18\x03 \x01(\x05R\x04page\"\xa7\x01\n" +
	"\fDomainRecord\x12\x1b\n" +
	"\tdomain_id\x18\x01 \x01(\x03R\bdomainId\x12\x1d\n" +
	"\n" +
	"created_at\x18\x02 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"created_by\x18\x03 \x01(\x03R\tcreatedBy\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x04 \x01(\x03R\tupdatedAt\x12\x1d\n" +
	"\n" +
	"updated_by\x18\x05 \x01(\x03R\tupdatedBy\"r\n" +
	"\x12ListForItemRequest\x12\x1b\n" +
	"\tdomain_id\x18\x01 \x01(\x03R\bdomainId\x12\x17\n" +
	"\aitem_id\x18\x02 \x01(\x03R\x06itemId\x12\x12\n" +
	"\x04size\x18\x03 \x01(\x05R\x04size\x12\x12\n" +
	"\x04page\x18\x04 \x01(\x05R\x04page\":\n" +
	"\vItemRequest\x12\x1b\n" +
	"\tdomain_id\x18\x01 \x01(\x03R\bdomainId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\x03R\x02id\"\"\n" +
	"\bResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\"\x19\n" +
	"\x03Tag\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"A\n" +
	"\bListTags\x12\x12\n" +
	"\x04next\x18\x01 \x01(\bR\x04next\x12!\n" +
	"\x05items\x18\x02 \x03(\v2\v.engine.TagR\x05items*0\n" +
	"\n" +
	"BoolFilter\x12\r\n" +
	"\tundefined\x10\x00\x12\b\n" +
	"\x04true\x10\x01\x12\t\n" +
	"\x05false\x10\x02B\"Z github.com/webitel/protos/engineb\x06proto3"

var (
	file_const_proto_rawDescOnce sync.Once
	file_const_proto_rawDescData []byte
)

func file_const_proto_rawDescGZIP() []byte {
	file_const_proto_rawDescOnce.Do(func() {
		file_const_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_const_proto_rawDesc), len(file_const_proto_rawDesc)))
	})
	return file_const_proto_rawDescData
}

var file_const_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_const_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_const_proto_goTypes = []any{
	(BoolFilter)(0),            // 0: engine.BoolFilter
	(*Lookup)(nil),             // 1: engine.Lookup
	(*FilterBetween)(nil),      // 2: engine.FilterBetween
	(*ListRequest)(nil),        // 3: engine.ListRequest
	(*DomainRecord)(nil),       // 4: engine.DomainRecord
	(*ListForItemRequest)(nil), // 5: engine.ListForItemRequest
	(*ItemRequest)(nil),        // 6: engine.ItemRequest
	(*Response)(nil),           // 7: engine.Response
	(*Tag)(nil),                // 8: engine.Tag
	(*ListTags)(nil),           // 9: engine.ListTags
}
var file_const_proto_depIdxs = []int32{
	8, // 0: engine.ListTags.items:type_name -> engine.Tag
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_const_proto_init() }
func file_const_proto_init() {
	if File_const_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_const_proto_rawDesc), len(file_const_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_const_proto_goTypes,
		DependencyIndexes: file_const_proto_depIdxs,
		EnumInfos:         file_const_proto_enumTypes,
		MessageInfos:      file_const_proto_msgTypes,
	}.Build()
	File_const_proto = out.File
	file_const_proto_goTypes = nil
	file_const_proto_depIdxs = nil
}
