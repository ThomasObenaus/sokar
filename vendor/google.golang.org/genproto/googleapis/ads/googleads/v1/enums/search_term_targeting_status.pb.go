// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/enums/search_term_targeting_status.proto

package enums

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Indicates whether the search term is one of your targeted or excluded
// keywords.
type SearchTermTargetingStatusEnum_SearchTermTargetingStatus int32

const (
	// Not specified.
	SearchTermTargetingStatusEnum_UNSPECIFIED SearchTermTargetingStatusEnum_SearchTermTargetingStatus = 0
	// Used for return value only. Represents value unknown in this version.
	SearchTermTargetingStatusEnum_UNKNOWN SearchTermTargetingStatusEnum_SearchTermTargetingStatus = 1
	// Search term is added to targeted keywords.
	SearchTermTargetingStatusEnum_ADDED SearchTermTargetingStatusEnum_SearchTermTargetingStatus = 2
	// Search term matches a negative keyword.
	SearchTermTargetingStatusEnum_EXCLUDED SearchTermTargetingStatusEnum_SearchTermTargetingStatus = 3
	// Search term has been both added and excluded.
	SearchTermTargetingStatusEnum_ADDED_EXCLUDED SearchTermTargetingStatusEnum_SearchTermTargetingStatus = 4
	// Search term is neither targeted nor excluded.
	SearchTermTargetingStatusEnum_NONE SearchTermTargetingStatusEnum_SearchTermTargetingStatus = 5
)

var SearchTermTargetingStatusEnum_SearchTermTargetingStatus_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "ADDED",
	3: "EXCLUDED",
	4: "ADDED_EXCLUDED",
	5: "NONE",
}

var SearchTermTargetingStatusEnum_SearchTermTargetingStatus_value = map[string]int32{
	"UNSPECIFIED":    0,
	"UNKNOWN":        1,
	"ADDED":          2,
	"EXCLUDED":       3,
	"ADDED_EXCLUDED": 4,
	"NONE":           5,
}

func (x SearchTermTargetingStatusEnum_SearchTermTargetingStatus) String() string {
	return proto.EnumName(SearchTermTargetingStatusEnum_SearchTermTargetingStatus_name, int32(x))
}

func (SearchTermTargetingStatusEnum_SearchTermTargetingStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e85ee4a345f96b8f, []int{0, 0}
}

// Container for enum indicating whether a search term is one of your targeted
// or excluded keywords.
type SearchTermTargetingStatusEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchTermTargetingStatusEnum) Reset()         { *m = SearchTermTargetingStatusEnum{} }
func (m *SearchTermTargetingStatusEnum) String() string { return proto.CompactTextString(m) }
func (*SearchTermTargetingStatusEnum) ProtoMessage()    {}
func (*SearchTermTargetingStatusEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_e85ee4a345f96b8f, []int{0}
}

func (m *SearchTermTargetingStatusEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchTermTargetingStatusEnum.Unmarshal(m, b)
}
func (m *SearchTermTargetingStatusEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchTermTargetingStatusEnum.Marshal(b, m, deterministic)
}
func (m *SearchTermTargetingStatusEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchTermTargetingStatusEnum.Merge(m, src)
}
func (m *SearchTermTargetingStatusEnum) XXX_Size() int {
	return xxx_messageInfo_SearchTermTargetingStatusEnum.Size(m)
}
func (m *SearchTermTargetingStatusEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchTermTargetingStatusEnum.DiscardUnknown(m)
}

var xxx_messageInfo_SearchTermTargetingStatusEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v1.enums.SearchTermTargetingStatusEnum_SearchTermTargetingStatus", SearchTermTargetingStatusEnum_SearchTermTargetingStatus_name, SearchTermTargetingStatusEnum_SearchTermTargetingStatus_value)
	proto.RegisterType((*SearchTermTargetingStatusEnum)(nil), "google.ads.googleads.v1.enums.SearchTermTargetingStatusEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/enums/search_term_targeting_status.proto", fileDescriptor_e85ee4a345f96b8f)
}

var fileDescriptor_e85ee4a345f96b8f = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x50, 0x4b, 0x4e, 0xc3, 0x30,
	0x10, 0x25, 0x69, 0x0b, 0xc5, 0x45, 0x10, 0x79, 0x07, 0xa2, 0x48, 0xed, 0x01, 0x1c, 0x45, 0xec,
	0xcc, 0x86, 0xb4, 0x09, 0x55, 0x05, 0x4a, 0x2b, 0xf5, 0x03, 0x42, 0x91, 0x22, 0xd3, 0x58, 0x26,
	0x52, 0x63, 0x47, 0xb1, 0xdb, 0x7b, 0x70, 0x05, 0x96, 0x1c, 0x85, 0xa3, 0xb0, 0xe5, 0x02, 0x28,
	0x4e, 0x93, 0x5d, 0xd8, 0x58, 0xcf, 0xf3, 0x66, 0xde, 0x9b, 0x79, 0xe0, 0x9e, 0x09, 0xc1, 0xb6,
	0xd4, 0x26, 0xb1, 0xb4, 0x4b, 0x58, 0xa0, 0xbd, 0x63, 0x53, 0xbe, 0x4b, 0xa5, 0x2d, 0x29, 0xc9,
	0x37, 0xef, 0x91, 0xa2, 0x79, 0x1a, 0x29, 0x92, 0x33, 0xaa, 0x12, 0xce, 0x22, 0xa9, 0x88, 0xda,
	0x49, 0x94, 0xe5, 0x42, 0x09, 0xd8, 0x2f, 0xc7, 0x10, 0x89, 0x25, 0xaa, 0x15, 0xd0, 0xde, 0x41,
	0x5a, 0xe1, 0xea, 0xba, 0x32, 0xc8, 0x12, 0x9b, 0x70, 0x2e, 0x14, 0x51, 0x89, 0xe0, 0x87, 0xe1,
	0xe1, 0x87, 0x01, 0xfa, 0x0b, 0xed, 0xb1, 0xa4, 0x79, 0xba, 0xac, 0x1c, 0x16, 0xda, 0xc0, 0xe7,
	0xbb, 0x74, 0x98, 0x81, 0xcb, 0xc6, 0x06, 0x78, 0x01, 0x7a, 0xab, 0x60, 0x31, 0xf7, 0xc7, 0xd3,
	0x87, 0xa9, 0xef, 0x59, 0x47, 0xb0, 0x07, 0x4e, 0x56, 0xc1, 0x63, 0x30, 0x7b, 0x0e, 0x2c, 0x03,
	0x9e, 0x82, 0x8e, 0xeb, 0x79, 0xbe, 0x67, 0x99, 0xf0, 0x0c, 0x74, 0xfd, 0x97, 0xf1, 0xd3, 0xaa,
	0xf8, 0xb5, 0x20, 0x04, 0xe7, 0x9a, 0x88, 0xea, 0x5a, 0x1b, 0x76, 0x41, 0x3b, 0x98, 0x05, 0xbe,
	0xd5, 0x19, 0xfd, 0x1a, 0x60, 0xb0, 0x11, 0x29, 0xfa, 0xf7, 0xae, 0xd1, 0x4d, 0xe3, 0x56, 0xf3,
	0xe2, 0xb2, 0xb9, 0xf1, 0x3a, 0x3a, 0x08, 0x30, 0xb1, 0x25, 0x9c, 0x21, 0x91, 0x33, 0x9b, 0x51,
	0xae, 0xef, 0xae, 0xa2, 0xce, 0x12, 0xd9, 0x90, 0xfc, 0x9d, 0x7e, 0x3f, 0xcd, 0xd6, 0xc4, 0x75,
	0xbf, 0xcc, 0xfe, 0xa4, 0x94, 0x72, 0x63, 0x89, 0x4a, 0x58, 0xa0, 0xb5, 0x83, 0x8a, 0x88, 0xe4,
	0x77, 0xc5, 0x87, 0x6e, 0x2c, 0xc3, 0x9a, 0x0f, 0xd7, 0x4e, 0xa8, 0xf9, 0x1f, 0x73, 0x50, 0x16,
	0x31, 0x76, 0x63, 0x89, 0x71, 0xdd, 0x81, 0xf1, 0xda, 0xc1, 0x58, 0xf7, 0xbc, 0x1d, 0xeb, 0xc5,
	0x6e, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x13, 0xc5, 0xfd, 0x2b, 0x11, 0x02, 0x00, 0x00,
}