// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: import_template.proto

package storage

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	engine "github.com/webitel/logger/api/engine"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
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

type ImportSourceType int32

const (
	ImportSourceType_DefaultSourceType ImportSourceType = 0
	ImportSourceType_Dialer            ImportSourceType = 1
)

// Enum value maps for ImportSourceType.
var (
	ImportSourceType_name = map[int32]string{
		0: "DefaultSourceType",
		1: "Dialer",
	}
	ImportSourceType_value = map[string]int32{
		"DefaultSourceType": 0,
		"Dialer":            1,
	}
)

func (x ImportSourceType) Enum() *ImportSourceType {
	p := new(ImportSourceType)
	*p = x
	return p
}

func (x ImportSourceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ImportSourceType) Descriptor() protoreflect.EnumDescriptor {
	return file_import_template_proto_enumTypes[0].Descriptor()
}

func (ImportSourceType) Type() protoreflect.EnumType {
	return &file_import_template_proto_enumTypes[0]
}

func (x ImportSourceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ImportSourceType.Descriptor instead.
func (ImportSourceType) EnumDescriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{0}
}

type ImportTemplate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	SourceType    ImportSourceType       `protobuf:"varint,4,opt,name=source_type,json=sourceType,proto3,enum=storage.ImportSourceType" json:"source_type,omitempty"`
	SourceId      int64                  `protobuf:"varint,5,opt,name=source_id,json=sourceId,proto3" json:"source_id,omitempty"`
	Parameters    *structpb.Struct       `protobuf:"bytes,6,opt,name=parameters,proto3" json:"parameters,omitempty"`
	Source        *engine.Lookup         `protobuf:"bytes,7,opt,name=source,proto3" json:"source,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ImportTemplate) Reset() {
	*x = ImportTemplate{}
	mi := &file_import_template_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ImportTemplate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportTemplate) ProtoMessage() {}

func (x *ImportTemplate) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportTemplate.ProtoReflect.Descriptor instead.
func (*ImportTemplate) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{0}
}

func (x *ImportTemplate) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ImportTemplate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ImportTemplate) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ImportTemplate) GetSourceType() ImportSourceType {
	if x != nil {
		return x.SourceType
	}
	return ImportSourceType_DefaultSourceType
}

func (x *ImportTemplate) GetSourceId() int64 {
	if x != nil {
		return x.SourceId
	}
	return 0
}

func (x *ImportTemplate) GetParameters() *structpb.Struct {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *ImportTemplate) GetSource() *engine.Lookup {
	if x != nil {
		return x.Source
	}
	return nil
}

type CreateImportTemplateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	SourceType    ImportSourceType       `protobuf:"varint,3,opt,name=source_type,json=sourceType,proto3,enum=storage.ImportSourceType" json:"source_type,omitempty"`
	SourceId      int64                  `protobuf:"varint,4,opt,name=source_id,json=sourceId,proto3" json:"source_id,omitempty"`
	Parameters    *structpb.Struct       `protobuf:"bytes,5,opt,name=parameters,proto3" json:"parameters,omitempty"`
	Source        *engine.Lookup         `protobuf:"bytes,6,opt,name=source,proto3" json:"source,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateImportTemplateRequest) Reset() {
	*x = CreateImportTemplateRequest{}
	mi := &file_import_template_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateImportTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateImportTemplateRequest) ProtoMessage() {}

func (x *CreateImportTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateImportTemplateRequest.ProtoReflect.Descriptor instead.
func (*CreateImportTemplateRequest) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{1}
}

func (x *CreateImportTemplateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateImportTemplateRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateImportTemplateRequest) GetSourceType() ImportSourceType {
	if x != nil {
		return x.SourceType
	}
	return ImportSourceType_DefaultSourceType
}

func (x *CreateImportTemplateRequest) GetSourceId() int64 {
	if x != nil {
		return x.SourceId
	}
	return 0
}

func (x *CreateImportTemplateRequest) GetParameters() *structpb.Struct {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *CreateImportTemplateRequest) GetSource() *engine.Lookup {
	if x != nil {
		return x.Source
	}
	return nil
}

type SearchImportTemplateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q             string                 `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	Sort          string                 `protobuf:"bytes,4,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields        []string               `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            []int32                `protobuf:"varint,6,rep,packed,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchImportTemplateRequest) Reset() {
	*x = SearchImportTemplateRequest{}
	mi := &file_import_template_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchImportTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchImportTemplateRequest) ProtoMessage() {}

func (x *SearchImportTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchImportTemplateRequest.ProtoReflect.Descriptor instead.
func (*SearchImportTemplateRequest) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{2}
}

func (x *SearchImportTemplateRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchImportTemplateRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchImportTemplateRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchImportTemplateRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchImportTemplateRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchImportTemplateRequest) GetId() []int32 {
	if x != nil {
		return x.Id
	}
	return nil
}

type ListImportTemplate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*ImportTemplate      `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListImportTemplate) Reset() {
	*x = ListImportTemplate{}
	mi := &file_import_template_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListImportTemplate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListImportTemplate) ProtoMessage() {}

func (x *ListImportTemplate) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListImportTemplate.ProtoReflect.Descriptor instead.
func (*ListImportTemplate) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{3}
}

func (x *ListImportTemplate) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListImportTemplate) GetItems() []*ImportTemplate {
	if x != nil {
		return x.Items
	}
	return nil
}

type ReadImportTemplateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadImportTemplateRequest) Reset() {
	*x = ReadImportTemplateRequest{}
	mi := &file_import_template_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadImportTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadImportTemplateRequest) ProtoMessage() {}

func (x *ReadImportTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadImportTemplateRequest.ProtoReflect.Descriptor instead.
func (*ReadImportTemplateRequest) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{4}
}

func (x *ReadImportTemplateRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UpdateImportTemplateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Parameters    *structpb.Struct       `protobuf:"bytes,4,opt,name=parameters,proto3" json:"parameters,omitempty"`
	Source        *engine.Lookup         `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateImportTemplateRequest) Reset() {
	*x = UpdateImportTemplateRequest{}
	mi := &file_import_template_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateImportTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateImportTemplateRequest) ProtoMessage() {}

func (x *UpdateImportTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateImportTemplateRequest.ProtoReflect.Descriptor instead.
func (*UpdateImportTemplateRequest) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateImportTemplateRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateImportTemplateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateImportTemplateRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateImportTemplateRequest) GetParameters() *structpb.Struct {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *UpdateImportTemplateRequest) GetSource() *engine.Lookup {
	if x != nil {
		return x.Source
	}
	return nil
}

type PatchImportTemplateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Parameters    *structpb.Struct       `protobuf:"bytes,4,opt,name=parameters,proto3" json:"parameters,omitempty"`
	Fields        []string               `protobuf:"bytes,50,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchImportTemplateRequest) Reset() {
	*x = PatchImportTemplateRequest{}
	mi := &file_import_template_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchImportTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchImportTemplateRequest) ProtoMessage() {}

func (x *PatchImportTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchImportTemplateRequest.ProtoReflect.Descriptor instead.
func (*PatchImportTemplateRequest) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{6}
}

func (x *PatchImportTemplateRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PatchImportTemplateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PatchImportTemplateRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PatchImportTemplateRequest) GetParameters() *structpb.Struct {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *PatchImportTemplateRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type DeleteImportTemplateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteImportTemplateRequest) Reset() {
	*x = DeleteImportTemplateRequest{}
	mi := &file_import_template_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteImportTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteImportTemplateRequest) ProtoMessage() {}

func (x *DeleteImportTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_import_template_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteImportTemplateRequest.ProtoReflect.Descriptor instead.
func (*DeleteImportTemplateRequest) Descriptor() ([]byte, []int) {
	return file_import_template_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteImportTemplateRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_import_template_proto protoreflect.FileDescriptor

const file_import_template_proto_rawDesc = "" +
	"\n" +
	"\x15import_template.proto\x12\astorage\x1a\vconst.proto\x1a\x1cgoogle/protobuf/struct.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\x90\x02\n" +
	"\x0eImportTemplate\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12:\n" +
	"\vsource_type\x18\x04 \x01(\x0e2\x19.storage.ImportSourceTypeR\n" +
	"sourceType\x12\x1b\n" +
	"\tsource_id\x18\x05 \x01(\x03R\bsourceId\x127\n" +
	"\n" +
	"parameters\x18\x06 \x01(\v2\x17.google.protobuf.StructR\n" +
	"parameters\x12&\n" +
	"\x06source\x18\a \x01(\v2\x0e.engine.LookupR\x06source\"\xf9\x02\n" +
	"\x1bCreateImportTemplateRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12:\n" +
	"\vsource_type\x18\x03 \x01(\x0e2\x19.storage.ImportSourceTypeR\n" +
	"sourceType\x12\x1b\n" +
	"\tsource_id\x18\x04 \x01(\x03R\bsourceId\x127\n" +
	"\n" +
	"parameters\x18\x05 \x01(\v2\x17.google.protobuf.StructR\n" +
	"parameters\x12&\n" +
	"\x06source\x18\x06 \x01(\v2\x0e.engine.LookupR\x06source:j\x92Ag\n" +
	"e*#Create import template request body2\x1eCreate import template for CSV\xd2\x01\x04name\xd2\x01\tsource_id\xd2\x01\n" +
	"parameters\"\x8f\x01\n" +
	"\x1bSearchImportTemplateRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\f\n" +
	"\x01q\x18\x03 \x01(\tR\x01q\x12\x12\n" +
	"\x04sort\x18\x04 \x01(\tR\x04sort\x12\x16\n" +
	"\x06fields\x18\x05 \x03(\tR\x06fields\x12\x0e\n" +
	"\x02id\x18\x06 \x03(\x05R\x02id\"W\n" +
	"\x12ListImportTemplate\x12\x12\n" +
	"\x04next\x18\x01 \x01(\bR\x04next\x12-\n" +
	"\x05items\x18\x02 \x03(\v2\x17.storage.ImportTemplateR\x05items\"+\n" +
	"\x19ReadImportTemplateRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\"\x9c\x02\n" +
	"\x1bUpdateImportTemplateRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x127\n" +
	"\n" +
	"parameters\x18\x04 \x01(\v2\x17.google.protobuf.StructR\n" +
	"parameters\x12&\n" +
	"\x06source\x18\x05 \x01(\v2\x0e.engine.LookupR\x06source:V\x92AS\n" +
	"Q*#Update import template request body2\x1eUpdate import template for CSV\xd2\x01\x02id\xd2\x01\x04name\"\x82\x02\n" +
	"\x1aPatchImportTemplateRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x127\n" +
	"\n" +
	"parameters\x18\x04 \x01(\v2\x17.google.protobuf.StructR\n" +
	"parameters\x12\x16\n" +
	"\x06fields\x182 \x03(\tR\x06fields:M\x92AJ\n" +
	"H*\"Patch import template request body2\x1dPatch import template for CSV\xd2\x01\x02id\"q\n" +
	"\x1bDeleteImportTemplateRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id:B\x92A?\n" +
	"=*\x1eDelete import template request2\x16Delete import template\xd2\x01\x02id*5\n" +
	"\x10ImportSourceType\x12\x15\n" +
	"\x11DefaultSourceType\x10\x00\x12\n" +
	"\n" +
	"\x06Dialer\x10\x012\x8f\x06\n" +
	"\x15ImportTemplateService\x12{\n" +
	"\x14CreateImportTemplate\x12$.storage.CreateImportTemplateRequest\x1a\x17.storage.ImportTemplate\"$\x82\xd3\xe4\x93\x02\x1e:\x01*\"\x19/storage/import_templates\x12|\n" +
	"\x14SearchImportTemplate\x12$.storage.SearchImportTemplateRequest\x1a\x1b.storage.ListImportTemplate\"!\x82\xd3\xe4\x93\x02\x1b\x12\x19/storage/import_templates\x12y\n" +
	"\x12ReadImportTemplate\x12\".storage.ReadImportTemplateRequest\x1a\x17.storage.ImportTemplate\"&\x82\xd3\xe4\x93\x02 \x12\x1e/storage/import_templates/{id}\x12\x80\x01\n" +
	"\x14UpdateImportTemplate\x12$.storage.UpdateImportTemplateRequest\x1a\x17.storage.ImportTemplate\")\x82\xd3\xe4\x93\x02#:\x01*\x1a\x1e/storage/import_templates/{id}\x12~\n" +
	"\x13PatchImportTemplate\x12#.storage.PatchImportTemplateRequest\x1a\x17.storage.ImportTemplate\")\x82\xd3\xe4\x93\x02#:\x01*2\x1e/storage/import_templates/{id}\x12}\n" +
	"\x14DeleteImportTemplate\x12$.storage.DeleteImportTemplateRequest\x1a\x17.storage.ImportTemplate\"&\x82\xd3\xe4\x93\x02 *\x1e/storage/import_templates/{id}B\x81\x01\n" +
	"\vcom.storageB\x13ImportTemplateProtoP\x01Z!github.com/webitel/protos/storage\xa2\x02\x03SXX\xaa\x02\aStorage\xca\x02\aStorage\xe2\x02\x13Storage\\GPBMetadata\xea\x02\aStorageb\x06proto3"

var (
	file_import_template_proto_rawDescOnce sync.Once
	file_import_template_proto_rawDescData []byte
)

func file_import_template_proto_rawDescGZIP() []byte {
	file_import_template_proto_rawDescOnce.Do(func() {
		file_import_template_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_import_template_proto_rawDesc), len(file_import_template_proto_rawDesc)))
	})
	return file_import_template_proto_rawDescData
}

var file_import_template_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_import_template_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_import_template_proto_goTypes = []any{
	(ImportSourceType)(0),               // 0: storage.ImportSourceType
	(*ImportTemplate)(nil),              // 1: storage.ImportTemplate
	(*CreateImportTemplateRequest)(nil), // 2: storage.CreateImportTemplateRequest
	(*SearchImportTemplateRequest)(nil), // 3: storage.SearchImportTemplateRequest
	(*ListImportTemplate)(nil),          // 4: storage.ListImportTemplate
	(*ReadImportTemplateRequest)(nil),   // 5: storage.ReadImportTemplateRequest
	(*UpdateImportTemplateRequest)(nil), // 6: storage.UpdateImportTemplateRequest
	(*PatchImportTemplateRequest)(nil),  // 7: storage.PatchImportTemplateRequest
	(*DeleteImportTemplateRequest)(nil), // 8: storage.DeleteImportTemplateRequest
	(*structpb.Struct)(nil),             // 9: google.protobuf.Struct
	(*engine.Lookup)(nil),               // 10: engine.Lookup
}
var file_import_template_proto_depIdxs = []int32{
	0,  // 0: storage.ImportTemplate.source_type:type_name -> storage.ImportSourceType
	9,  // 1: storage.ImportTemplate.parameters:type_name -> google.protobuf.Struct
	10, // 2: storage.ImportTemplate.source:type_name -> engine.Lookup
	0,  // 3: storage.CreateImportTemplateRequest.source_type:type_name -> storage.ImportSourceType
	9,  // 4: storage.CreateImportTemplateRequest.parameters:type_name -> google.protobuf.Struct
	10, // 5: storage.CreateImportTemplateRequest.source:type_name -> engine.Lookup
	1,  // 6: storage.ListImportTemplate.items:type_name -> storage.ImportTemplate
	9,  // 7: storage.UpdateImportTemplateRequest.parameters:type_name -> google.protobuf.Struct
	10, // 8: storage.UpdateImportTemplateRequest.source:type_name -> engine.Lookup
	9,  // 9: storage.PatchImportTemplateRequest.parameters:type_name -> google.protobuf.Struct
	2,  // 10: storage.ImportTemplateService.CreateImportTemplate:input_type -> storage.CreateImportTemplateRequest
	3,  // 11: storage.ImportTemplateService.SearchImportTemplate:input_type -> storage.SearchImportTemplateRequest
	5,  // 12: storage.ImportTemplateService.ReadImportTemplate:input_type -> storage.ReadImportTemplateRequest
	6,  // 13: storage.ImportTemplateService.UpdateImportTemplate:input_type -> storage.UpdateImportTemplateRequest
	7,  // 14: storage.ImportTemplateService.PatchImportTemplate:input_type -> storage.PatchImportTemplateRequest
	8,  // 15: storage.ImportTemplateService.DeleteImportTemplate:input_type -> storage.DeleteImportTemplateRequest
	1,  // 16: storage.ImportTemplateService.CreateImportTemplate:output_type -> storage.ImportTemplate
	4,  // 17: storage.ImportTemplateService.SearchImportTemplate:output_type -> storage.ListImportTemplate
	1,  // 18: storage.ImportTemplateService.ReadImportTemplate:output_type -> storage.ImportTemplate
	1,  // 19: storage.ImportTemplateService.UpdateImportTemplate:output_type -> storage.ImportTemplate
	1,  // 20: storage.ImportTemplateService.PatchImportTemplate:output_type -> storage.ImportTemplate
	1,  // 21: storage.ImportTemplateService.DeleteImportTemplate:output_type -> storage.ImportTemplate
	16, // [16:22] is the sub-list for method output_type
	10, // [10:16] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_import_template_proto_init() }
func file_import_template_proto_init() {
	if File_import_template_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_import_template_proto_rawDesc), len(file_import_template_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_import_template_proto_goTypes,
		DependencyIndexes: file_import_template_proto_depIdxs,
		EnumInfos:         file_import_template_proto_enumTypes,
		MessageInfos:      file_import_template_proto_msgTypes,
	}.Build()
	File_import_template_proto = out.File
	file_import_template_proto_goTypes = nil
	file_import_template_proto_depIdxs = nil
}
