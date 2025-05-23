// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: knowledgebase/space.proto

package spaces

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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

type Space struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the association. Never changes.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// READONLY. Operational attributes
	// Version of the latest update. Numeric sequence.
	Ver int32 `protobuf:"varint,2,opt,name=ver,proto3" json:"ver,omitempty"`
	// Unique ID of the latest version of the update.
	// This ID changes after any update to the underlying value(s).
	Etag string `protobuf:"bytes,3,opt,name=etag,proto3" json:"etag,omitempty"`
	// [R]ecord[b]ased[A]ccess[C]ontrol mode granted.
	Mode string `protobuf:"bytes,4,opt,name=mode,proto3" json:"mode,omitempty"`
	// READONLY. The space's metadata.
	Domain *Lookup `protobuf:"bytes,5,opt,name=domain,proto3" json:"domain,omitempty"`
	// The timestamp when the space was created (in Unix time).
	CreatedAt int64 `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The user who created the space.
	CreatedBy *Lookup `protobuf:"bytes,7,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// The timestamp when the space was last updated (in Unix time).
	UpdatedAt int64 `protobuf:"varint,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The user who last updated the space.
	UpdatedBy *Lookup `protobuf:"bytes,9,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// The name of the space.
	Name string `protobuf:"bytes,10,opt,name=name,proto3" json:"name,omitempty"`
	// BIO. Short description about the space.
	HomePage string `protobuf:"bytes,11,opt,name=home_page,json=homePage,proto3" json:"home_page,omitempty"`
	// The state of the space.
	State bool `protobuf:"varint,12,opt,name=state,proto3" json:"state,omitempty"`
	// Indicates if the space has children.
	HasChildren   bool `protobuf:"varint,13,opt,name=has_children,json=hasChildren,proto3" json:"has_children,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Space) Reset() {
	*x = Space{}
	mi := &file_knowledgebase_space_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Space) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Space) ProtoMessage() {}

func (x *Space) ProtoReflect() protoreflect.Message {
	mi := &file_knowledgebase_space_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Space.ProtoReflect.Descriptor instead.
func (*Space) Descriptor() ([]byte, []int) {
	return file_knowledgebase_space_proto_rawDescGZIP(), []int{0}
}

func (x *Space) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Space) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *Space) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *Space) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

func (x *Space) GetDomain() *Lookup {
	if x != nil {
		return x.Domain
	}
	return nil
}

func (x *Space) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Space) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Space) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Space) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *Space) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Space) GetHomePage() string {
	if x != nil {
		return x.HomePage
	}
	return ""
}

func (x *Space) GetState() bool {
	if x != nil {
		return x.State
	}
	return false
}

func (x *Space) GetHasChildren() bool {
	if x != nil {
		return x.HasChildren
	}
	return false
}

// The Space principal input.
type InputSpace struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique ID of the latest version of an existing resorce.
	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	// Represents the name of the knowledge base space.
	Name string `protobuf:"bytes,10,opt,name=name,proto3" json:"name,omitempty"`
	// BIO. Short description about the space.
	// OPTIONAL. Multi-lined text.
	HomePage string `protobuf:"bytes,11,opt,name=home_page,json=homePage,proto3" json:"home_page,omitempty"`
	// The state of the space.
	State         bool `protobuf:"varint,12,opt,name=state,proto3" json:"state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputSpace) Reset() {
	*x = InputSpace{}
	mi := &file_knowledgebase_space_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputSpace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputSpace) ProtoMessage() {}

func (x *InputSpace) ProtoReflect() protoreflect.Message {
	mi := &file_knowledgebase_space_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputSpace.ProtoReflect.Descriptor instead.
func (*InputSpace) Descriptor() ([]byte, []int) {
	return file_knowledgebase_space_proto_rawDescGZIP(), []int{1}
}

func (x *InputSpace) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *InputSpace) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputSpace) GetHomePage() string {
	if x != nil {
		return x.HomePage
	}
	return ""
}

func (x *InputSpace) GetState() bool {
	if x != nil {
		return x.State
	}
	return false
}

var File_knowledgebase_space_proto protoreflect.FileDescriptor

const file_knowledgebase_space_proto_rawDesc = "" +
	"\n" +
	"\x19knowledgebase/space.proto\x12\x15webitel.knowledgebase\x1a\x1aknowledgebase/fields.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\xac\x03\n" +
	"\x05Space\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x10\n" +
	"\x03ver\x18\x02 \x01(\x05R\x03ver\x12\x12\n" +
	"\x04etag\x18\x03 \x01(\tR\x04etag\x12\x12\n" +
	"\x04mode\x18\x04 \x01(\tR\x04mode\x125\n" +
	"\x06domain\x18\x05 \x01(\v2\x1d.webitel.knowledgebase.LookupR\x06domain\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\x03R\tcreatedAt\x12<\n" +
	"\n" +
	"created_by\x18\a \x01(\v2\x1d.webitel.knowledgebase.LookupR\tcreatedBy\x12\x1d\n" +
	"\n" +
	"updated_at\x18\b \x01(\x03R\tupdatedAt\x12<\n" +
	"\n" +
	"updated_by\x18\t \x01(\v2\x1d.webitel.knowledgebase.LookupR\tupdatedBy\x12\x12\n" +
	"\x04name\x18\n" +
	" \x01(\tR\x04name\x12\x1b\n" +
	"\thome_page\x18\v \x01(\tR\bhomePage\x12\x14\n" +
	"\x05state\x18\f \x01(\bR\x05state\x12!\n" +
	"\fhas_children\x18\r \x01(\bR\vhasChildren\"v\n" +
	"\n" +
	"InputSpace\x12!\n" +
	"\x04etag\x18\x01 \x01(\tB\r\x92A\n" +
	"\xca>\a\xfa\x02\x04etagR\x04etag\x12\x12\n" +
	"\x04name\x18\n" +
	" \x01(\tR\x04name\x12\x1b\n" +
	"\thome_page\x18\v \x01(\tR\bhomePage\x12\x14\n" +
	"\x05state\x18\f \x01(\bR\x05stateB\xba\x01\n" +
	"\x19com.webitel.knowledgebaseB\n" +
	"SpaceProtoP\x01Z\x1cwebitel.go/api/spaces;spaces\xa2\x02\x03WKX\xaa\x02\x15Webitel.Knowledgebase\xca\x02\x15Webitel\\Knowledgebase\xe2\x02!Webitel\\Knowledgebase\\GPBMetadata\xea\x02\x16Webitel::Knowledgebaseb\x06proto3"

var (
	file_knowledgebase_space_proto_rawDescOnce sync.Once
	file_knowledgebase_space_proto_rawDescData []byte
)

func file_knowledgebase_space_proto_rawDescGZIP() []byte {
	file_knowledgebase_space_proto_rawDescOnce.Do(func() {
		file_knowledgebase_space_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_knowledgebase_space_proto_rawDesc), len(file_knowledgebase_space_proto_rawDesc)))
	})
	return file_knowledgebase_space_proto_rawDescData
}

var file_knowledgebase_space_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_knowledgebase_space_proto_goTypes = []any{
	(*Space)(nil),      // 0: webitel.knowledgebase.Space
	(*InputSpace)(nil), // 1: webitel.knowledgebase.InputSpace
	(*Lookup)(nil),     // 2: webitel.knowledgebase.Lookup
}
var file_knowledgebase_space_proto_depIdxs = []int32{
	2, // 0: webitel.knowledgebase.Space.domain:type_name -> webitel.knowledgebase.Lookup
	2, // 1: webitel.knowledgebase.Space.created_by:type_name -> webitel.knowledgebase.Lookup
	2, // 2: webitel.knowledgebase.Space.updated_by:type_name -> webitel.knowledgebase.Lookup
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_knowledgebase_space_proto_init() }
func file_knowledgebase_space_proto_init() {
	if File_knowledgebase_space_proto != nil {
		return
	}
	file_knowledgebase_fields_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_knowledgebase_space_proto_rawDesc), len(file_knowledgebase_space_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_knowledgebase_space_proto_goTypes,
		DependencyIndexes: file_knowledgebase_space_proto_depIdxs,
		MessageInfos:      file_knowledgebase_space_proto_msgTypes,
	}.Build()
	File_knowledgebase_space_proto = out.File
	file_knowledgebase_space_proto_goTypes = nil
	file_knowledgebase_space_proto_depIdxs = nil
}
