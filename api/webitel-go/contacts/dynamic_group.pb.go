// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: contacts/dynamic_group.proto

package contacts

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/visibility"
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

// Dynamic Group condition
type DynamicCondition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the condition.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The query or condition expression used to evaluate the group.
	Expression string `protobuf:"bytes,2,opt,name=expression,proto3" json:"expression,omitempty"`
	// The ID of the static group that should be assigned if the condition is met.
	Group *Lookup `protobuf:"bytes,3,opt,name=group,proto3" json:"group,omitempty"`
	// The ID of the assignee that should be assigned if the condition is met (optional).
	Assignee      *Lookup `protobuf:"bytes,4,opt,name=assignee,proto3" json:"assignee,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DynamicCondition) Reset() {
	*x = DynamicCondition{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DynamicCondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DynamicCondition) ProtoMessage() {}

func (x *DynamicCondition) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DynamicCondition.ProtoReflect.Descriptor instead.
func (*DynamicCondition) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{0}
}

func (x *DynamicCondition) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DynamicCondition) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *DynamicCondition) GetGroup() *Lookup {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *DynamicCondition) GetAssignee() *Lookup {
	if x != nil {
		return x.Assignee
	}
	return nil
}

// Dynamic Group condition
type InputDynamicCondition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The query or condition expression used to evaluate the group.
	Expression string `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
	// The ID of the static group that should be assigned if the condition is met.
	Group *Lookup `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"`
	// The ID of the assignee that should be assigned if the condition is met (optional).
	Assignee      *Lookup `protobuf:"bytes,3,opt,name=assignee,proto3" json:"assignee,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputDynamicCondition) Reset() {
	*x = InputDynamicCondition{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputDynamicCondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputDynamicCondition) ProtoMessage() {}

func (x *InputDynamicCondition) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputDynamicCondition.ProtoReflect.Descriptor instead.
func (*InputDynamicCondition) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{1}
}

func (x *InputDynamicCondition) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *InputDynamicCondition) GetGroup() *Lookup {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *InputDynamicCondition) GetAssignee() *Lookup {
	if x != nil {
		return x.Assignee
	}
	return nil
}

// Dynamic Group
type DynamicGroup struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the dynamic group. Never changes.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the dynamic group.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the dynamic group.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// Timestamp(milli) of the group's creation.
	CreatedAt int64 `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Timestamp(milli) of the last group update.
	UpdatedAt int64 `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The user who created this dynamic group.
	CreatedBy *Lookup `protobuf:"bytes,7,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// The user who performed the last update.
	UpdatedBy *Lookup `protobuf:"bytes,8,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// List of dynamic conditions associated with the group.
	Conditions []*DynamicCondition `protobuf:"bytes,9,rep,name=conditions,proto3" json:"conditions,omitempty"`
	// Default static group to be assigned if no conditions are met.
	DefaultGroup *Lookup `protobuf:"bytes,10,opt,name=default_group,json=defaultGroup,proto3" json:"default_group,omitempty"`
	// Enabled status of the group: active or inactive.
	Enabled       bool `protobuf:"varint,11,opt,name=enabled,proto3" json:"enabled,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DynamicGroup) Reset() {
	*x = DynamicGroup{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DynamicGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DynamicGroup) ProtoMessage() {}

func (x *DynamicGroup) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DynamicGroup.ProtoReflect.Descriptor instead.
func (*DynamicGroup) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{2}
}

func (x *DynamicGroup) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DynamicGroup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DynamicGroup) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *DynamicGroup) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *DynamicGroup) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *DynamicGroup) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *DynamicGroup) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *DynamicGroup) GetConditions() []*DynamicCondition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

func (x *DynamicGroup) GetDefaultGroup() *Lookup {
	if x != nil {
		return x.DefaultGroup
	}
	return nil
}

func (x *DynamicGroup) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

// A list of Dynamic Groups.
type DynamicGroupList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of the partial result.
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Have more records.
	Next bool `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	// List of dynamic groups.
	Items         []*DynamicGroup `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DynamicGroupList) Reset() {
	*x = DynamicGroupList{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DynamicGroupList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DynamicGroupList) ProtoMessage() {}

func (x *DynamicGroupList) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DynamicGroupList.ProtoReflect.Descriptor instead.
func (*DynamicGroupList) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{3}
}

func (x *DynamicGroupList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *DynamicGroupList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *DynamicGroupList) GetItems() []*DynamicGroup {
	if x != nil {
		return x.Items
	}
	return nil
}

// Input message for creating/updating a dynamic group.
type DynamicGroupInput struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the dynamic group.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the dynamic group.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Default static group to assign if no conditions are met.
	DefaultGroup *Lookup `protobuf:"bytes,3,opt,name=default_group,json=defaultGroup,proto3" json:"default_group,omitempty"`
	// Enabled status of the dynamic group: active/inactive.
	Enabled       bool `protobuf:"varint,6,opt,name=enabled,proto3" json:"enabled,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DynamicGroupInput) Reset() {
	*x = DynamicGroupInput{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DynamicGroupInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DynamicGroupInput) ProtoMessage() {}

func (x *DynamicGroupInput) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DynamicGroupInput.ProtoReflect.Descriptor instead.
func (*DynamicGroupInput) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{4}
}

func (x *DynamicGroupInput) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DynamicGroupInput) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *DynamicGroupInput) GetDefaultGroup() *Lookup {
	if x != nil {
		return x.DefaultGroup
	}
	return nil
}

func (x *DynamicGroupInput) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

// Request message for creating a new dynamic group.
type CreateDynamicGroupRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the dynamic group.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the dynamic group.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Default static group to assign if no conditions are met.
	DefaultGroup *Lookup `protobuf:"bytes,3,opt,name=default_group,json=defaultGroup,proto3" json:"default_group,omitempty"`
	// Enabled status of the dynamic group: active/inactive.
	Enabled bool `protobuf:"varint,4,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// Input details for the dynamic group.
	Condition     []*InputDynamicCondition `protobuf:"bytes,5,rep,name=condition,proto3" json:"condition,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDynamicGroupRequest) Reset() {
	*x = CreateDynamicGroupRequest{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDynamicGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDynamicGroupRequest) ProtoMessage() {}

func (x *CreateDynamicGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDynamicGroupRequest.ProtoReflect.Descriptor instead.
func (*CreateDynamicGroupRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{5}
}

func (x *CreateDynamicGroupRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateDynamicGroupRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateDynamicGroupRequest) GetDefaultGroup() *Lookup {
	if x != nil {
		return x.DefaultGroup
	}
	return nil
}

func (x *CreateDynamicGroupRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *CreateDynamicGroupRequest) GetCondition() []*InputDynamicCondition {
	if x != nil {
		return x.Condition
	}
	return nil
}

// Request message for locating a dynamic group by ID.
type LocateDynamicGroupRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Fields        []string               `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateDynamicGroupRequest) Reset() {
	*x = LocateDynamicGroupRequest{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateDynamicGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateDynamicGroupRequest) ProtoMessage() {}

func (x *LocateDynamicGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateDynamicGroupRequest.ProtoReflect.Descriptor instead.
func (*LocateDynamicGroupRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{6}
}

func (x *LocateDynamicGroupRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateDynamicGroupRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Response message for locating a dynamic group by ID.
type LocateDynamicGroupResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Group         *DynamicGroup          `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateDynamicGroupResponse) Reset() {
	*x = LocateDynamicGroupResponse{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateDynamicGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateDynamicGroupResponse) ProtoMessage() {}

func (x *LocateDynamicGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateDynamicGroupResponse.ProtoReflect.Descriptor instead.
func (*LocateDynamicGroupResponse) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{7}
}

func (x *LocateDynamicGroupResponse) GetGroup() *DynamicGroup {
	if x != nil {
		return x.Group
	}
	return nil
}

// Request message for updating an existing dynamic group.
type UpdateDynamicGroupRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the dynamic group to update.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Input details for the dynamic group.
	Input *DynamicGroupInput `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,3,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDynamicGroupRequest) Reset() {
	*x = UpdateDynamicGroupRequest{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDynamicGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDynamicGroupRequest) ProtoMessage() {}

func (x *UpdateDynamicGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDynamicGroupRequest.ProtoReflect.Descriptor instead.
func (*UpdateDynamicGroupRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateDynamicGroupRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateDynamicGroupRequest) GetInput() *DynamicGroupInput {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateDynamicGroupRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// Request message for deleting a dynamic group.
type DeleteDynamicGroupRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the dynamic group to delete.
	Id            int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteDynamicGroupRequest) Reset() {
	*x = DeleteDynamicGroupRequest{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteDynamicGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDynamicGroupRequest) ProtoMessage() {}

func (x *DeleteDynamicGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDynamicGroupRequest.ProtoReflect.Descriptor instead.
func (*DeleteDynamicGroupRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteDynamicGroupRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// Request message for listing dynamic groups.
type ListDynamicGroupsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of result dataset records. offset = (page*size)
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Size count of records on result page. limit = (size++)
	Size int32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Sort the result according to fields.
	Sort []string `protobuf:"bytes,4,rep,name=sort,proto3" json:"sort,omitempty"`
	// Filter by unique IDs.
	Id []int64 `protobuf:"varint,5,rep,packed,name=id,proto3" json:"id,omitempty"`
	// Search term: group name;
	// `?` - matches any one character
	// `*` - matches 0 or more characters
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Filter by group name.
	Name          string `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDynamicGroupsRequest) Reset() {
	*x = ListDynamicGroupsRequest{}
	mi := &file_contacts_dynamic_group_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDynamicGroupsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDynamicGroupsRequest) ProtoMessage() {}

func (x *ListDynamicGroupsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_group_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDynamicGroupsRequest.ProtoReflect.Descriptor instead.
func (*ListDynamicGroupsRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_group_proto_rawDescGZIP(), []int{10}
}

func (x *ListDynamicGroupsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListDynamicGroupsRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListDynamicGroupsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListDynamicGroupsRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListDynamicGroupsRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListDynamicGroupsRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListDynamicGroupsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_contacts_dynamic_group_proto protoreflect.FileDescriptor

const file_contacts_dynamic_group_proto_rawDesc = "" +
	"\n" +
	"\x1ccontacts/dynamic_group.proto\x12\x10webitel.contacts\x1a\x1bgoogle/api/visibility.proto\x1a\x15contacts/fields.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\xa8\x01\n" +
	"\x10DynamicCondition\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x1e\n" +
	"\n" +
	"expression\x18\x02 \x01(\tR\n" +
	"expression\x12.\n" +
	"\x05group\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\x05group\x124\n" +
	"\bassignee\x18\x04 \x01(\v2\x18.webitel.contacts.LookupR\bassignee\"\x9d\x01\n" +
	"\x15InputDynamicCondition\x12\x1e\n" +
	"\n" +
	"expression\x18\x01 \x01(\tR\n" +
	"expression\x12.\n" +
	"\x05group\x18\x02 \x01(\v2\x18.webitel.contacts.LookupR\x05group\x124\n" +
	"\bassignee\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\bassignee\"\xa1\x03\n" +
	"\fDynamicGroup\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x04 \x01(\tR\vdescription\x12\x1d\n" +
	"\n" +
	"created_at\x18\x05 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x06 \x01(\x03R\tupdatedAt\x127\n" +
	"\n" +
	"created_by\x18\a \x01(\v2\x18.webitel.contacts.LookupR\tcreatedBy\x127\n" +
	"\n" +
	"updated_by\x18\b \x01(\v2\x18.webitel.contacts.LookupR\tupdatedBy\x12B\n" +
	"\n" +
	"conditions\x18\t \x03(\v2\".webitel.contacts.DynamicConditionR\n" +
	"conditions\x12=\n" +
	"\rdefault_group\x18\n" +
	" \x01(\v2\x18.webitel.contacts.LookupR\fdefaultGroup\x12\x18\n" +
	"\aenabled\x18\v \x01(\bR\aenabled\"p\n" +
	"\x10DynamicGroupList\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x02 \x01(\bR\x04next\x124\n" +
	"\x05items\x18\x03 \x03(\v2\x1e.webitel.contacts.DynamicGroupR\x05items\"\xa2\x01\n" +
	"\x11DynamicGroupInput\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12=\n" +
	"\rdefault_group\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\fdefaultGroup\x12\x18\n" +
	"\aenabled\x18\x06 \x01(\bR\aenabled\"\xff\x01\n" +
	"\x19CreateDynamicGroupRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12=\n" +
	"\rdefault_group\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\fdefaultGroup\x12\x18\n" +
	"\aenabled\x18\x04 \x01(\bR\aenabled\x12E\n" +
	"\tcondition\x18\x05 \x03(\v2'.webitel.contacts.InputDynamicConditionR\tcondition:\f\x92A\t\n" +
	"\a\xd2\x01\x04name\"C\n" +
	"\x19LocateDynamicGroupRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x16\n" +
	"\x06fields\x18\x02 \x03(\tR\x06fields\"R\n" +
	"\x1aLocateDynamicGroupResponse\x124\n" +
	"\x05group\x18\x01 \x01(\v2\x1e.webitel.contacts.DynamicGroupR\x05group\"\xad\x01\n" +
	"\x19UpdateDynamicGroupRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x129\n" +
	"\x05input\x18\x02 \x01(\v2#.webitel.contacts.DynamicGroupInputR\x05input\x129\n" +
	"\vx_json_mask\x18\x03 \x03(\tB\x19\x92A\a@\x01\x8a\x01\x02^$\xfa\xd2\xe4\x93\x02\t\x12\aPREVIEWR\txJsonMask:\n" +
	"\x92A\a\n" +
	"\x05\xd2\x01\x02id\"7\n" +
	"\x19DeleteDynamicGroupRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id:\n" +
	"\x92A\a\n" +
	"\x05\xd2\x01\x02id\"\xa0\x01\n" +
	"\x18ListDynamicGroupsRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\x12\x12\n" +
	"\x04sort\x18\x04 \x03(\tR\x04sort\x12\x0e\n" +
	"\x02id\x18\x05 \x03(\x03R\x02id\x12\f\n" +
	"\x01q\x18\x06 \x01(\tR\x01q\x12\x12\n" +
	"\x04name\x18\a \x01(\tR\x04name2\xbd\a\n" +
	"\rDynamicGroups\x12\xc4\x01\n" +
	"\x11ListDynamicGroups\x12*.webitel.contacts.ListDynamicGroupsRequest\x1a\".webitel.contacts.DynamicGroupList\"_\x92A<\x12:Retrieve a list of dynamic groups or search dynamic groups\x82\xd3\xe4\x93\x02\x1a\x12\x18/contacts/groups/dynamic\x12\xa5\x01\n" +
	"\x12CreateDynamicGroup\x12+.webitel.contacts.CreateDynamicGroupRequest\x1a\x1e.webitel.contacts.DynamicGroup\"B\x92A\x1c\x12\x1aCreate a new dynamic group\x82\xd3\xe4\x93\x02\x1d:\x01*\"\x18/contacts/groups/dynamic\x12\xdc\x01\n" +
	"\x12UpdateDynamicGroup\x12+.webitel.contacts.UpdateDynamicGroupRequest\x1a\x1e.webitel.contacts.DynamicGroup\"y\x92A\"\x12 Update an existing dynamic group\x82\xd3\xe4\x93\x02N:\x05inputZ&:\x05input2\x1d/contacts/groups/{id}/dynamic\x1a\x1d/contacts/groups/{id}/dynamic\x12\xa3\x01\n" +
	"\x12DeleteDynamicGroup\x12+.webitel.contacts.DeleteDynamicGroupRequest\x1a\x1e.webitel.contacts.DynamicGroup\"@\x92A\x18\x12\x16Delete a dynamic group\x82\xd3\xe4\x93\x02\x1f*\x1d/contacts/groups/{id}/dynamic\x12\xb7\x01\n" +
	"\x12LocateDynamicGroup\x12+.webitel.contacts.LocateDynamicGroupRequest\x1a,.webitel.contacts.LocateDynamicGroupResponse\"F\x92A\x1e\x12\x1cLocate a dynamic group by ID\x82\xd3\xe4\x93\x02\x1f\x12\x1d/contacts/groups/{id}/dynamicB\xac\x01\n" +
	"\x14com.webitel.contactsB\x11DynamicGroupProtoP\x01Z webitel.go/api/contacts;contacts\xa2\x02\x03WCX\xaa\x02\x10Webitel.Contacts\xca\x02\x10Webitel\\Contacts\xe2\x02\x1cWebitel\\Contacts\\GPBMetadata\xea\x02\x11Webitel::Contactsb\x06proto3"

var (
	file_contacts_dynamic_group_proto_rawDescOnce sync.Once
	file_contacts_dynamic_group_proto_rawDescData []byte
)

func file_contacts_dynamic_group_proto_rawDescGZIP() []byte {
	file_contacts_dynamic_group_proto_rawDescOnce.Do(func() {
		file_contacts_dynamic_group_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contacts_dynamic_group_proto_rawDesc), len(file_contacts_dynamic_group_proto_rawDesc)))
	})
	return file_contacts_dynamic_group_proto_rawDescData
}

var file_contacts_dynamic_group_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_contacts_dynamic_group_proto_goTypes = []any{
	(*DynamicCondition)(nil),           // 0: webitel.contacts.DynamicCondition
	(*InputDynamicCondition)(nil),      // 1: webitel.contacts.InputDynamicCondition
	(*DynamicGroup)(nil),               // 2: webitel.contacts.DynamicGroup
	(*DynamicGroupList)(nil),           // 3: webitel.contacts.DynamicGroupList
	(*DynamicGroupInput)(nil),          // 4: webitel.contacts.DynamicGroupInput
	(*CreateDynamicGroupRequest)(nil),  // 5: webitel.contacts.CreateDynamicGroupRequest
	(*LocateDynamicGroupRequest)(nil),  // 6: webitel.contacts.LocateDynamicGroupRequest
	(*LocateDynamicGroupResponse)(nil), // 7: webitel.contacts.LocateDynamicGroupResponse
	(*UpdateDynamicGroupRequest)(nil),  // 8: webitel.contacts.UpdateDynamicGroupRequest
	(*DeleteDynamicGroupRequest)(nil),  // 9: webitel.contacts.DeleteDynamicGroupRequest
	(*ListDynamicGroupsRequest)(nil),   // 10: webitel.contacts.ListDynamicGroupsRequest
	(*Lookup)(nil),                     // 11: webitel.contacts.Lookup
}
var file_contacts_dynamic_group_proto_depIdxs = []int32{
	11, // 0: webitel.contacts.DynamicCondition.group:type_name -> webitel.contacts.Lookup
	11, // 1: webitel.contacts.DynamicCondition.assignee:type_name -> webitel.contacts.Lookup
	11, // 2: webitel.contacts.InputDynamicCondition.group:type_name -> webitel.contacts.Lookup
	11, // 3: webitel.contacts.InputDynamicCondition.assignee:type_name -> webitel.contacts.Lookup
	11, // 4: webitel.contacts.DynamicGroup.created_by:type_name -> webitel.contacts.Lookup
	11, // 5: webitel.contacts.DynamicGroup.updated_by:type_name -> webitel.contacts.Lookup
	0,  // 6: webitel.contacts.DynamicGroup.conditions:type_name -> webitel.contacts.DynamicCondition
	11, // 7: webitel.contacts.DynamicGroup.default_group:type_name -> webitel.contacts.Lookup
	2,  // 8: webitel.contacts.DynamicGroupList.items:type_name -> webitel.contacts.DynamicGroup
	11, // 9: webitel.contacts.DynamicGroupInput.default_group:type_name -> webitel.contacts.Lookup
	11, // 10: webitel.contacts.CreateDynamicGroupRequest.default_group:type_name -> webitel.contacts.Lookup
	1,  // 11: webitel.contacts.CreateDynamicGroupRequest.condition:type_name -> webitel.contacts.InputDynamicCondition
	2,  // 12: webitel.contacts.LocateDynamicGroupResponse.group:type_name -> webitel.contacts.DynamicGroup
	4,  // 13: webitel.contacts.UpdateDynamicGroupRequest.input:type_name -> webitel.contacts.DynamicGroupInput
	10, // 14: webitel.contacts.DynamicGroups.ListDynamicGroups:input_type -> webitel.contacts.ListDynamicGroupsRequest
	5,  // 15: webitel.contacts.DynamicGroups.CreateDynamicGroup:input_type -> webitel.contacts.CreateDynamicGroupRequest
	8,  // 16: webitel.contacts.DynamicGroups.UpdateDynamicGroup:input_type -> webitel.contacts.UpdateDynamicGroupRequest
	9,  // 17: webitel.contacts.DynamicGroups.DeleteDynamicGroup:input_type -> webitel.contacts.DeleteDynamicGroupRequest
	6,  // 18: webitel.contacts.DynamicGroups.LocateDynamicGroup:input_type -> webitel.contacts.LocateDynamicGroupRequest
	3,  // 19: webitel.contacts.DynamicGroups.ListDynamicGroups:output_type -> webitel.contacts.DynamicGroupList
	2,  // 20: webitel.contacts.DynamicGroups.CreateDynamicGroup:output_type -> webitel.contacts.DynamicGroup
	2,  // 21: webitel.contacts.DynamicGroups.UpdateDynamicGroup:output_type -> webitel.contacts.DynamicGroup
	2,  // 22: webitel.contacts.DynamicGroups.DeleteDynamicGroup:output_type -> webitel.contacts.DynamicGroup
	7,  // 23: webitel.contacts.DynamicGroups.LocateDynamicGroup:output_type -> webitel.contacts.LocateDynamicGroupResponse
	19, // [19:24] is the sub-list for method output_type
	14, // [14:19] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_contacts_dynamic_group_proto_init() }
func file_contacts_dynamic_group_proto_init() {
	if File_contacts_dynamic_group_proto != nil {
		return
	}
	file_contacts_fields_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contacts_dynamic_group_proto_rawDesc), len(file_contacts_dynamic_group_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contacts_dynamic_group_proto_goTypes,
		DependencyIndexes: file_contacts_dynamic_group_proto_depIdxs,
		MessageInfos:      file_contacts_dynamic_group_proto_msgTypes,
	}.Build()
	File_contacts_dynamic_group_proto = out.File
	file_contacts_dynamic_group_proto_goTypes = nil
	file_contacts_dynamic_group_proto_depIdxs = nil
}
