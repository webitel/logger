// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: contacts/dynamic_condition.proto

package contacts

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/visibility"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
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

// Condition
type Condition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the condition.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The query or condition expression used to evaluate the group.
	Expression string `protobuf:"bytes,2,opt,name=expression,proto3" json:"expression,omitempty"`
	// The ID of the static group that should be assigned if the condition is met.
	Group *Lookup `protobuf:"bytes,3,opt,name=group,proto3" json:"group,omitempty"`
	// The ID of the assignee that should be assigned if the condition is met (optional).
	Assignee *Lookup `protobuf:"bytes,4,opt,name=assignee,proto3" json:"assignee,omitempty"`
	// The user who created this condition.
	CreatedBy *Lookup `protobuf:"bytes,5,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// The user who performed the last update.
	UpdatedBy *Lookup `protobuf:"bytes,6,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// The timestamp (in milliseconds) of when the condition was created.
	CreatedAt int64 `protobuf:"varint,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The timestamp (in milliseconds) of the last update.
	UpdatedAt     int64 `protobuf:"varint,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Condition) Reset() {
	*x = Condition{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Condition) ProtoMessage() {}

func (x *Condition) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Condition.ProtoReflect.Descriptor instead.
func (*Condition) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{0}
}

func (x *Condition) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Condition) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *Condition) GetGroup() *Lookup {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *Condition) GetAssignee() *Lookup {
	if x != nil {
		return x.Assignee
	}
	return nil
}

func (x *Condition) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Condition) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *Condition) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Condition) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// Request message for listing conditions.
type ListConditionsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of result dataset records. offset = (page*size)
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Size count of records on result page. limit = (size++)
	Size int32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Sort the result according to fields.
	Sort []string `protobuf:"bytes,4,rep,name=sort,proto3" json:"sort,omitempty"`
	// Search term for conditions.
	Q string `protobuf:"bytes,5,opt,name=q,proto3" json:"q,omitempty"`
	// Filter by unique IDs.
	Id []int64 `protobuf:"varint,6,rep,packed,name=id,proto3" json:"id,omitempty"`
	// The ID of the group to which the conditions belong.
	GroupId       int64 `protobuf:"varint,7,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListConditionsRequest) Reset() {
	*x = ListConditionsRequest{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListConditionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListConditionsRequest) ProtoMessage() {}

func (x *ListConditionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListConditionsRequest.ProtoReflect.Descriptor instead.
func (*ListConditionsRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{1}
}

func (x *ListConditionsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListConditionsRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListConditionsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListConditionsRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListConditionsRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListConditionsRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListConditionsRequest) GetGroupId() int64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

// A list of Conditions.
type ConditionList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of the partial result.
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Have more records.
	Next bool `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	// List of conditions.
	Items         []*Condition `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConditionList) Reset() {
	*x = ConditionList{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConditionList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConditionList) ProtoMessage() {}

func (x *ConditionList) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConditionList.ProtoReflect.Descriptor instead.
func (*ConditionList) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{2}
}

func (x *ConditionList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ConditionList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ConditionList) GetItems() []*Condition {
	if x != nil {
		return x.Items
	}
	return nil
}

// Request message for creating a new condition.
type CreateConditionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The ID of the group to which the condition belongs.
	GroupId int64 `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	// The query or condition expression used to evaluate the group.
	Expression string `protobuf:"bytes,2,opt,name=expression,proto3" json:"expression,omitempty"`
	// The ID of the static group that should be assigned if the condition is met.
	Group *Lookup `protobuf:"bytes,3,opt,name=group,proto3" json:"group,omitempty"`
	// The ID of the assignee that should be assigned if the condition is met (optional).
	Assignee      *Lookup `protobuf:"bytes,4,opt,name=assignee,proto3" json:"assignee,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateConditionRequest) Reset() {
	*x = CreateConditionRequest{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateConditionRequest) ProtoMessage() {}

func (x *CreateConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateConditionRequest.ProtoReflect.Descriptor instead.
func (*CreateConditionRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{3}
}

func (x *CreateConditionRequest) GetGroupId() int64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *CreateConditionRequest) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *CreateConditionRequest) GetGroup() *Lookup {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *CreateConditionRequest) GetAssignee() *Lookup {
	if x != nil {
		return x.Assignee
	}
	return nil
}

// Position details for conditions in the group.
type Position struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The current position of the condition in the list.
	// if we set zero (0) index position -- set only cond_down -- cond_up should be ZERO
	CondUp int64 `protobuf:"varint,1,opt,name=cond_up,json=condUp,proto3" json:"cond_up,omitempty"`
	// The target position where the condition should be moved.
	// if we set last (n) index position -- set only cond_up -- cond_down should be ZERO
	CondDown      int64 `protobuf:"varint,2,opt,name=cond_down,json=condDown,proto3" json:"cond_down,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Position) Reset() {
	*x = Position{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Position) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Position) ProtoMessage() {}

func (x *Position) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Position.ProtoReflect.Descriptor instead.
func (*Position) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{4}
}

func (x *Position) GetCondUp() int64 {
	if x != nil {
		return x.CondUp
	}
	return 0
}

func (x *Position) GetCondDown() int64 {
	if x != nil {
		return x.CondDown
	}
	return 0
}

// Input message for creating/updating a condition.
type InputCondition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The query or condition expression used to evaluate the group.
	Expression string `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
	// The ID of the static group that should be assigned if the condition is met.
	Group int64 `protobuf:"varint,2,opt,name=group,proto3" json:"group,omitempty"`
	// The ID of the assignee that should be assigned if the condition is met (optional).
	Assignee *Lookup `protobuf:"bytes,3,opt,name=assignee,proto3" json:"assignee,omitempty"`
	// The position of the condition in the group.
	Position      *Position `protobuf:"bytes,4,opt,name=position,proto3" json:"position,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputCondition) Reset() {
	*x = InputCondition{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputCondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputCondition) ProtoMessage() {}

func (x *InputCondition) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputCondition.ProtoReflect.Descriptor instead.
func (*InputCondition) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{5}
}

func (x *InputCondition) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *InputCondition) GetGroup() int64 {
	if x != nil {
		return x.Group
	}
	return 0
}

func (x *InputCondition) GetAssignee() *Lookup {
	if x != nil {
		return x.Assignee
	}
	return nil
}

func (x *InputCondition) GetPosition() *Position {
	if x != nil {
		return x.Position
	}
	return nil
}

// Request message for locating a condition by ID.
type LocateConditionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // Unique ID of the condition.
	// Fields to be retrieved as a result.
	Fields        []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateConditionRequest) Reset() {
	*x = LocateConditionRequest{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateConditionRequest) ProtoMessage() {}

func (x *LocateConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateConditionRequest.ProtoReflect.Descriptor instead.
func (*LocateConditionRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{6}
}

func (x *LocateConditionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateConditionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Response message for locating a condition by ID.
type LocateConditionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Condition     *Condition             `protobuf:"bytes,1,opt,name=condition,proto3" json:"condition,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateConditionResponse) Reset() {
	*x = LocateConditionResponse{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateConditionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateConditionResponse) ProtoMessage() {}

func (x *LocateConditionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateConditionResponse.ProtoReflect.Descriptor instead.
func (*LocateConditionResponse) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{7}
}

func (x *LocateConditionResponse) GetCondition() *Condition {
	if x != nil {
		return x.Condition
	}
	return nil
}

// Request message for updating an existing condition.
type UpdateConditionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the condition to update.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Input details for the condition.
	Input *InputCondition `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,3,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateConditionRequest) Reset() {
	*x = UpdateConditionRequest{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateConditionRequest) ProtoMessage() {}

func (x *UpdateConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateConditionRequest.ProtoReflect.Descriptor instead.
func (*UpdateConditionRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateConditionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateConditionRequest) GetInput() *InputCondition {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateConditionRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// Request message for deleting a condition.
type DeleteConditionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the condition to delete.
	Id            int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // Unique ID of the condition.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteConditionRequest) Reset() {
	*x = DeleteConditionRequest{}
	mi := &file_contacts_dynamic_condition_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteConditionRequest) ProtoMessage() {}

func (x *DeleteConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_dynamic_condition_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteConditionRequest.ProtoReflect.Descriptor instead.
func (*DeleteConditionRequest) Descriptor() ([]byte, []int) {
	return file_contacts_dynamic_condition_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteConditionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_contacts_dynamic_condition_proto protoreflect.FileDescriptor

const file_contacts_dynamic_condition_proto_rawDesc = "" +
	"\n" +
	" contacts/dynamic_condition.proto\x12\x10webitel.contacts\x1a\x1bgoogle/api/visibility.proto\x1a\x15contacts/fields.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\x1a\x1egoogle/protobuf/wrappers.proto\"\xd1\x02\n" +
	"\tCondition\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x1e\n" +
	"\n" +
	"expression\x18\x02 \x01(\tR\n" +
	"expression\x12.\n" +
	"\x05group\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\x05group\x124\n" +
	"\bassignee\x18\x04 \x01(\v2\x18.webitel.contacts.LookupR\bassignee\x127\n" +
	"\n" +
	"created_by\x18\x05 \x01(\v2\x18.webitel.contacts.LookupR\tcreatedBy\x127\n" +
	"\n" +
	"updated_by\x18\x06 \x01(\v2\x18.webitel.contacts.LookupR\tupdatedBy\x12\x1d\n" +
	"\n" +
	"created_at\x18\a \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\b \x01(\x03R\tupdatedAt\"\xb6\x01\n" +
	"\x15ListConditionsRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\x12\x12\n" +
	"\x04sort\x18\x04 \x03(\tR\x04sort\x12\f\n" +
	"\x01q\x18\x05 \x01(\tR\x01q\x12\x0e\n" +
	"\x02id\x18\x06 \x03(\x03R\x02id\x12\x19\n" +
	"\bgroup_id\x18\a \x01(\x03R\agroupId:\x10\x92A\r\n" +
	"\v\xd2\x01\bgroup_id\"j\n" +
	"\rConditionList\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x02 \x01(\bR\x04next\x121\n" +
	"\x05items\x18\x03 \x03(\v2\x1b.webitel.contacts.ConditionR\x05items\"\xcb\x01\n" +
	"\x16CreateConditionRequest\x12\x19\n" +
	"\bgroup_id\x18\x01 \x01(\x03R\agroupId\x12\x1e\n" +
	"\n" +
	"expression\x18\x02 \x01(\tR\n" +
	"expression\x12.\n" +
	"\x05group\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\x05group\x124\n" +
	"\bassignee\x18\x04 \x01(\v2\x18.webitel.contacts.LookupR\bassignee:\x10\x92A\r\n" +
	"\v\xd2\x01\bgroup_id\"@\n" +
	"\bPosition\x12\x17\n" +
	"\acond_up\x18\x01 \x01(\x03R\x06condUp\x12\x1b\n" +
	"\tcond_down\x18\x02 \x01(\x03R\bcondDown\"\xb4\x01\n" +
	"\x0eInputCondition\x12\x1e\n" +
	"\n" +
	"expression\x18\x01 \x01(\tR\n" +
	"expression\x12\x14\n" +
	"\x05group\x18\x02 \x01(\x03R\x05group\x124\n" +
	"\bassignee\x18\x03 \x01(\v2\x18.webitel.contacts.LookupR\bassignee\x126\n" +
	"\bposition\x18\x04 \x01(\v2\x1a.webitel.contacts.PositionR\bposition\"@\n" +
	"\x16LocateConditionRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x16\n" +
	"\x06fields\x18\x02 \x03(\tR\x06fields\"T\n" +
	"\x17LocateConditionResponse\x129\n" +
	"\tcondition\x18\x01 \x01(\v2\x1b.webitel.contacts.ConditionR\tcondition\"\xa7\x01\n" +
	"\x16UpdateConditionRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x126\n" +
	"\x05input\x18\x02 \x01(\v2 .webitel.contacts.InputConditionR\x05input\x129\n" +
	"\vx_json_mask\x18\x03 \x03(\tB\x19\x92A\a@\x01\x8a\x01\x02^$\xfa\xd2\xe4\x93\x02\t\x12\aPREVIEWR\txJsonMask:\n" +
	"\x92A\a\n" +
	"\x05\xd2\x01\x02id\"4\n" +
	"\x16DeleteConditionRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id:\n" +
	"\x92A\a\n" +
	"\x05\xd2\x01\x02id2\x8f\a\n" +
	"\x11DynamicConditions\x12\xac\x01\n" +
	"\x0eListConditions\x12'.webitel.contacts.ListConditionsRequest\x1a\x1f.webitel.contacts.ConditionList\"P\x92A\x1f\x12\x1dRetrieve a list of conditions\x82\xd3\xe4\x93\x02(\x12&/contacts/groups/{group_id}/conditions\x12\xa6\x01\n" +
	"\x0fCreateCondition\x12(.webitel.contacts.CreateConditionRequest\x1a\x1b.webitel.contacts.Condition\"L\x92A\x18\x12\x16Create a new condition\x82\xd3\xe4\x93\x02+:\x01*\"&/contacts/groups/{group_id}/conditions\x12\xd5\x01\n" +
	"\x0fUpdateCondition\x12(.webitel.contacts.UpdateConditionRequest\x1a\x1b.webitel.contacts.Condition\"{\x92A\x1e\x12\x1cUpdate an existing condition\x82\xd3\xe4\x93\x02T:\x05inputZ):\x05input2 /contacts/groups/conditions/{id}\x1a /contacts/groups/conditions/{id}\x12\x99\x01\n" +
	"\x0fDeleteCondition\x12(.webitel.contacts.DeleteConditionRequest\x1a\x1b.webitel.contacts.Condition\"?\x92A\x14\x12\x12Delete a condition\x82\xd3\xe4\x93\x02\"* /contacts/groups/conditions/{id}\x12\xad\x01\n" +
	"\x0fLocateCondition\x12(.webitel.contacts.LocateConditionRequest\x1a).webitel.contacts.LocateConditionResponse\"E\x92A\x1a\x12\x18Locate a condition by ID\x82\xd3\xe4\x93\x02\"\x12 /contacts/groups/conditions/{id}B\xb0\x01\n" +
	"\x14com.webitel.contactsB\x15DynamicConditionProtoP\x01Z webitel.go/api/contacts;contacts\xa2\x02\x03WCX\xaa\x02\x10Webitel.Contacts\xca\x02\x10Webitel\\Contacts\xe2\x02\x1cWebitel\\Contacts\\GPBMetadata\xea\x02\x11Webitel::Contactsb\x06proto3"

var (
	file_contacts_dynamic_condition_proto_rawDescOnce sync.Once
	file_contacts_dynamic_condition_proto_rawDescData []byte
)

func file_contacts_dynamic_condition_proto_rawDescGZIP() []byte {
	file_contacts_dynamic_condition_proto_rawDescOnce.Do(func() {
		file_contacts_dynamic_condition_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contacts_dynamic_condition_proto_rawDesc), len(file_contacts_dynamic_condition_proto_rawDesc)))
	})
	return file_contacts_dynamic_condition_proto_rawDescData
}

var file_contacts_dynamic_condition_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_contacts_dynamic_condition_proto_goTypes = []any{
	(*Condition)(nil),               // 0: webitel.contacts.Condition
	(*ListConditionsRequest)(nil),   // 1: webitel.contacts.ListConditionsRequest
	(*ConditionList)(nil),           // 2: webitel.contacts.ConditionList
	(*CreateConditionRequest)(nil),  // 3: webitel.contacts.CreateConditionRequest
	(*Position)(nil),                // 4: webitel.contacts.Position
	(*InputCondition)(nil),          // 5: webitel.contacts.InputCondition
	(*LocateConditionRequest)(nil),  // 6: webitel.contacts.LocateConditionRequest
	(*LocateConditionResponse)(nil), // 7: webitel.contacts.LocateConditionResponse
	(*UpdateConditionRequest)(nil),  // 8: webitel.contacts.UpdateConditionRequest
	(*DeleteConditionRequest)(nil),  // 9: webitel.contacts.DeleteConditionRequest
	(*Lookup)(nil),                  // 10: webitel.contacts.Lookup
}
var file_contacts_dynamic_condition_proto_depIdxs = []int32{
	10, // 0: webitel.contacts.Condition.group:type_name -> webitel.contacts.Lookup
	10, // 1: webitel.contacts.Condition.assignee:type_name -> webitel.contacts.Lookup
	10, // 2: webitel.contacts.Condition.created_by:type_name -> webitel.contacts.Lookup
	10, // 3: webitel.contacts.Condition.updated_by:type_name -> webitel.contacts.Lookup
	0,  // 4: webitel.contacts.ConditionList.items:type_name -> webitel.contacts.Condition
	10, // 5: webitel.contacts.CreateConditionRequest.group:type_name -> webitel.contacts.Lookup
	10, // 6: webitel.contacts.CreateConditionRequest.assignee:type_name -> webitel.contacts.Lookup
	10, // 7: webitel.contacts.InputCondition.assignee:type_name -> webitel.contacts.Lookup
	4,  // 8: webitel.contacts.InputCondition.position:type_name -> webitel.contacts.Position
	0,  // 9: webitel.contacts.LocateConditionResponse.condition:type_name -> webitel.contacts.Condition
	5,  // 10: webitel.contacts.UpdateConditionRequest.input:type_name -> webitel.contacts.InputCondition
	1,  // 11: webitel.contacts.DynamicConditions.ListConditions:input_type -> webitel.contacts.ListConditionsRequest
	3,  // 12: webitel.contacts.DynamicConditions.CreateCondition:input_type -> webitel.contacts.CreateConditionRequest
	8,  // 13: webitel.contacts.DynamicConditions.UpdateCondition:input_type -> webitel.contacts.UpdateConditionRequest
	9,  // 14: webitel.contacts.DynamicConditions.DeleteCondition:input_type -> webitel.contacts.DeleteConditionRequest
	6,  // 15: webitel.contacts.DynamicConditions.LocateCondition:input_type -> webitel.contacts.LocateConditionRequest
	2,  // 16: webitel.contacts.DynamicConditions.ListConditions:output_type -> webitel.contacts.ConditionList
	0,  // 17: webitel.contacts.DynamicConditions.CreateCondition:output_type -> webitel.contacts.Condition
	0,  // 18: webitel.contacts.DynamicConditions.UpdateCondition:output_type -> webitel.contacts.Condition
	0,  // 19: webitel.contacts.DynamicConditions.DeleteCondition:output_type -> webitel.contacts.Condition
	7,  // 20: webitel.contacts.DynamicConditions.LocateCondition:output_type -> webitel.contacts.LocateConditionResponse
	16, // [16:21] is the sub-list for method output_type
	11, // [11:16] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_contacts_dynamic_condition_proto_init() }
func file_contacts_dynamic_condition_proto_init() {
	if File_contacts_dynamic_condition_proto != nil {
		return
	}
	file_contacts_fields_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contacts_dynamic_condition_proto_rawDesc), len(file_contacts_dynamic_condition_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contacts_dynamic_condition_proto_goTypes,
		DependencyIndexes: file_contacts_dynamic_condition_proto_depIdxs,
		MessageInfos:      file_contacts_dynamic_condition_proto_msgTypes,
	}.Build()
	File_contacts_dynamic_condition_proto = out.File
	file_contacts_dynamic_condition_proto_goTypes = nil
	file_contacts_dynamic_condition_proto_depIdxs = nil
}
