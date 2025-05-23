// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: contacts/media.proto

package contacts

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

type MediaAttribute struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Attribute:
	//
	//	*MediaAttribute_Image_
	//	*MediaAttribute_Audio_
	//	*MediaAttribute_Video_
	//	*MediaAttribute_Filename_
	Attribute     isMediaAttribute_Attribute `protobuf_oneof:"attribute"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaAttribute) Reset() {
	*x = MediaAttribute{}
	mi := &file_contacts_media_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaAttribute) ProtoMessage() {}

func (x *MediaAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaAttribute.ProtoReflect.Descriptor instead.
func (*MediaAttribute) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{0}
}

func (x *MediaAttribute) GetAttribute() isMediaAttribute_Attribute {
	if x != nil {
		return x.Attribute
	}
	return nil
}

func (x *MediaAttribute) GetImage() *MediaAttribute_Image {
	if x != nil {
		if x, ok := x.Attribute.(*MediaAttribute_Image_); ok {
			return x.Image
		}
	}
	return nil
}

func (x *MediaAttribute) GetAudio() *MediaAttribute_Audio {
	if x != nil {
		if x, ok := x.Attribute.(*MediaAttribute_Audio_); ok {
			return x.Audio
		}
	}
	return nil
}

func (x *MediaAttribute) GetVideo() *MediaAttribute_Video {
	if x != nil {
		if x, ok := x.Attribute.(*MediaAttribute_Video_); ok {
			return x.Video
		}
	}
	return nil
}

func (x *MediaAttribute) GetFilename() string {
	if x != nil {
		if x, ok := x.Attribute.(*MediaAttribute_Filename_); ok {
			return x.Filename
		}
	}
	return ""
}

type isMediaAttribute_Attribute interface {
	isMediaAttribute_Attribute()
}

type MediaAttribute_Image_ struct {
	Image *MediaAttribute_Image `protobuf:"bytes,1,opt,name=image,proto3,oneof"`
}

type MediaAttribute_Audio_ struct {
	Audio *MediaAttribute_Audio `protobuf:"bytes,2,opt,name=audio,proto3,oneof"`
}

type MediaAttribute_Video_ struct {
	Video *MediaAttribute_Video `protobuf:"bytes,3,opt,name=video,proto3,oneof"`
}

type MediaAttribute_Filename_ struct {
	Filename string `protobuf:"bytes,4,opt,name=filename,proto3,oneof"`
}

func (*MediaAttribute_Image_) isMediaAttribute_Attribute() {}

func (*MediaAttribute_Audio_) isMediaAttribute_Attribute() {}

func (*MediaAttribute_Video_) isMediaAttribute_Attribute() {}

func (*MediaAttribute_Filename_) isMediaAttribute_Attribute() {}

type ImageSize struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Image width
	Width int32 `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	// Image height
	Height int32 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	// Size in bytes
	Size          int32 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ImageSize) Reset() {
	*x = ImageSize{}
	mi := &file_contacts_media_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ImageSize) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageSize) ProtoMessage() {}

func (x *ImageSize) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageSize.ProtoReflect.Descriptor instead.
func (*ImageSize) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{1}
}

func (x *ImageSize) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *ImageSize) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *ImageSize) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

// Animated profile picture in MPEG4 format
type VideoSize struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Video width
	Width uint32 `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	// Video height
	Height uint32 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	// File size in bytes
	Size uint32 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	// Timestamp that should be shown as static preview to the user (seconds)
	StartTs       uint32 `protobuf:"varint,4,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VideoSize) Reset() {
	*x = VideoSize{}
	mi := &file_contacts_media_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoSize) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoSize) ProtoMessage() {}

func (x *VideoSize) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoSize.ProtoReflect.Descriptor instead.
func (*VideoSize) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{2}
}

func (x *VideoSize) GetWidth() uint32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *VideoSize) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *VideoSize) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *VideoSize) GetStartTs() uint32 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

type MediaFile struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// File unique ID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Creation date; timestamp(milli)
	Date int64 `protobuf:"varint,2,opt,name=date,proto3" json:"date,omitempty"`
	// Size in bytes
	Size uint32 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	// MIME type
	Type string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	// Binary content
	Data []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	// Check sum, dependent on unique ID
	Hash map[string]string `protobuf:"bytes,6,rep,name=hash,proto3" json:"hash,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // md5, sha256, ..
	// File attributes metadata
	Meta []*MediaAttribute `protobuf:"bytes,7,rep,name=meta,proto3" json:"meta,omitempty"`
	// Thumbnails
	Thumb []*ImageSize `protobuf:"bytes,8,rep,name=thumb,proto3" json:"thumb,omitempty"`
	// Video Thumbnails
	VideoThumb    []*VideoSize `protobuf:"bytes,9,rep,name=video_thumb,json=videoThumb,proto3" json:"video_thumb,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaFile) Reset() {
	*x = MediaFile{}
	mi := &file_contacts_media_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaFile) ProtoMessage() {}

func (x *MediaFile) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaFile.ProtoReflect.Descriptor instead.
func (*MediaFile) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{3}
}

func (x *MediaFile) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaFile) GetDate() int64 {
	if x != nil {
		return x.Date
	}
	return 0
}

func (x *MediaFile) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *MediaFile) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MediaFile) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *MediaFile) GetHash() map[string]string {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *MediaFile) GetMeta() []*MediaAttribute {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *MediaFile) GetThumb() []*ImageSize {
	if x != nil {
		return x.Thumb
	}
	return nil
}

func (x *MediaFile) GetVideoThumb() []*VideoSize {
	if x != nil {
		return x.VideoThumb
	}
	return nil
}

// An Image or Photo
type MediaImage struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// File unique ID.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Date of upload; timestamp(milli)
	Date int64 `protobuf:"varint,2,opt,name=date,proto3" json:"date,omitempty"`
	// MIME type
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	// Available sizes for download
	Sizes         []*ImageSize `protobuf:"bytes,4,rep,name=sizes,proto3" json:"sizes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaImage) Reset() {
	*x = MediaImage{}
	mi := &file_contacts_media_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaImage) ProtoMessage() {}

func (x *MediaImage) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaImage.ProtoReflect.Descriptor instead.
func (*MediaImage) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{4}
}

func (x *MediaImage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaImage) GetDate() int64 {
	if x != nil {
		return x.Date
	}
	return 0
}

func (x *MediaImage) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MediaImage) GetSizes() []*ImageSize {
	if x != nil {
		return x.Sizes
	}
	return nil
}

type InputMedia struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the uploaded file.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Remote resource URL.
	Url           string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputMedia) Reset() {
	*x = InputMedia{}
	mi := &file_contacts_media_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputMedia) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputMedia) ProtoMessage() {}

func (x *InputMedia) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputMedia.ProtoReflect.Descriptor instead.
func (*InputMedia) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{5}
}

func (x *InputMedia) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *InputMedia) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// Defines the width and height of an image uploaded
type MediaAttribute_Image struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Width of image
	Width uint32 `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	// Height of image
	Height uint32 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	// Defines an animated GIF
	Animated      bool `protobuf:"varint,3,opt,name=animated,proto3" json:"animated,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaAttribute_Image) Reset() {
	*x = MediaAttribute_Image{}
	mi := &file_contacts_media_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaAttribute_Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaAttribute_Image) ProtoMessage() {}

func (x *MediaAttribute_Image) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaAttribute_Image.ProtoReflect.Descriptor instead.
func (*MediaAttribute_Image) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MediaAttribute_Image) GetWidth() uint32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *MediaAttribute_Image) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *MediaAttribute_Image) GetAnimated() bool {
	if x != nil {
		return x.Animated
	}
	return false
}

// // Defines an animated GIF
// message Animated {}
// Defines an audio
type MediaAttribute_Audio struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Name of the song
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Song Performer
	Performer string `protobuf:"bytes,2,opt,name=performer,proto3" json:"performer,omitempty"`
	// Duration in seconds
	Duration uint32 `protobuf:"varint,3,opt,name=duration,proto3" json:"duration,omitempty"`
	// Waveform: consists in a series of bitpacked 5-bit values.
	Waveform      []byte `protobuf:"bytes,4,opt,name=waveform,proto3" json:"waveform,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaAttribute_Audio) Reset() {
	*x = MediaAttribute_Audio{}
	mi := &file_contacts_media_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaAttribute_Audio) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaAttribute_Audio) ProtoMessage() {}

func (x *MediaAttribute_Audio) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaAttribute_Audio.ProtoReflect.Descriptor instead.
func (*MediaAttribute_Audio) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{0, 1}
}

func (x *MediaAttribute_Audio) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *MediaAttribute_Audio) GetPerformer() string {
	if x != nil {
		return x.Performer
	}
	return ""
}

func (x *MediaAttribute_Audio) GetDuration() uint32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *MediaAttribute_Audio) GetWaveform() []byte {
	if x != nil {
		return x.Waveform
	}
	return nil
}

// Defines a video
type MediaAttribute_Video struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Video width
	Width uint32 `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	// Video height
	Height uint32 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	// Duration in seconds
	Duration      uint32 `protobuf:"varint,3,opt,name=duration,proto3" json:"duration,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaAttribute_Video) Reset() {
	*x = MediaAttribute_Video{}
	mi := &file_contacts_media_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaAttribute_Video) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaAttribute_Video) ProtoMessage() {}

func (x *MediaAttribute_Video) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaAttribute_Video.ProtoReflect.Descriptor instead.
func (*MediaAttribute_Video) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{0, 2}
}

func (x *MediaAttribute_Video) GetWidth() uint32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *MediaAttribute_Video) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *MediaAttribute_Video) GetDuration() uint32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

// A simple document with a file name
type MediaAttribute_Filename struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The file name
	FileName      string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MediaAttribute_Filename) Reset() {
	*x = MediaAttribute_Filename{}
	mi := &file_contacts_media_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MediaAttribute_Filename) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaAttribute_Filename) ProtoMessage() {}

func (x *MediaAttribute_Filename) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_media_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaAttribute_Filename.ProtoReflect.Descriptor instead.
func (*MediaAttribute_Filename) Descriptor() ([]byte, []int) {
	return file_contacts_media_proto_rawDescGZIP(), []int{0, 3}
}

func (x *MediaAttribute_Filename) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

var File_contacts_media_proto protoreflect.FileDescriptor

const file_contacts_media_proto_rawDesc = "" +
	"\n" +
	"\x14contacts/media.proto\x12\x10webitel.contacts\"\xbf\x04\n" +
	"\x0eMediaAttribute\x12>\n" +
	"\x05image\x18\x01 \x01(\v2&.webitel.contacts.MediaAttribute.ImageH\x00R\x05image\x12>\n" +
	"\x05audio\x18\x02 \x01(\v2&.webitel.contacts.MediaAttribute.AudioH\x00R\x05audio\x12>\n" +
	"\x05video\x18\x03 \x01(\v2&.webitel.contacts.MediaAttribute.VideoH\x00R\x05video\x12\x1c\n" +
	"\bfilename\x18\x04 \x01(\tH\x00R\bfilename\x1aQ\n" +
	"\x05Image\x12\x14\n" +
	"\x05width\x18\x01 \x01(\rR\x05width\x12\x16\n" +
	"\x06height\x18\x02 \x01(\rR\x06height\x12\x1a\n" +
	"\banimated\x18\x03 \x01(\bR\banimated\x1as\n" +
	"\x05Audio\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12\x1c\n" +
	"\tperformer\x18\x02 \x01(\tR\tperformer\x12\x1a\n" +
	"\bduration\x18\x03 \x01(\rR\bduration\x12\x1a\n" +
	"\bwaveform\x18\x04 \x01(\fR\bwaveform\x1aQ\n" +
	"\x05Video\x12\x14\n" +
	"\x05width\x18\x01 \x01(\rR\x05width\x12\x16\n" +
	"\x06height\x18\x02 \x01(\rR\x06height\x12\x1a\n" +
	"\bduration\x18\x03 \x01(\rR\bduration\x1a'\n" +
	"\bFilename\x12\x1b\n" +
	"\tfile_name\x18\x01 \x01(\tR\bfileNameB\v\n" +
	"\tattribute\"M\n" +
	"\tImageSize\x12\x14\n" +
	"\x05width\x18\x01 \x01(\x05R\x05width\x12\x16\n" +
	"\x06height\x18\x02 \x01(\x05R\x06height\x12\x12\n" +
	"\x04size\x18\x03 \x01(\x05R\x04size\"h\n" +
	"\tVideoSize\x12\x14\n" +
	"\x05width\x18\x01 \x01(\rR\x05width\x12\x16\n" +
	"\x06height\x18\x02 \x01(\rR\x06height\x12\x12\n" +
	"\x04size\x18\x03 \x01(\rR\x04size\x12\x19\n" +
	"\bstart_ts\x18\x04 \x01(\rR\astartTs\"\x86\x03\n" +
	"\tMediaFile\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04date\x18\x02 \x01(\x03R\x04date\x12\x12\n" +
	"\x04size\x18\x03 \x01(\rR\x04size\x12\x12\n" +
	"\x04type\x18\x04 \x01(\tR\x04type\x12\x12\n" +
	"\x04data\x18\x05 \x01(\fR\x04data\x129\n" +
	"\x04hash\x18\x06 \x03(\v2%.webitel.contacts.MediaFile.HashEntryR\x04hash\x124\n" +
	"\x04meta\x18\a \x03(\v2 .webitel.contacts.MediaAttributeR\x04meta\x121\n" +
	"\x05thumb\x18\b \x03(\v2\x1b.webitel.contacts.ImageSizeR\x05thumb\x12<\n" +
	"\vvideo_thumb\x18\t \x03(\v2\x1b.webitel.contacts.VideoSizeR\n" +
	"videoThumb\x1a7\n" +
	"\tHashEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"w\n" +
	"\n" +
	"MediaImage\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04date\x18\x02 \x01(\x03R\x04date\x12\x12\n" +
	"\x04type\x18\x03 \x01(\tR\x04type\x121\n" +
	"\x05sizes\x18\x04 \x03(\v2\x1b.webitel.contacts.ImageSizeR\x05sizes\".\n" +
	"\n" +
	"InputMedia\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x10\n" +
	"\x03url\x18\x02 \x01(\tR\x03urlB\xa5\x01\n" +
	"\x14com.webitel.contactsB\n" +
	"MediaProtoP\x01Z webitel.go/api/contacts;contacts\xa2\x02\x03WCX\xaa\x02\x10Webitel.Contacts\xca\x02\x10Webitel\\Contacts\xe2\x02\x1cWebitel\\Contacts\\GPBMetadata\xea\x02\x11Webitel::Contactsb\x06proto3"

var (
	file_contacts_media_proto_rawDescOnce sync.Once
	file_contacts_media_proto_rawDescData []byte
)

func file_contacts_media_proto_rawDescGZIP() []byte {
	file_contacts_media_proto_rawDescOnce.Do(func() {
		file_contacts_media_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contacts_media_proto_rawDesc), len(file_contacts_media_proto_rawDesc)))
	})
	return file_contacts_media_proto_rawDescData
}

var file_contacts_media_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_contacts_media_proto_goTypes = []any{
	(*MediaAttribute)(nil),          // 0: webitel.contacts.MediaAttribute
	(*ImageSize)(nil),               // 1: webitel.contacts.ImageSize
	(*VideoSize)(nil),               // 2: webitel.contacts.VideoSize
	(*MediaFile)(nil),               // 3: webitel.contacts.MediaFile
	(*MediaImage)(nil),              // 4: webitel.contacts.MediaImage
	(*InputMedia)(nil),              // 5: webitel.contacts.InputMedia
	(*MediaAttribute_Image)(nil),    // 6: webitel.contacts.MediaAttribute.Image
	(*MediaAttribute_Audio)(nil),    // 7: webitel.contacts.MediaAttribute.Audio
	(*MediaAttribute_Video)(nil),    // 8: webitel.contacts.MediaAttribute.Video
	(*MediaAttribute_Filename)(nil), // 9: webitel.contacts.MediaAttribute.Filename
	nil,                             // 10: webitel.contacts.MediaFile.HashEntry
}
var file_contacts_media_proto_depIdxs = []int32{
	6,  // 0: webitel.contacts.MediaAttribute.image:type_name -> webitel.contacts.MediaAttribute.Image
	7,  // 1: webitel.contacts.MediaAttribute.audio:type_name -> webitel.contacts.MediaAttribute.Audio
	8,  // 2: webitel.contacts.MediaAttribute.video:type_name -> webitel.contacts.MediaAttribute.Video
	10, // 3: webitel.contacts.MediaFile.hash:type_name -> webitel.contacts.MediaFile.HashEntry
	0,  // 4: webitel.contacts.MediaFile.meta:type_name -> webitel.contacts.MediaAttribute
	1,  // 5: webitel.contacts.MediaFile.thumb:type_name -> webitel.contacts.ImageSize
	2,  // 6: webitel.contacts.MediaFile.video_thumb:type_name -> webitel.contacts.VideoSize
	1,  // 7: webitel.contacts.MediaImage.sizes:type_name -> webitel.contacts.ImageSize
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_contacts_media_proto_init() }
func file_contacts_media_proto_init() {
	if File_contacts_media_proto != nil {
		return
	}
	file_contacts_media_proto_msgTypes[0].OneofWrappers = []any{
		(*MediaAttribute_Image_)(nil),
		(*MediaAttribute_Audio_)(nil),
		(*MediaAttribute_Video_)(nil),
		(*MediaAttribute_Filename_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contacts_media_proto_rawDesc), len(file_contacts_media_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contacts_media_proto_goTypes,
		DependencyIndexes: file_contacts_media_proto_depIdxs,
		MessageInfos:      file_contacts_media_proto_msgTypes,
	}.Build()
	File_contacts_media_proto = out.File
	file_contacts_media_proto_goTypes = nil
	file_contacts_media_proto_depIdxs = nil
}
