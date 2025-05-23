// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: oauth2service.proto

package api

import (
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

type Claim struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Claim) Reset() {
	*x = Claim{}
	mi := &file_oauth2service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Claim) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Claim) ProtoMessage() {}

func (x *Claim) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Claim.ProtoReflect.Descriptor instead.
func (*Claim) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{0}
}

func (x *Claim) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Claim) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// OAuth 2.0 Authentication [S]ervice [P]rovider Application Configuration
type OAuthService struct {
	state  protoimpl.MessageState `protogen:"open.v1"`
	Id     int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // display
	Domain *ObjectId              `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Type   string                 `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"` // well-known vendor; provider
	Logo   string                 `protobuf:"bytes,5,opt,name=logo,proto3" json:"logo,omitempty"`
	// Scopes to be requested
	Scopes []string `protobuf:"bytes,6,rep,name=scopes,proto3" json:"scopes,omitempty"`
	// Identity claims policy rules
	// NOTE: Order matters
	Claims        []*Claim         `protobuf:"bytes,7,rep,name=claims,proto3" json:"claims,omitempty"` // google.protobuf.Struct claims = 7;
	Enabled       bool             `protobuf:"varint,8,opt,name=enabled,proto3" json:"enabled,omitempty"`
	ClientId      string           `protobuf:"bytes,9,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret  string           `protobuf:"bytes,10,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	AuthUrl       string           `protobuf:"bytes,11,opt,name=auth_url,json=authUrl,proto3" json:"auth_url,omitempty"`                // OAuth 2.0 Authorization Endpoint
	TokenUrl      string           `protobuf:"bytes,12,opt,name=token_url,json=tokenUrl,proto3" json:"token_url,omitempty"`             // OAuth 2.0 Token Endpoint
	UserinfoUrl   string           `protobuf:"bytes,13,opt,name=userinfo_url,json=userinfoUrl,proto3" json:"userinfo_url,omitempty"`    // OpenID Connect Userinfo Endpoint
	DiscoveryUrl  string           `protobuf:"bytes,14,opt,name=discovery_url,json=discoveryUrl,proto3" json:"discovery_url,omitempty"` // OpenID Connect Service Discovery
	Metadata      *structpb.Struct `protobuf:"bytes,15,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CreatedAt     int64            `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"` // unix
	CreatedBy     *ObjectId        `protobuf:"bytes,21,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`  // user
	UpdatedAt     int64            `protobuf:"varint,22,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"` // unix
	UpdatedBy     *ObjectId        `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`  // user
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OAuthService) Reset() {
	*x = OAuthService{}
	mi := &file_oauth2service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OAuthService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuthService) ProtoMessage() {}

func (x *OAuthService) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuthService.ProtoReflect.Descriptor instead.
func (*OAuthService) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{1}
}

func (x *OAuthService) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *OAuthService) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OAuthService) GetDomain() *ObjectId {
	if x != nil {
		return x.Domain
	}
	return nil
}

func (x *OAuthService) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *OAuthService) GetLogo() string {
	if x != nil {
		return x.Logo
	}
	return ""
}

func (x *OAuthService) GetScopes() []string {
	if x != nil {
		return x.Scopes
	}
	return nil
}

func (x *OAuthService) GetClaims() []*Claim {
	if x != nil {
		return x.Claims
	}
	return nil
}

func (x *OAuthService) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *OAuthService) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *OAuthService) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *OAuthService) GetAuthUrl() string {
	if x != nil {
		return x.AuthUrl
	}
	return ""
}

func (x *OAuthService) GetTokenUrl() string {
	if x != nil {
		return x.TokenUrl
	}
	return ""
}

func (x *OAuthService) GetUserinfoUrl() string {
	if x != nil {
		return x.UserinfoUrl
	}
	return ""
}

func (x *OAuthService) GetDiscoveryUrl() string {
	if x != nil {
		return x.DiscoveryUrl
	}
	return ""
}

func (x *OAuthService) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *OAuthService) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *OAuthService) GetCreatedBy() *ObjectId {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *OAuthService) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *OAuthService) GetUpdatedBy() *ObjectId {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// SearchOAuthServiceRequest Options
type SearchOAuthServiceRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ----- Select Options -------------------------
	Page   int32    `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`    // default: 1
	Size   int32    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`    // default: 16
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"` // attributes list
	Sort   []string `protobuf:"bytes,4,rep,name=sort,proto3" json:"sort,omitempty"`     // e.g.: "updated_at" - ASC; "!updated_at" - DESC;
	// ----- Search Basic Filters ---------------------------
	Id            []int64 `protobuf:"varint,5,rep,packed,name=id,proto3" json:"id,omitempty"`    // selection: by unique identifier
	Q             string  `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`              // term-of-search: lookup[name]
	Name          string  `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`        // case-ignore substring match: ILIKE '*' - any; '?' - one
	Access        string  `protobuf:"bytes,8,opt,name=access,proto3" json:"access,omitempty"`    // [M]andatory[A]ccess[C]ontrol: with access mode (action) granted!
	Enabled       bool    `protobuf:"varint,9,opt,name=enabled,proto3" json:"enabled,omitempty"` // ----- OAuthService-Specific Filters ----------------
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchOAuthServiceRequest) Reset() {
	*x = SearchOAuthServiceRequest{}
	mi := &file_oauth2service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchOAuthServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchOAuthServiceRequest) ProtoMessage() {}

func (x *SearchOAuthServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchOAuthServiceRequest.ProtoReflect.Descriptor instead.
func (*SearchOAuthServiceRequest) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{2}
}

func (x *SearchOAuthServiceRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchOAuthServiceRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchOAuthServiceRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchOAuthServiceRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *SearchOAuthServiceRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SearchOAuthServiceRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchOAuthServiceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SearchOAuthServiceRequest) GetAccess() string {
	if x != nil {
		return x.Access
	}
	return ""
}

func (x *SearchOAuthServiceRequest) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

type SearchOAuthServiceResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"` // result: {page} number
	Next          bool                   `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"` // result: has {next} page ?
	Items         []*OAuthService        `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchOAuthServiceResponse) Reset() {
	*x = SearchOAuthServiceResponse{}
	mi := &file_oauth2service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchOAuthServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchOAuthServiceResponse) ProtoMessage() {}

func (x *SearchOAuthServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchOAuthServiceResponse.ProtoReflect.Descriptor instead.
func (*SearchOAuthServiceResponse) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{3}
}

func (x *SearchOAuthServiceResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchOAuthServiceResponse) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *SearchOAuthServiceResponse) GetItems() []*OAuthService {
	if x != nil {
		return x.Items
	}
	return nil
}

type UpdateOAuthServiceRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Fields for partial update. PATCH
	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	// Configuration changes.
	Changes       *OAuthService `protobuf:"bytes,2,opt,name=changes,proto3" json:"changes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateOAuthServiceRequest) Reset() {
	*x = UpdateOAuthServiceRequest{}
	mi := &file_oauth2service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateOAuthServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateOAuthServiceRequest) ProtoMessage() {}

func (x *UpdateOAuthServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateOAuthServiceRequest.ProtoReflect.Descriptor instead.
func (*UpdateOAuthServiceRequest) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateOAuthServiceRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *UpdateOAuthServiceRequest) GetChanges() *OAuthService {
	if x != nil {
		return x.Changes
	}
	return nil
}

type DeleteOAuthServiceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []int64                `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	Permanent     bool                   `protobuf:"varint,2,opt,name=permanent,proto3" json:"permanent,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteOAuthServiceRequest) Reset() {
	*x = DeleteOAuthServiceRequest{}
	mi := &file_oauth2service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteOAuthServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteOAuthServiceRequest) ProtoMessage() {}

func (x *DeleteOAuthServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteOAuthServiceRequest.ProtoReflect.Descriptor instead.
func (*DeleteOAuthServiceRequest) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteOAuthServiceRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *DeleteOAuthServiceRequest) GetPermanent() bool {
	if x != nil {
		return x.Permanent
	}
	return false
}

type DeleteOAuthServiceResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteOAuthServiceResponse) Reset() {
	*x = DeleteOAuthServiceResponse{}
	mi := &file_oauth2service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteOAuthServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteOAuthServiceResponse) ProtoMessage() {}

func (x *DeleteOAuthServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_oauth2service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteOAuthServiceResponse.ProtoReflect.Descriptor instead.
func (*DeleteOAuthServiceResponse) Descriptor() ([]byte, []int) {
	return file_oauth2service_proto_rawDescGZIP(), []int{6}
}

var File_oauth2service_proto protoreflect.FileDescriptor

const file_oauth2service_proto_rawDesc = "" +
	"\n" +
	"\x13oauth2service.proto\x12\x03api\x1a\toid.proto\x1a\x1cgoogle/protobuf/struct.proto\x1a\x1cgoogle/api/annotations.proto\"1\n" +
	"\x05Claim\x12\x12\n" +
	"\x04type\x18\x01 \x01(\tR\x04type\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value\"\xe8\x04\n" +
	"\fOAuthService\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12%\n" +
	"\x06domain\x18\x03 \x01(\v2\r.api.ObjectIdR\x06domain\x12\x12\n" +
	"\x04type\x18\x04 \x01(\tR\x04type\x12\x12\n" +
	"\x04logo\x18\x05 \x01(\tR\x04logo\x12\x16\n" +
	"\x06scopes\x18\x06 \x03(\tR\x06scopes\x12\"\n" +
	"\x06claims\x18\a \x03(\v2\n" +
	".api.ClaimR\x06claims\x12\x18\n" +
	"\aenabled\x18\b \x01(\bR\aenabled\x12\x1b\n" +
	"\tclient_id\x18\t \x01(\tR\bclientId\x12#\n" +
	"\rclient_secret\x18\n" +
	" \x01(\tR\fclientSecret\x12\x19\n" +
	"\bauth_url\x18\v \x01(\tR\aauthUrl\x12\x1b\n" +
	"\ttoken_url\x18\f \x01(\tR\btokenUrl\x12!\n" +
	"\fuserinfo_url\x18\r \x01(\tR\vuserinfoUrl\x12#\n" +
	"\rdiscovery_url\x18\x0e \x01(\tR\fdiscoveryUrl\x123\n" +
	"\bmetadata\x18\x0f \x01(\v2\x17.google.protobuf.StructR\bmetadata\x12\x1d\n" +
	"\n" +
	"created_at\x18\x14 \x01(\x03R\tcreatedAt\x12,\n" +
	"\n" +
	"created_by\x18\x15 \x01(\v2\r.api.ObjectIdR\tcreatedBy\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x16 \x01(\x03R\tupdatedAt\x12,\n" +
	"\n" +
	"updated_by\x18\x17 \x01(\v2\r.api.ObjectIdR\tupdatedBy\"\xd3\x01\n" +
	"\x19SearchOAuthServiceRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\x12\x12\n" +
	"\x04sort\x18\x04 \x03(\tR\x04sort\x12\x0e\n" +
	"\x02id\x18\x05 \x03(\x03R\x02id\x12\f\n" +
	"\x01q\x18\x06 \x01(\tR\x01q\x12\x12\n" +
	"\x04name\x18\a \x01(\tR\x04name\x12\x16\n" +
	"\x06access\x18\b \x01(\tR\x06access\x12\x18\n" +
	"\aenabled\x18\t \x01(\bR\aenabled\"m\n" +
	"\x1aSearchOAuthServiceResponse\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x02 \x01(\bR\x04next\x12'\n" +
	"\x05items\x18\x03 \x03(\v2\x11.api.OAuthServiceR\x05items\"`\n" +
	"\x19UpdateOAuthServiceRequest\x12\x16\n" +
	"\x06fields\x18\x01 \x03(\tR\x06fields\x12+\n" +
	"\achanges\x18\x02 \x01(\v2\x11.api.OAuthServiceR\achanges\"I\n" +
	"\x19DeleteOAuthServiceRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x03(\x03R\x02id\x12\x1c\n" +
	"\tpermanent\x18\x02 \x01(\bR\tpermanent\"\x1c\n" +
	"\x1aDeleteOAuthServiceResponse2\xca\x04\n" +
	"\x10OAuth2Federation\x12j\n" +
	"\x12SearchOAuthService\x12\x1e.api.SearchOAuthServiceRequest\x1a\x1f.api.SearchOAuthServiceResponse\"\x13\x82\xd3\xe4\x93\x02\r\x12\v/oauth/apps\x12R\n" +
	"\x12CreateOAuthService\x12\x11.api.OAuthService\x1a\x11.api.OAuthService\"\x16\x82\xd3\xe4\x93\x02\x10:\x01*\"\v/oauth/apps\x12a\n" +
	"\x12LocateOAuthService\x12\x1e.api.SearchOAuthServiceRequest\x1a\x11.api.OAuthService\"\x18\x82\xd3\xe4\x93\x02\x12\x12\x10/oauth/apps/{id}\x12\x8b\x01\n" +
	"\x12UpdateOAuthService\x12\x1e.api.UpdateOAuthServiceRequest\x1a\x11.api.OAuthService\"B\x82\xd3\xe4\x93\x02<:\x01*Z\x1d:\x01*2\x18/oauth/apps/{changes.id}\x1a\x18/oauth/apps/{changes.id}\x12\x84\x01\n" +
	"\x12DeleteOAuthService\x12\x1e.api.DeleteOAuthServiceRequest\x1a\x1f.api.DeleteOAuthServiceResponse\"-\x82\xd3\xe4\x93\x02':\x01*Z\x15:\x01**\x10/oauth/apps/{id}*\v/oauth/appsB]\n" +
	"\acom.apiB\x12Oauth2serviceProtoP\x01Z\x12webitel.go/api;api\xa2\x02\x03AXX\xaa\x02\x03Api\xca\x02\x03Api\xe2\x02\x0fApi\\GPBMetadata\xea\x02\x03Apib\x06proto3"

var (
	file_oauth2service_proto_rawDescOnce sync.Once
	file_oauth2service_proto_rawDescData []byte
)

func file_oauth2service_proto_rawDescGZIP() []byte {
	file_oauth2service_proto_rawDescOnce.Do(func() {
		file_oauth2service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_oauth2service_proto_rawDesc), len(file_oauth2service_proto_rawDesc)))
	})
	return file_oauth2service_proto_rawDescData
}

var file_oauth2service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_oauth2service_proto_goTypes = []any{
	(*Claim)(nil),                      // 0: api.Claim
	(*OAuthService)(nil),               // 1: api.OAuthService
	(*SearchOAuthServiceRequest)(nil),  // 2: api.SearchOAuthServiceRequest
	(*SearchOAuthServiceResponse)(nil), // 3: api.SearchOAuthServiceResponse
	(*UpdateOAuthServiceRequest)(nil),  // 4: api.UpdateOAuthServiceRequest
	(*DeleteOAuthServiceRequest)(nil),  // 5: api.DeleteOAuthServiceRequest
	(*DeleteOAuthServiceResponse)(nil), // 6: api.DeleteOAuthServiceResponse
	(*ObjectId)(nil),                   // 7: api.ObjectId
	(*structpb.Struct)(nil),            // 8: google.protobuf.Struct
}
var file_oauth2service_proto_depIdxs = []int32{
	7,  // 0: api.OAuthService.domain:type_name -> api.ObjectId
	0,  // 1: api.OAuthService.claims:type_name -> api.Claim
	8,  // 2: api.OAuthService.metadata:type_name -> google.protobuf.Struct
	7,  // 3: api.OAuthService.created_by:type_name -> api.ObjectId
	7,  // 4: api.OAuthService.updated_by:type_name -> api.ObjectId
	1,  // 5: api.SearchOAuthServiceResponse.items:type_name -> api.OAuthService
	1,  // 6: api.UpdateOAuthServiceRequest.changes:type_name -> api.OAuthService
	2,  // 7: api.OAuth2Federation.SearchOAuthService:input_type -> api.SearchOAuthServiceRequest
	1,  // 8: api.OAuth2Federation.CreateOAuthService:input_type -> api.OAuthService
	2,  // 9: api.OAuth2Federation.LocateOAuthService:input_type -> api.SearchOAuthServiceRequest
	4,  // 10: api.OAuth2Federation.UpdateOAuthService:input_type -> api.UpdateOAuthServiceRequest
	5,  // 11: api.OAuth2Federation.DeleteOAuthService:input_type -> api.DeleteOAuthServiceRequest
	3,  // 12: api.OAuth2Federation.SearchOAuthService:output_type -> api.SearchOAuthServiceResponse
	1,  // 13: api.OAuth2Federation.CreateOAuthService:output_type -> api.OAuthService
	1,  // 14: api.OAuth2Federation.LocateOAuthService:output_type -> api.OAuthService
	1,  // 15: api.OAuth2Federation.UpdateOAuthService:output_type -> api.OAuthService
	6,  // 16: api.OAuth2Federation.DeleteOAuthService:output_type -> api.DeleteOAuthServiceResponse
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_oauth2service_proto_init() }
func file_oauth2service_proto_init() {
	if File_oauth2service_proto != nil {
		return
	}
	file_oid_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_oauth2service_proto_rawDesc), len(file_oauth2service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_oauth2service_proto_goTypes,
		DependencyIndexes: file_oauth2service_proto_depIdxs,
		MessageInfos:      file_oauth2service_proto_msgTypes,
	}.Build()
	File_oauth2service_proto = out.File
	file_oauth2service_proto_goTypes = nil
	file_oauth2service_proto_depIdxs = nil
}
