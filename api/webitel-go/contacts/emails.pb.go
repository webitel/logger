// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: contacts/emails.proto

package contacts

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/visibility"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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

// The Contact's email address.
type EmailAddress struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the association. Never changes.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Version of the latest update. Numeric sequence.
	Ver int32 `protobuf:"varint,2,opt,name=ver,proto3" json:"ver,omitempty"`
	// Unique ID of the latest version of the update.
	// This ID changes after any update to the underlying value(s).
	Etag string `protobuf:"bytes,3,opt,name=etag,proto3" json:"etag,omitempty"`
	// The user who created this Field.
	CreatedAt int64 `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Timestamp(milli) of the Field creation.
	CreatedBy *Lookup `protobuf:"bytes,6,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// Timestamp(milli) of the last Field update.
	// Take part in Etag generation.
	UpdatedAt int64 `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The user who performed last Update.
	UpdatedBy *Lookup `protobuf:"bytes,8,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// Indicates whether this phone number is default within other channels of the same type(phone).
	Primary bool `protobuf:"varint,11,opt,name=primary,proto3" json:"primary,omitempty"`
	// Indicate whether Contact, as a Person, realy owns this associated phone number.
	// In other words: whether Contact is reachable thru this 'email' communication channel ?
	Verified bool `protobuf:"varint,12,opt,name=verified,proto3" json:"verified,omitempty"`
	// The email address.
	Email string `protobuf:"bytes,13,opt,name=email,proto3" json:"email,omitempty"`
	// The type of the email address.
	// Lookup value from CommunicationType dictionary.
	// The type can be custom or one of these predefined values:
	// - home
	// - work
	// - other
	Type          *Lookup `protobuf:"bytes,14,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailAddress) Reset() {
	*x = EmailAddress{}
	mi := &file_contacts_emails_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailAddress) ProtoMessage() {}

func (x *EmailAddress) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailAddress.ProtoReflect.Descriptor instead.
func (*EmailAddress) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{0}
}

func (x *EmailAddress) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EmailAddress) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *EmailAddress) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *EmailAddress) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *EmailAddress) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *EmailAddress) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *EmailAddress) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *EmailAddress) GetPrimary() bool {
	if x != nil {
		return x.Primary
	}
	return false
}

func (x *EmailAddress) GetVerified() bool {
	if x != nil {
		return x.Verified
	}
	return false
}

func (x *EmailAddress) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *EmailAddress) GetType() *Lookup {
	if x != nil {
		return x.Type
	}
	return nil
}

// Input of the Contact's email address.
type InputEmailAddress struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique ID of the latest version of an existing resorce.
	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	// Indicates whether this phone number is default within other channels of the same type(phone).
	Primary bool `protobuf:"varint,11,opt,name=primary,proto3" json:"primary,omitempty"`
	// Indicate whether Contact, as a Person, realy owns this associated phone number.
	// In other words: whether Contact is reachable thru this 'email' communication channel ?
	Verified bool `protobuf:"varint,12,opt,name=verified,proto3" json:"verified,omitempty"`
	// The email address.
	Email string `protobuf:"bytes,13,opt,name=email,proto3" json:"email,omitempty"`
	// The type of the email address.
	// Lookup value from CommunicationType dictionary.
	// The type can be custom or one of these predefined values:
	// - home
	// - work
	// - other
	Type          *Lookup `protobuf:"bytes,14,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputEmailAddress) Reset() {
	*x = InputEmailAddress{}
	mi := &file_contacts_emails_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputEmailAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputEmailAddress) ProtoMessage() {}

func (x *InputEmailAddress) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputEmailAddress.ProtoReflect.Descriptor instead.
func (*InputEmailAddress) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{1}
}

func (x *InputEmailAddress) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *InputEmailAddress) GetPrimary() bool {
	if x != nil {
		return x.Primary
	}
	return false
}

func (x *InputEmailAddress) GetVerified() bool {
	if x != nil {
		return x.Verified
	}
	return false
}

func (x *InputEmailAddress) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *InputEmailAddress) GetType() *Lookup {
	if x != nil {
		return x.Type
	}
	return nil
}

// Email dataset.
type EmailList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// EmailAddress dataset page.
	Data []*EmailAddress `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	// The page number of the partial result.
	Page int32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	// Indicates that this is a partial result.
	// More data available upon query: ?size=${data.len}&page=${page++}
	Next          bool `protobuf:"varint,3,opt,name=next,proto3" json:"next,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailList) Reset() {
	*x = EmailList{}
	mi := &file_contacts_emails_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailList) ProtoMessage() {}

func (x *EmailList) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailList.ProtoReflect.Descriptor instead.
func (*EmailList) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{2}
}

func (x *EmailList) GetData() []*EmailAddress {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *EmailList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *EmailList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

// Locate single Link by unique IDentifier.
type LocateEmailRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Fields to be retrieved into result.
	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	// Contact source ID.
	ContactId string `protobuf:"bytes,2,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// Unique mail address link IDentifier.
	// Accept: `etag` (obsolete+) or `id`.
	Etag          string `protobuf:"bytes,3,opt,name=etag,proto3" json:"etag,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateEmailRequest) Reset() {
	*x = LocateEmailRequest{}
	mi := &file_contacts_emails_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateEmailRequest) ProtoMessage() {}

func (x *LocateEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateEmailRequest.ProtoReflect.Descriptor instead.
func (*LocateEmailRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{3}
}

func (x *LocateEmailRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *LocateEmailRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *LocateEmailRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

type ListEmailsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of result dataset records. offset = (page*size)
	Page int32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	// Size count of records on result page. limit = (size++)
	Size int32 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	// Search term: email address.
	// `?` - matches any one character
	// `*` - matches 0 or more characters
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Sort the result according to fields.
	Sort []string `protobuf:"bytes,3,rep,name=sort,proto3" json:"sort,omitempty"`
	// Fields to be retrieved into result.
	Fields []string `protobuf:"bytes,4,rep,name=fields,proto3" json:"fields,omitempty"`
	// The Contact ID linked with.
	ContactId string `protobuf:"bytes,5,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// Link(s) with unique ID only.
	Id []string `protobuf:"bytes,7,rep,name=id,proto3" json:"id,omitempty"`
	// Primary email address only.
	Primary *wrapperspb.BoolValue `protobuf:"bytes,8,opt,name=primary,proto3" json:"primary,omitempty"`
	// Verified email addresses only.
	Verified *wrapperspb.BoolValue `protobuf:"bytes,9,opt,name=verified,proto3" json:"verified,omitempty"`
	// Certain communication type associated with the address.
	Type          *Lookup `protobuf:"bytes,10,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListEmailsRequest) Reset() {
	*x = ListEmailsRequest{}
	mi := &file_contacts_emails_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEmailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEmailsRequest) ProtoMessage() {}

func (x *ListEmailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEmailsRequest.ProtoReflect.Descriptor instead.
func (*ListEmailsRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{4}
}

func (x *ListEmailsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListEmailsRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListEmailsRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListEmailsRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListEmailsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListEmailsRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *ListEmailsRequest) GetId() []string {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListEmailsRequest) GetPrimary() *wrapperspb.BoolValue {
	if x != nil {
		return x.Primary
	}
	return nil
}

func (x *ListEmailsRequest) GetVerified() *wrapperspb.BoolValue {
	if x != nil {
		return x.Verified
	}
	return nil
}

func (x *ListEmailsRequest) GetType() *Lookup {
	if x != nil {
		return x.Type
	}
	return nil
}

type MergeEmailsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// JSON PATCH fields mask.
	// List of JPath fields specified in body(input).
	XJsonMask []string `protobuf:"bytes,1,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	// Fields to be retrieved into result of changes.
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	// Link contact ID.
	ContactId string `protobuf:"bytes,3,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// Fixed set of email address(es) to be linked with the contact.
	// Email address(es) that conflicts(email) with already linked will be updated.
	Input         []*InputEmailAddress `protobuf:"bytes,4,rep,name=input,proto3" json:"input,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MergeEmailsRequest) Reset() {
	*x = MergeEmailsRequest{}
	mi := &file_contacts_emails_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MergeEmailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergeEmailsRequest) ProtoMessage() {}

func (x *MergeEmailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergeEmailsRequest.ProtoReflect.Descriptor instead.
func (*MergeEmailsRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{5}
}

func (x *MergeEmailsRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

func (x *MergeEmailsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *MergeEmailsRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *MergeEmailsRequest) GetInput() []*InputEmailAddress {
	if x != nil {
		return x.Input
	}
	return nil
}

type ResetEmailsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Fields to be retrieved into result of changes.
	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	// Link contact ID.
	ContactId string `protobuf:"bytes,2,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// Final set of email address(es) to be linked with the contact.
	// Email address(es) that are already linked with the contact
	// but not given in here will be removed.
	Input         []*InputEmailAddress `protobuf:"bytes,3,rep,name=input,proto3" json:"input,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResetEmailsRequest) Reset() {
	*x = ResetEmailsRequest{}
	mi := &file_contacts_emails_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetEmailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetEmailsRequest) ProtoMessage() {}

func (x *ResetEmailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetEmailsRequest.ProtoReflect.Descriptor instead.
func (*ResetEmailsRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{6}
}

func (x *ResetEmailsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ResetEmailsRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *ResetEmailsRequest) GetInput() []*InputEmailAddress {
	if x != nil {
		return x.Input
	}
	return nil
}

type UpdateEmailRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// JSON PATCH fields mask.
	// List of JPath fields specified in body(input).
	XJsonMask []string `protobuf:"bytes,1,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	// Fields to be retrieved into result of changes.
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	// Link contact ID.
	ContactId string `protobuf:"bytes,3,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// NEW Update of the email address link.
	Input         *InputEmailAddress `protobuf:"bytes,4,opt,name=input,proto3" json:"input,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateEmailRequest) Reset() {
	*x = UpdateEmailRequest{}
	mi := &file_contacts_emails_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEmailRequest) ProtoMessage() {}

func (x *UpdateEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateEmailRequest.ProtoReflect.Descriptor instead.
func (*UpdateEmailRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateEmailRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

func (x *UpdateEmailRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *UpdateEmailRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *UpdateEmailRequest) GetInput() *InputEmailAddress {
	if x != nil {
		return x.Input
	}
	return nil
}

type DeleteEmailsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	// Contact ID associated with.
	ContactId string `protobuf:"bytes,2,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// Set of unique ID(s) to remove.
	Etag          []string `protobuf:"bytes,3,rep,name=etag,proto3" json:"etag,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteEmailsRequest) Reset() {
	*x = DeleteEmailsRequest{}
	mi := &file_contacts_emails_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteEmailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEmailsRequest) ProtoMessage() {}

func (x *DeleteEmailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEmailsRequest.ProtoReflect.Descriptor instead.
func (*DeleteEmailsRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteEmailsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *DeleteEmailsRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *DeleteEmailsRequest) GetEtag() []string {
	if x != nil {
		return x.Etag
	}
	return nil
}

type DeleteEmailRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	// Contact ID associated with.
	ContactId string `protobuf:"bytes,2,opt,name=contact_id,json=contactId,proto3" json:"contact_id,omitempty"`
	// Unique ID to remove.
	Etag          string `protobuf:"bytes,3,opt,name=etag,proto3" json:"etag,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteEmailRequest) Reset() {
	*x = DeleteEmailRequest{}
	mi := &file_contacts_emails_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEmailRequest) ProtoMessage() {}

func (x *DeleteEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_emails_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEmailRequest.ProtoReflect.Descriptor instead.
func (*DeleteEmailRequest) Descriptor() ([]byte, []int) {
	return file_contacts_emails_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteEmailRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *DeleteEmailRequest) GetContactId() string {
	if x != nil {
		return x.ContactId
	}
	return ""
}

func (x *DeleteEmailRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

var File_contacts_emails_proto protoreflect.FileDescriptor

const file_contacts_emails_proto_rawDesc = "" +
	"\n" +
	"\x15contacts/emails.proto\x12\x10webitel.contacts\x1a\x15contacts/fields.proto\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x1bgoogle/api/visibility.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\xf3\x02\n" +
	"\fEmailAddress\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x10\n" +
	"\x03ver\x18\x02 \x01(\x05R\x03ver\x12\x12\n" +
	"\x04etag\x18\x03 \x01(\tR\x04etag\x12\x1d\n" +
	"\n" +
	"created_at\x18\x05 \x01(\x03R\tcreatedAt\x127\n" +
	"\n" +
	"created_by\x18\x06 \x01(\v2\x18.webitel.contacts.LookupR\tcreatedBy\x12\x1d\n" +
	"\n" +
	"updated_at\x18\a \x01(\x03R\tupdatedAt\x127\n" +
	"\n" +
	"updated_by\x18\b \x01(\v2\x18.webitel.contacts.LookupR\tupdatedBy\x12\x18\n" +
	"\aprimary\x18\v \x01(\bR\aprimary\x12\x1a\n" +
	"\bverified\x18\f \x01(\bR\bverified\x12\x14\n" +
	"\x05email\x18\r \x01(\tR\x05email\x12,\n" +
	"\x04type\x18\x0e \x01(\v2\x18.webitel.contacts.LookupR\x04type:\x03\x92A\x00\"\xaa\x02\n" +
	"\x11InputEmailAddress\x12!\n" +
	"\x04etag\x18\x01 \x01(\tB\r\x92A\n" +
	"\xca>\a\xfa\x02\x04etagR\x04etag\x12\x18\n" +
	"\aprimary\x18\v \x01(\bR\aprimary\x12\x1a\n" +
	"\bverified\x18\f \x01(\bR\bverified\x12\x14\n" +
	"\x05email\x18\r \x01(\tR\x05email\x12,\n" +
	"\x04type\x18\x0e \x01(\v2\x18.webitel.contacts.LookupR\x04type:x\x92Au\n" +
	"\b\xd2\x01\x05email2i{\"etag\":\"1679792219687\",\"verified\":false,\"primary\":true,\"email\":\"user@domain\",\"type\":{\"name\":\"personal\"}}\"l\n" +
	"\tEmailList\x122\n" +
	"\x04data\x18\x01 \x03(\v2\x1e.webitel.contacts.EmailAddressR\x04data\x12\x12\n" +
	"\x04page\x18\x02 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x03 \x01(\bR\x04next:\x03\x92A\x00\"j\n" +
	"\x12LocateEmailRequest\x12\x16\n" +
	"\x06fields\x18\x01 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x02 \x01(\tR\tcontactId\x12\x1d\n" +
	"\x04etag\x18\x03 \x01(\tB\t\x92A\x06\xa2\x02\x03\\w+R\x04etag\"\xc0\x02\n" +
	"\x11ListEmailsRequest\x12\x12\n" +
	"\x04page\x18\x02 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x01 \x01(\x05R\x04size\x12\f\n" +
	"\x01q\x18\x06 \x01(\tR\x01q\x12\x12\n" +
	"\x04sort\x18\x03 \x03(\tR\x04sort\x12\x16\n" +
	"\x06fields\x18\x04 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x05 \x01(\tR\tcontactId\x12\x0e\n" +
	"\x02id\x18\a \x03(\tR\x02id\x124\n" +
	"\aprimary\x18\b \x01(\v2\x1a.google.protobuf.BoolValueR\aprimary\x126\n" +
	"\bverified\x18\t \x01(\v2\x1a.google.protobuf.BoolValueR\bverified\x12,\n" +
	"\x04type\x18\n" +
	" \x01(\v2\x18.webitel.contacts.LookupR\x04type\"\x9d\x02\n" +
	"\x12MergeEmailsRequest\x129\n" +
	"\vx_json_mask\x18\x01 \x03(\tB\x19\x92A\a@\x01\x8a\x01\x02^$\xfa\xd2\xe4\x93\x02\t\x12\aPREVIEWR\txJsonMask\x12\x16\n" +
	"\x06fields\x18\x02 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x03 \x01(\tR\tcontactId\x12\x94\x01\n" +
	"\x05input\x18\x04 \x03(\v2#.webitel.contacts.InputEmailAddressBY\x92AVJT[{\"verified\":false,\"primary\":true,\"email\":\"user@domain\",\"type\":{\"name\":\"personal\"}}]R\x05input\"\xe1\x02\n" +
	"\x12ResetEmailsRequest\x12\x16\n" +
	"\x06fields\x18\x01 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x02 \x01(\tR\tcontactId\x12\x84\x02\n" +
	"\x05input\x18\x03 \x03(\v2#.webitel.contacts.InputEmailAddressB\xc8\x01\x92A\xc4\x01J\xbb\x01[{\"verified\":true,\"email\":\"johndoe_43@gmail.com\",\"type\":{\"name\":\"personal\"}},{\"primary\":true,\"etag\":\"k0WqvUn4IJGnuCyG\",\"email\":\"j.doe@x-company.org\",\"type\":{\"id\":\"11\",\"name\":\"business\"}}]\xa8\x01\x01\xb0\x01\x01R\x05input:\r\x92A\n" +
	"\n" +
	"\b\xd2\x01\x05input\"\xd0\x01\n" +
	"\x12UpdateEmailRequest\x129\n" +
	"\vx_json_mask\x18\x01 \x03(\tB\x19\x92A\a@\x01\x8a\x01\x02^$\xfa\xd2\xe4\x93\x02\t\x12\aPREVIEWR\txJsonMask\x12\x16\n" +
	"\x06fields\x18\x02 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x03 \x01(\tR\tcontactId\x129\n" +
	"\x05input\x18\x04 \x01(\v2#.webitel.contacts.InputEmailAddressR\x05input:\r\x92A\n" +
	"\n" +
	"\b\xd2\x01\x05input\"\x8d\x01\n" +
	"\x13DeleteEmailsRequest\x12\x16\n" +
	"\x06fields\x18\x01 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x02 \x01(\tR\tcontactId\x12$\n" +
	"\x04etag\x18\x03 \x03(\tB\x10\x92A\r\x8a\x01\x04^.+$\xa8\x01\x01\xb0\x01\x01R\x04etag:\x19\x92A\x16\n" +
	"\x14\xd2\x01\n" +
	"contact_id\xd2\x01\x04etag\"\x87\x01\n" +
	"\x12DeleteEmailRequest\x12\x16\n" +
	"\x06fields\x18\x01 \x03(\tR\x06fields\x12\x1d\n" +
	"\n" +
	"contact_id\x18\x02 \x01(\tR\tcontactId\x12\x1f\n" +
	"\x04etag\x18\x03 \x01(\tB\v\x92A\b\x8a\x01\x05^\\.+$R\x04etag:\x19\x92A\x16\n" +
	"\x14\xd2\x01\n" +
	"contact_id\xd2\x01\x04etag2\xc5\a\n" +
	"\x06Emails\x12u\n" +
	"\n" +
	"ListEmails\x12#.webitel.contacts.ListEmailsRequest\x1a\x1b.webitel.contacts.EmailList\"%\x82\xd3\xe4\x93\x02\x1f\x12\x1d/contacts/{contact_id}/emails\x12~\n" +
	"\vMergeEmails\x12$.webitel.contacts.MergeEmailsRequest\x1a\x1b.webitel.contacts.EmailList\",\x82\xd3\xe4\x93\x02&:\x05input\"\x1d/contacts/{contact_id}/emails\x12~\n" +
	"\vResetEmails\x12$.webitel.contacts.ResetEmailsRequest\x1a\x1b.webitel.contacts.EmailList\",\x82\xd3\xe4\x93\x02&:\x05input\x1a\x1d/contacts/{contact_id}/emails\x12y\n" +
	"\fDeleteEmails\x12%.webitel.contacts.DeleteEmailsRequest\x1a\x1b.webitel.contacts.EmailList\"%\x82\xd3\xe4\x93\x02\x1f*\x1d/contacts/{contact_id}/emails\x12\x81\x01\n" +
	"\vLocateEmail\x12$.webitel.contacts.LocateEmailRequest\x1a\x1e.webitel.contacts.EmailAddress\",\x82\xd3\xe4\x93\x02&\x12$/contacts/{contact_id}/emails/{etag}\x12\xc0\x01\n" +
	"\vUpdateEmail\x12$.webitel.contacts.UpdateEmailRequest\x1a\x1b.webitel.contacts.EmailList\"n\x82\xd3\xe4\x93\x02h:\x05inputZ3:\x05input2*/contacts/{contact_id}/emails/{input.etag}\x1a*/contacts/{contact_id}/emails/{input.etag}\x12\x81\x01\n" +
	"\vDeleteEmail\x12$.webitel.contacts.DeleteEmailRequest\x1a\x1e.webitel.contacts.EmailAddress\",\x82\xd3\xe4\x93\x02&*$/contacts/{contact_id}/emails/{etag}B\xa6\x01\n" +
	"\x14com.webitel.contactsB\vEmailsProtoP\x01Z webitel.go/api/contacts;contacts\xa2\x02\x03WCX\xaa\x02\x10Webitel.Contacts\xca\x02\x10Webitel\\Contacts\xe2\x02\x1cWebitel\\Contacts\\GPBMetadata\xea\x02\x11Webitel::Contactsb\x06proto3"

var (
	file_contacts_emails_proto_rawDescOnce sync.Once
	file_contacts_emails_proto_rawDescData []byte
)

func file_contacts_emails_proto_rawDescGZIP() []byte {
	file_contacts_emails_proto_rawDescOnce.Do(func() {
		file_contacts_emails_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contacts_emails_proto_rawDesc), len(file_contacts_emails_proto_rawDesc)))
	})
	return file_contacts_emails_proto_rawDescData
}

var file_contacts_emails_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_contacts_emails_proto_goTypes = []any{
	(*EmailAddress)(nil),         // 0: webitel.contacts.EmailAddress
	(*InputEmailAddress)(nil),    // 1: webitel.contacts.InputEmailAddress
	(*EmailList)(nil),            // 2: webitel.contacts.EmailList
	(*LocateEmailRequest)(nil),   // 3: webitel.contacts.LocateEmailRequest
	(*ListEmailsRequest)(nil),    // 4: webitel.contacts.ListEmailsRequest
	(*MergeEmailsRequest)(nil),   // 5: webitel.contacts.MergeEmailsRequest
	(*ResetEmailsRequest)(nil),   // 6: webitel.contacts.ResetEmailsRequest
	(*UpdateEmailRequest)(nil),   // 7: webitel.contacts.UpdateEmailRequest
	(*DeleteEmailsRequest)(nil),  // 8: webitel.contacts.DeleteEmailsRequest
	(*DeleteEmailRequest)(nil),   // 9: webitel.contacts.DeleteEmailRequest
	(*Lookup)(nil),               // 10: webitel.contacts.Lookup
	(*wrapperspb.BoolValue)(nil), // 11: google.protobuf.BoolValue
}
var file_contacts_emails_proto_depIdxs = []int32{
	10, // 0: webitel.contacts.EmailAddress.created_by:type_name -> webitel.contacts.Lookup
	10, // 1: webitel.contacts.EmailAddress.updated_by:type_name -> webitel.contacts.Lookup
	10, // 2: webitel.contacts.EmailAddress.type:type_name -> webitel.contacts.Lookup
	10, // 3: webitel.contacts.InputEmailAddress.type:type_name -> webitel.contacts.Lookup
	0,  // 4: webitel.contacts.EmailList.data:type_name -> webitel.contacts.EmailAddress
	11, // 5: webitel.contacts.ListEmailsRequest.primary:type_name -> google.protobuf.BoolValue
	11, // 6: webitel.contacts.ListEmailsRequest.verified:type_name -> google.protobuf.BoolValue
	10, // 7: webitel.contacts.ListEmailsRequest.type:type_name -> webitel.contacts.Lookup
	1,  // 8: webitel.contacts.MergeEmailsRequest.input:type_name -> webitel.contacts.InputEmailAddress
	1,  // 9: webitel.contacts.ResetEmailsRequest.input:type_name -> webitel.contacts.InputEmailAddress
	1,  // 10: webitel.contacts.UpdateEmailRequest.input:type_name -> webitel.contacts.InputEmailAddress
	4,  // 11: webitel.contacts.Emails.ListEmails:input_type -> webitel.contacts.ListEmailsRequest
	5,  // 12: webitel.contacts.Emails.MergeEmails:input_type -> webitel.contacts.MergeEmailsRequest
	6,  // 13: webitel.contacts.Emails.ResetEmails:input_type -> webitel.contacts.ResetEmailsRequest
	8,  // 14: webitel.contacts.Emails.DeleteEmails:input_type -> webitel.contacts.DeleteEmailsRequest
	3,  // 15: webitel.contacts.Emails.LocateEmail:input_type -> webitel.contacts.LocateEmailRequest
	7,  // 16: webitel.contacts.Emails.UpdateEmail:input_type -> webitel.contacts.UpdateEmailRequest
	9,  // 17: webitel.contacts.Emails.DeleteEmail:input_type -> webitel.contacts.DeleteEmailRequest
	2,  // 18: webitel.contacts.Emails.ListEmails:output_type -> webitel.contacts.EmailList
	2,  // 19: webitel.contacts.Emails.MergeEmails:output_type -> webitel.contacts.EmailList
	2,  // 20: webitel.contacts.Emails.ResetEmails:output_type -> webitel.contacts.EmailList
	2,  // 21: webitel.contacts.Emails.DeleteEmails:output_type -> webitel.contacts.EmailList
	0,  // 22: webitel.contacts.Emails.LocateEmail:output_type -> webitel.contacts.EmailAddress
	2,  // 23: webitel.contacts.Emails.UpdateEmail:output_type -> webitel.contacts.EmailList
	0,  // 24: webitel.contacts.Emails.DeleteEmail:output_type -> webitel.contacts.EmailAddress
	18, // [18:25] is the sub-list for method output_type
	11, // [11:18] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_contacts_emails_proto_init() }
func file_contacts_emails_proto_init() {
	if File_contacts_emails_proto != nil {
		return
	}
	file_contacts_fields_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contacts_emails_proto_rawDesc), len(file_contacts_emails_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contacts_emails_proto_goTypes,
		DependencyIndexes: file_contacts_emails_proto_depIdxs,
		MessageInfos:      file_contacts_emails_proto_msgTypes,
	}.Build()
	File_contacts_emails_proto = out.File
	file_contacts_emails_proto_goTypes = nil
	file_contacts_emails_proto_depIdxs = nil
}
