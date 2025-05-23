// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: config_service.proto

package logger

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

type Configs struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next          bool                   `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*Config              `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Configs) Reset() {
	*x = Configs{}
	mi := &file_config_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Configs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configs) ProtoMessage() {}

func (x *Configs) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configs.ProtoReflect.Descriptor instead.
func (*Configs) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{0}
}

func (x *Configs) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *Configs) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *Configs) GetItems() []*Config {
	if x != nil {
		return x.Items
	}
	return nil
}

type SystemObjects struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Lookup              `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SystemObjects) Reset() {
	*x = SystemObjects{}
	mi := &file_config_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SystemObjects) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SystemObjects) ProtoMessage() {}

func (x *SystemObjects) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SystemObjects.ProtoReflect.Descriptor instead.
func (*SystemObjects) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{1}
}

func (x *SystemObjects) GetItems() []*Lookup {
	if x != nil {
		return x.Items
	}
	return nil
}

type ReadSystemObjectsRequest struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	IncludeExisting bool                   `protobuf:"varint,1,opt,name=include_existing,json=includeExisting,proto3" json:"include_existing,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ReadSystemObjectsRequest) Reset() {
	*x = ReadSystemObjectsRequest{}
	mi := &file_config_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadSystemObjectsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadSystemObjectsRequest) ProtoMessage() {}

func (x *ReadSystemObjectsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadSystemObjectsRequest.ProtoReflect.Descriptor instead.
func (*ReadSystemObjectsRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{2}
}

func (x *ReadSystemObjectsRequest) GetIncludeExisting() bool {
	if x != nil {
		return x.IncludeExisting
	}
	return false
}

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_config_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{3}
}

type Config struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Object        *Lookup                `protobuf:"bytes,2,opt,name=object,proto3" json:"object,omitempty"`
	Enabled       bool                   `protobuf:"varint,3,opt,name=enabled,proto3" json:"enabled,omitempty"`
	DaysToStore   int32                  `protobuf:"varint,4,opt,name=days_to_store,json=daysToStore,proto3" json:"days_to_store,omitempty"`
	Period        int32                  `protobuf:"varint,5,opt,name=period,proto3" json:"period,omitempty"`
	Storage       *Lookup                `protobuf:"bytes,6,opt,name=storage,proto3" json:"storage,omitempty"`
	Description   string                 `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	LogsSize      string                 `protobuf:"bytes,8,opt,name=logs_size,json=logsSize,proto3" json:"logs_size,omitempty"`
	LogsCount     int64                  `protobuf:"varint,9,opt,name=logs_count,json=logsCount,proto3" json:"logs_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_config_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{4}
}

func (x *Config) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Config) GetObject() *Lookup {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *Config) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Config) GetDaysToStore() int32 {
	if x != nil {
		return x.DaysToStore
	}
	return 0
}

func (x *Config) GetPeriod() int32 {
	if x != nil {
		return x.Period
	}
	return 0
}

func (x *Config) GetStorage() *Lookup {
	if x != nil {
		return x.Storage
	}
	return nil
}

func (x *Config) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Config) GetLogsSize() string {
	if x != nil {
		return x.LogsSize
	}
	return ""
}

func (x *Config) GetLogsCount() int64 {
	if x != nil {
		return x.LogsCount
	}
	return 0
}

type DeleteConfigRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigId      int32                  `protobuf:"varint,1,opt,name=config_id,json=configId,proto3" json:"config_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteConfigRequest) Reset() {
	*x = DeleteConfigRequest{}
	mi := &file_config_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteConfigRequest) ProtoMessage() {}

func (x *DeleteConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteConfigRequest.ProtoReflect.Descriptor instead.
func (*DeleteConfigRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteConfigRequest) GetConfigId() int32 {
	if x != nil {
		return x.ConfigId
	}
	return 0
}

type ConfigStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IsEnabled     bool                   `protobuf:"varint,1,opt,name=is_enabled,json=isEnabled,proto3" json:"is_enabled,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConfigStatus) Reset() {
	*x = ConfigStatus{}
	mi := &file_config_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConfigStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigStatus) ProtoMessage() {}

func (x *ConfigStatus) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigStatus.ProtoReflect.Descriptor instead.
func (*ConfigStatus) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{6}
}

func (x *ConfigStatus) GetIsEnabled() bool {
	if x != nil {
		return x.IsEnabled
	}
	return false
}

type DeleteConfigBulkRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ids           []int32                `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteConfigBulkRequest) Reset() {
	*x = DeleteConfigBulkRequest{}
	mi := &file_config_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteConfigBulkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteConfigBulkRequest) ProtoMessage() {}

func (x *DeleteConfigBulkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteConfigBulkRequest.ProtoReflect.Descriptor instead.
func (*DeleteConfigBulkRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteConfigBulkRequest) GetIds() []int32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type ReadConfigByObjectIdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ObjectId      int32                  `protobuf:"varint,1,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	DomainId      int32                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadConfigByObjectIdRequest) Reset() {
	*x = ReadConfigByObjectIdRequest{}
	mi := &file_config_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadConfigByObjectIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadConfigByObjectIdRequest) ProtoMessage() {}

func (x *ReadConfigByObjectIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadConfigByObjectIdRequest.ProtoReflect.Descriptor instead.
func (*ReadConfigByObjectIdRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{8}
}

func (x *ReadConfigByObjectIdRequest) GetObjectId() int32 {
	if x != nil {
		return x.ObjectId
	}
	return 0
}

func (x *ReadConfigByObjectIdRequest) GetDomainId() int32 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

type CheckConfigStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ObjectName    string                 `protobuf:"bytes,1,opt,name=object_name,json=objectName,proto3" json:"object_name,omitempty"`
	DomainId      int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckConfigStatusRequest) Reset() {
	*x = CheckConfigStatusRequest{}
	mi := &file_config_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckConfigStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckConfigStatusRequest) ProtoMessage() {}

func (x *CheckConfigStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckConfigStatusRequest.ProtoReflect.Descriptor instead.
func (*CheckConfigStatusRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{9}
}

func (x *CheckConfigStatusRequest) GetObjectName() string {
	if x != nil {
		return x.ObjectName
	}
	return ""
}

func (x *CheckConfigStatusRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

type ReadConfigRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigId      int32                  `protobuf:"varint,1,opt,name=config_id,json=configId,proto3" json:"config_id,omitempty"` //  int32 domainId = 8;
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadConfigRequest) Reset() {
	*x = ReadConfigRequest{}
	mi := &file_config_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadConfigRequest) ProtoMessage() {}

func (x *ReadConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadConfigRequest.ProtoReflect.Descriptor instead.
func (*ReadConfigRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{10}
}

func (x *ReadConfigRequest) GetConfigId() int32 {
	if x != nil {
		return x.ConfigId
	}
	return 0
}

type SearchConfigRequest struct {
	state  protoimpl.MessageState `protogen:"open.v1"`
	Page   int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size   int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q      string                 `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	Sort   string                 `protobuf:"bytes,4,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields []string               `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	// NOT USED
	Object        []AvailableSystemObjects `protobuf:"varint,6,rep,packed,name=object,proto3,enum=logger.AvailableSystemObjects" json:"object,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchConfigRequest) Reset() {
	*x = SearchConfigRequest{}
	mi := &file_config_service_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchConfigRequest) ProtoMessage() {}

func (x *SearchConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchConfigRequest.ProtoReflect.Descriptor instead.
func (*SearchConfigRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{11}
}

func (x *SearchConfigRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchConfigRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchConfigRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchConfigRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchConfigRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchConfigRequest) GetObject() []AvailableSystemObjects {
	if x != nil {
		return x.Object
	}
	return nil
}

type UpdateConfigRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigId      int32                  `protobuf:"varint,1,opt,name=config_id,json=configId,proto3" json:"config_id,omitempty"`
	Enabled       bool                   `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	DaysToStore   int32                  `protobuf:"varint,3,opt,name=days_to_store,json=daysToStore,proto3" json:"days_to_store,omitempty"`
	Period        int32                  `protobuf:"varint,4,opt,name=period,proto3" json:"period,omitempty"`
	Storage       *Lookup                `protobuf:"bytes,5,opt,name=storage,proto3" json:"storage,omitempty"`
	Description   string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateConfigRequest) Reset() {
	*x = UpdateConfigRequest{}
	mi := &file_config_service_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateConfigRequest) ProtoMessage() {}

func (x *UpdateConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateConfigRequest.ProtoReflect.Descriptor instead.
func (*UpdateConfigRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{12}
}

func (x *UpdateConfigRequest) GetConfigId() int32 {
	if x != nil {
		return x.ConfigId
	}
	return 0
}

func (x *UpdateConfigRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *UpdateConfigRequest) GetDaysToStore() int32 {
	if x != nil {
		return x.DaysToStore
	}
	return 0
}

func (x *UpdateConfigRequest) GetPeriod() int32 {
	if x != nil {
		return x.Period
	}
	return 0
}

func (x *UpdateConfigRequest) GetStorage() *Lookup {
	if x != nil {
		return x.Storage
	}
	return nil
}

func (x *UpdateConfigRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type PatchConfigRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigId      int32                  `protobuf:"varint,1,opt,name=config_id,json=configId,proto3" json:"config_id,omitempty"`
	Enabled       bool                   `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	DaysToStore   int32                  `protobuf:"varint,3,opt,name=days_to_store,json=daysToStore,proto3" json:"days_to_store,omitempty"`
	Period        int32                  `protobuf:"varint,4,opt,name=period,proto3" json:"period,omitempty"`
	Storage       *Lookup                `protobuf:"bytes,5,opt,name=storage,proto3" json:"storage,omitempty"`
	Description   string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Fields        []string               `protobuf:"bytes,7,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchConfigRequest) Reset() {
	*x = PatchConfigRequest{}
	mi := &file_config_service_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchConfigRequest) ProtoMessage() {}

func (x *PatchConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchConfigRequest.ProtoReflect.Descriptor instead.
func (*PatchConfigRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{13}
}

func (x *PatchConfigRequest) GetConfigId() int32 {
	if x != nil {
		return x.ConfigId
	}
	return 0
}

func (x *PatchConfigRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *PatchConfigRequest) GetDaysToStore() int32 {
	if x != nil {
		return x.DaysToStore
	}
	return 0
}

func (x *PatchConfigRequest) GetPeriod() int32 {
	if x != nil {
		return x.Period
	}
	return 0
}

func (x *PatchConfigRequest) GetStorage() *Lookup {
	if x != nil {
		return x.Storage
	}
	return nil
}

func (x *PatchConfigRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PatchConfigRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type CreateConfigRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Object        *Lookup                `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	Enabled       bool                   `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	DaysToStore   int32                  `protobuf:"varint,3,opt,name=days_to_store,json=daysToStore,proto3" json:"days_to_store,omitempty"`
	Period        int32                  `protobuf:"varint,4,opt,name=period,proto3" json:"period,omitempty"`
	Storage       *Lookup                `protobuf:"bytes,5,opt,name=storage,proto3" json:"storage,omitempty"`
	Description   string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateConfigRequest) Reset() {
	*x = CreateConfigRequest{}
	mi := &file_config_service_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateConfigRequest) ProtoMessage() {}

func (x *CreateConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_config_service_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateConfigRequest.ProtoReflect.Descriptor instead.
func (*CreateConfigRequest) Descriptor() ([]byte, []int) {
	return file_config_service_proto_rawDescGZIP(), []int{14}
}

func (x *CreateConfigRequest) GetObject() *Lookup {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *CreateConfigRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *CreateConfigRequest) GetDaysToStore() int32 {
	if x != nil {
		return x.DaysToStore
	}
	return 0
}

func (x *CreateConfigRequest) GetPeriod() int32 {
	if x != nil {
		return x.Period
	}
	return 0
}

func (x *CreateConfigRequest) GetStorage() *Lookup {
	if x != nil {
		return x.Storage
	}
	return nil
}

func (x *CreateConfigRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var File_config_service_proto protoreflect.FileDescriptor

const file_config_service_proto_rawDesc = "" +
	"\n" +
	"\x14config_service.proto\x12\x06logger\x1a\fimport.proto\x1a\x1cgoogle/api/annotations.proto\"W\n" +
	"\aConfigs\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x02 \x01(\bR\x04next\x12$\n" +
	"\x05items\x18\x03 \x03(\v2\x0e.logger.ConfigR\x05items\"5\n" +
	"\rSystemObjects\x12$\n" +
	"\x05items\x18\x03 \x03(\v2\x0e.logger.LookupR\x05items\"E\n" +
	"\x18ReadSystemObjectsRequest\x12)\n" +
	"\x10include_existing\x18\x01 \x01(\bR\x0fincludeExisting\"\a\n" +
	"\x05Empty\"\x9e\x02\n" +
	"\x06Config\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12&\n" +
	"\x06object\x18\x02 \x01(\v2\x0e.logger.LookupR\x06object\x12\x18\n" +
	"\aenabled\x18\x03 \x01(\bR\aenabled\x12\"\n" +
	"\rdays_to_store\x18\x04 \x01(\x05R\vdaysToStore\x12\x16\n" +
	"\x06period\x18\x05 \x01(\x05R\x06period\x12(\n" +
	"\astorage\x18\x06 \x01(\v2\x0e.logger.LookupR\astorage\x12 \n" +
	"\vdescription\x18\a \x01(\tR\vdescription\x12\x1b\n" +
	"\tlogs_size\x18\b \x01(\tR\blogsSize\x12\x1d\n" +
	"\n" +
	"logs_count\x18\t \x01(\x03R\tlogsCount\"2\n" +
	"\x13DeleteConfigRequest\x12\x1b\n" +
	"\tconfig_id\x18\x01 \x01(\x05R\bconfigId\"-\n" +
	"\fConfigStatus\x12\x1d\n" +
	"\n" +
	"is_enabled\x18\x01 \x01(\bR\tisEnabled\"+\n" +
	"\x17DeleteConfigBulkRequest\x12\x10\n" +
	"\x03ids\x18\x01 \x03(\x05R\x03ids\"W\n" +
	"\x1bReadConfigByObjectIdRequest\x12\x1b\n" +
	"\tobject_id\x18\x01 \x01(\x05R\bobjectId\x12\x1b\n" +
	"\tdomain_id\x18\x02 \x01(\x05R\bdomainId\"X\n" +
	"\x18CheckConfigStatusRequest\x12\x1f\n" +
	"\vobject_name\x18\x01 \x01(\tR\n" +
	"objectName\x12\x1b\n" +
	"\tdomain_id\x18\x02 \x01(\x03R\bdomainId\"0\n" +
	"\x11ReadConfigRequest\x12\x1b\n" +
	"\tconfig_id\x18\x01 \x01(\x05R\bconfigId\"\xaf\x01\n" +
	"\x13SearchConfigRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\f\n" +
	"\x01q\x18\x03 \x01(\tR\x01q\x12\x12\n" +
	"\x04sort\x18\x04 \x01(\tR\x04sort\x12\x16\n" +
	"\x06fields\x18\x05 \x03(\tR\x06fields\x126\n" +
	"\x06object\x18\x06 \x03(\x0e2\x1e.logger.AvailableSystemObjectsR\x06object\"\xd4\x01\n" +
	"\x13UpdateConfigRequest\x12\x1b\n" +
	"\tconfig_id\x18\x01 \x01(\x05R\bconfigId\x12\x18\n" +
	"\aenabled\x18\x02 \x01(\bR\aenabled\x12\"\n" +
	"\rdays_to_store\x18\x03 \x01(\x05R\vdaysToStore\x12\x16\n" +
	"\x06period\x18\x04 \x01(\x05R\x06period\x12(\n" +
	"\astorage\x18\x05 \x01(\v2\x0e.logger.LookupR\astorage\x12 \n" +
	"\vdescription\x18\x06 \x01(\tR\vdescription\"\xeb\x01\n" +
	"\x12PatchConfigRequest\x12\x1b\n" +
	"\tconfig_id\x18\x01 \x01(\x05R\bconfigId\x12\x18\n" +
	"\aenabled\x18\x02 \x01(\bR\aenabled\x12\"\n" +
	"\rdays_to_store\x18\x03 \x01(\x05R\vdaysToStore\x12\x16\n" +
	"\x06period\x18\x04 \x01(\x05R\x06period\x12(\n" +
	"\astorage\x18\x05 \x01(\v2\x0e.logger.LookupR\astorage\x12 \n" +
	"\vdescription\x18\x06 \x01(\tR\vdescription\x12\x16\n" +
	"\x06fields\x18\a \x03(\tR\x06fields\"\xdf\x01\n" +
	"\x13CreateConfigRequest\x12&\n" +
	"\x06object\x18\x01 \x01(\v2\x0e.logger.LookupR\x06object\x12\x18\n" +
	"\aenabled\x18\x02 \x01(\bR\aenabled\x12\"\n" +
	"\rdays_to_store\x18\x03 \x01(\x05R\vdaysToStore\x12\x16\n" +
	"\x06period\x18\x04 \x01(\x05R\x06period\x12(\n" +
	"\astorage\x18\x05 \x01(\v2\x0e.logger.LookupR\astorage\x12 \n" +
	"\vdescription\x18\x06 \x01(\tR\vdescription2\xcf\x06\n" +
	"\rConfigService\x12b\n" +
	"\fUpdateConfig\x12\x1b.logger.UpdateConfigRequest\x1a\x0e.logger.Config\"%\x82\xd3\xe4\x93\x02\x1f:\x01*\x1a\x1a/logger/config/{config_id}\x12`\n" +
	"\vPatchConfig\x12\x1a.logger.PatchConfigRequest\x1a\x0e.logger.Config\"%\x82\xd3\xe4\x93\x02\x1f:\x01*2\x1a/logger/config/{config_id}\x12V\n" +
	"\fCreateConfig\x12\x1b.logger.CreateConfigRequest\x1a\x0e.logger.Config\"\x19\x82\xd3\xe4\x93\x02\x13:\x01*\"\x0e/logger/config\x12^\n" +
	"\fDeleteConfig\x12\x1b.logger.DeleteConfigRequest\x1a\r.logger.Empty\"\"\x82\xd3\xe4\x93\x02\x1c*\x1a/logger/config/{config_id}\x12M\n" +
	"\x14ReadConfigByObjectId\x12#.logger.ReadConfigByObjectIdRequest\x1a\x0e.logger.Config\"\x00\x12M\n" +
	"\x11CheckConfigStatus\x12 .logger.CheckConfigStatusRequest\x1a\x14.logger.ConfigStatus\"\x00\x12o\n" +
	"\x11ReadSystemObjects\x12 .logger.ReadSystemObjectsRequest\x1a\x15.logger.SystemObjects\"!\x82\xd3\xe4\x93\x02\x1b\x12\x19/logger/available_objects\x12[\n" +
	"\n" +
	"ReadConfig\x12\x19.logger.ReadConfigRequest\x1a\x0e.logger.Config\"\"\x82\xd3\xe4\x93\x02\x1c\x12\x1a/logger/config/{config_id}\x12T\n" +
	"\fSearchConfig\x12\x1b.logger.SearchConfigRequest\x1a\x0f.logger.Configs\"\x16\x82\xd3\xe4\x93\x02\x10\x12\x0e/logger/configB~\n" +
	"\n" +
	"com.loggerB\x12ConfigServiceProtoP\x01Z$github.com/webitel/api/logger;logger\xa2\x02\x03LXX\xaa\x02\x06Logger\xca\x02\x06Logger\xe2\x02\x12Logger\\GPBMetadata\xea\x02\x06Loggerb\x06proto3"

var (
	file_config_service_proto_rawDescOnce sync.Once
	file_config_service_proto_rawDescData []byte
)

func file_config_service_proto_rawDescGZIP() []byte {
	file_config_service_proto_rawDescOnce.Do(func() {
		file_config_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_config_service_proto_rawDesc), len(file_config_service_proto_rawDesc)))
	})
	return file_config_service_proto_rawDescData
}

var file_config_service_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_config_service_proto_goTypes = []any{
	(*Configs)(nil),                     // 0: logger.Configs
	(*SystemObjects)(nil),               // 1: logger.SystemObjects
	(*ReadSystemObjectsRequest)(nil),    // 2: logger.ReadSystemObjectsRequest
	(*Empty)(nil),                       // 3: logger.Empty
	(*Config)(nil),                      // 4: logger.Config
	(*DeleteConfigRequest)(nil),         // 5: logger.DeleteConfigRequest
	(*ConfigStatus)(nil),                // 6: logger.ConfigStatus
	(*DeleteConfigBulkRequest)(nil),     // 7: logger.DeleteConfigBulkRequest
	(*ReadConfigByObjectIdRequest)(nil), // 8: logger.ReadConfigByObjectIdRequest
	(*CheckConfigStatusRequest)(nil),    // 9: logger.CheckConfigStatusRequest
	(*ReadConfigRequest)(nil),           // 10: logger.ReadConfigRequest
	(*SearchConfigRequest)(nil),         // 11: logger.SearchConfigRequest
	(*UpdateConfigRequest)(nil),         // 12: logger.UpdateConfigRequest
	(*PatchConfigRequest)(nil),          // 13: logger.PatchConfigRequest
	(*CreateConfigRequest)(nil),         // 14: logger.CreateConfigRequest
	(*Lookup)(nil),                      // 15: logger.Lookup
	(AvailableSystemObjects)(0),         // 16: logger.AvailableSystemObjects
}
var file_config_service_proto_depIdxs = []int32{
	4,  // 0: logger.Configs.items:type_name -> logger.Config
	15, // 1: logger.SystemObjects.items:type_name -> logger.Lookup
	15, // 2: logger.Config.object:type_name -> logger.Lookup
	15, // 3: logger.Config.storage:type_name -> logger.Lookup
	16, // 4: logger.SearchConfigRequest.object:type_name -> logger.AvailableSystemObjects
	15, // 5: logger.UpdateConfigRequest.storage:type_name -> logger.Lookup
	15, // 6: logger.PatchConfigRequest.storage:type_name -> logger.Lookup
	15, // 7: logger.CreateConfigRequest.object:type_name -> logger.Lookup
	15, // 8: logger.CreateConfigRequest.storage:type_name -> logger.Lookup
	12, // 9: logger.ConfigService.UpdateConfig:input_type -> logger.UpdateConfigRequest
	13, // 10: logger.ConfigService.PatchConfig:input_type -> logger.PatchConfigRequest
	14, // 11: logger.ConfigService.CreateConfig:input_type -> logger.CreateConfigRequest
	5,  // 12: logger.ConfigService.DeleteConfig:input_type -> logger.DeleteConfigRequest
	8,  // 13: logger.ConfigService.ReadConfigByObjectId:input_type -> logger.ReadConfigByObjectIdRequest
	9,  // 14: logger.ConfigService.CheckConfigStatus:input_type -> logger.CheckConfigStatusRequest
	2,  // 15: logger.ConfigService.ReadSystemObjects:input_type -> logger.ReadSystemObjectsRequest
	10, // 16: logger.ConfigService.ReadConfig:input_type -> logger.ReadConfigRequest
	11, // 17: logger.ConfigService.SearchConfig:input_type -> logger.SearchConfigRequest
	4,  // 18: logger.ConfigService.UpdateConfig:output_type -> logger.Config
	4,  // 19: logger.ConfigService.PatchConfig:output_type -> logger.Config
	4,  // 20: logger.ConfigService.CreateConfig:output_type -> logger.Config
	3,  // 21: logger.ConfigService.DeleteConfig:output_type -> logger.Empty
	4,  // 22: logger.ConfigService.ReadConfigByObjectId:output_type -> logger.Config
	6,  // 23: logger.ConfigService.CheckConfigStatus:output_type -> logger.ConfigStatus
	1,  // 24: logger.ConfigService.ReadSystemObjects:output_type -> logger.SystemObjects
	4,  // 25: logger.ConfigService.ReadConfig:output_type -> logger.Config
	0,  // 26: logger.ConfigService.SearchConfig:output_type -> logger.Configs
	18, // [18:27] is the sub-list for method output_type
	9,  // [9:18] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_config_service_proto_init() }
func file_config_service_proto_init() {
	if File_config_service_proto != nil {
		return
	}
	file_import_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_config_service_proto_rawDesc), len(file_config_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_config_service_proto_goTypes,
		DependencyIndexes: file_config_service_proto_depIdxs,
		MessageInfos:      file_config_service_proto_msgTypes,
	}.Build()
	File_config_service_proto = out.File
	file_config_service_proto_goTypes = nil
	file_config_service_proto_depIdxs = nil
}
