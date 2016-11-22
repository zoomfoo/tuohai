// Code generated by protoc-gen-go.
// source: internal/pb/IM_Login/IM.Login.proto
// DO NOT EDIT!

/*
Package IM_Login is a generated protocol buffer package.

It is generated from these files:
	internal/pb/IM_Login/IM.Login.proto

It has these top-level messages:
	IMLoginReq
	IMLoginRes
	IMLogoutReq
	IMLogoutRsp
*/
package IM_Login

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import IM_BaseDefine "internal/pb/IM_BaseDefine"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type IMLoginReq struct {
	// cmd id:		0x0103
	UserId        string                     `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	Passwd        string                     `protobuf:"bytes,2,opt,name=passwd" json:"passwd,omitempty"`
	OnlineStatus  IM_BaseDefine.UserStatType `protobuf:"varint,3,opt,name=online_status,json=onlineStatus,enum=IM.BaseDefine.UserStatType" json:"online_status,omitempty"`
	ClientType    IM_BaseDefine.ClientType   `protobuf:"varint,4,opt,name=client_type,json=clientType,enum=IM.BaseDefine.ClientType" json:"client_type,omitempty"`
	DeviceId      string                     `protobuf:"bytes,5,opt,name=device_id,json=deviceId" json:"device_id,omitempty"`
	ClientVersion string                     `protobuf:"bytes,6,opt,name=client_version,json=clientVersion" json:"client_version,omitempty"`
	Token         string                     `protobuf:"bytes,7,opt,name=token" json:"token,omitempty"`
}

func (m *IMLoginReq) Reset()                    { *m = IMLoginReq{} }
func (m *IMLoginReq) String() string            { return proto.CompactTextString(m) }
func (*IMLoginReq) ProtoMessage()               {}
func (*IMLoginReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type IMLoginRes struct {
	// cmd id:		0x0104
	ServerTime   uint32                     `protobuf:"varint,1,opt,name=server_time,json=serverTime" json:"server_time,omitempty"`
	ErrCode      IM_BaseDefine.ResultType   `protobuf:"varint,2,opt,name=err_code,json=errCode,enum=IM.BaseDefine.ResultType" json:"err_code,omitempty"`
	ErrMsg       string                     `protobuf:"bytes,3,opt,name=err_msg,json=errMsg" json:"err_msg,omitempty"`
	OnlineStatus IM_BaseDefine.UserStatType `protobuf:"varint,4,opt,name=online_status,json=onlineStatus,enum=IM.BaseDefine.UserStatType" json:"online_status,omitempty"`
}

func (m *IMLoginRes) Reset()                    { *m = IMLoginRes{} }
func (m *IMLoginRes) String() string            { return proto.CompactTextString(m) }
func (*IMLoginRes) ProtoMessage()               {}
func (*IMLoginRes) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type IMLogoutReq struct {
}

func (m *IMLogoutReq) Reset()                    { *m = IMLogoutReq{} }
func (m *IMLogoutReq) String() string            { return proto.CompactTextString(m) }
func (*IMLogoutReq) ProtoMessage()               {}
func (*IMLogoutReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type IMLogoutRsp struct {
	// cmd id:		0x0106
	ErrCode uint32 `protobuf:"varint,1,opt,name=err_code,json=errCode" json:"err_code,omitempty"`
}

func (m *IMLogoutRsp) Reset()                    { *m = IMLogoutRsp{} }
func (m *IMLogoutRsp) String() string            { return proto.CompactTextString(m) }
func (*IMLogoutRsp) ProtoMessage()               {}
func (*IMLogoutRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*IMLoginReq)(nil), "IM.Login.IMLoginReq")
	proto.RegisterType((*IMLoginRes)(nil), "IM.Login.IMLoginRes")
	proto.RegisterType((*IMLogoutReq)(nil), "IM.Login.IMLogoutReq")
	proto.RegisterType((*IMLogoutRsp)(nil), "IM.Login.IMLogoutRsp")
}

func init() { proto.RegisterFile("internal/pb/IM_Login/IM.Login.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x52, 0x4d, 0x4f, 0xea, 0x40,
	0x14, 0x0d, 0x5f, 0x05, 0x2e, 0xaf, 0x2c, 0x26, 0x2f, 0xef, 0x95, 0xc7, 0xe2, 0x99, 0x1a, 0x13,
	0x36, 0x42, 0xa2, 0xae, 0x5c, 0x19, 0x70, 0x21, 0x89, 0x6c, 0x0a, 0xba, 0x6d, 0x4a, 0x7b, 0x25,
	0x13, 0xcb, 0x4c, 0xed, 0x4c, 0x31, 0xfe, 0x10, 0xff, 0x8e, 0xbf, 0xcd, 0xdb, 0x19, 0x51, 0x84,
	0x95, 0xbb, 0x39, 0xe7, 0x9e, 0x73, 0x39, 0xf7, 0x50, 0x38, 0xe6, 0x42, 0x63, 0x2e, 0xa2, 0x74,
	0x94, 0x2d, 0x47, 0xd3, 0x59, 0x78, 0x2b, 0x57, 0x5c, 0xd0, 0x63, 0x68, 0x1e, 0xc3, 0x2c, 0x97,
	0x5a, 0xb2, 0xd6, 0x16, 0xff, 0x3b, 0xdd, 0x93, 0x8f, 0x23, 0x85, 0xd7, 0xf8, 0xc0, 0x05, 0x96,
	0x9e, 0x2f, 0x64, 0x8d, 0xfe, 0x6b, 0x15, 0x60, 0x3a, 0x33, 0xd6, 0x00, 0x9f, 0xd8, 0x5f, 0x68,
	0x16, 0x0a, 0xf3, 0x90, 0x27, 0x5e, 0xe5, 0xa8, 0x32, 0x68, 0x07, 0x4e, 0x09, 0xa7, 0x09, 0xfb,
	0x03, 0x4e, 0x16, 0x29, 0xf5, 0x9c, 0x78, 0x55, 0xcb, 0x5b, 0xc4, 0xae, 0xc0, 0x95, 0x22, 0xa5,
	0x7d, 0xa1, 0xd2, 0x91, 0x2e, 0x94, 0x57, 0xa3, 0x71, 0xf7, 0xac, 0x3f, 0xfc, 0xfe, 0x63, 0x77,
	0xb4, 0x65, 0x4e, 0x82, 0xc5, 0x4b, 0x86, 0xc1, 0x2f, 0xeb, 0x98, 0x1b, 0x03, 0xbb, 0x84, 0x4e,
	0x9c, 0x72, 0x14, 0x3a, 0xd4, 0x34, 0xf4, 0xea, 0xc6, 0xdf, 0xdb, 0xf3, 0x4f, 0x8c, 0xc2, 0xb8,
	0x21, 0xfe, 0x7c, 0xb3, 0x3e, 0xb4, 0x13, 0xdc, 0xf0, 0x18, 0xcb, 0xc0, 0x0d, 0x13, 0xac, 0x65,
	0x09, 0x8a, 0x7c, 0x02, 0xdd, 0x8f, 0xc5, 0x1b, 0xcc, 0x15, 0x97, 0xc2, 0x73, 0x8c, 0xc2, 0xb5,
	0xec, 0xbd, 0x25, 0xd9, 0x6f, 0x68, 0x68, 0xf9, 0x88, 0xc2, 0x6b, 0x9a, 0xa9, 0x05, 0xfe, 0x5b,
	0x65, 0xa7, 0x17, 0xc5, 0xfe, 0x43, 0x87, 0x2e, 0xa0, 0x3d, 0xa1, 0xe6, 0x6b, 0x34, 0xdd, 0xb8,
	0x01, 0x58, 0x6a, 0x41, 0x0c, 0xbb, 0x80, 0x16, 0xe6, 0x79, 0x18, 0xcb, 0x04, 0x4d, 0x43, 0x87,
	0x27, 0xd0, 0x9a, 0x22, 0xb5, 0x27, 0x34, 0x49, 0x3a, 0x21, 0x65, 0x59, 0x77, 0xe9, 0x5a, 0xab,
	0x95, 0xe9, 0x8d, 0x6a, 0x25, 0x38, 0x53, 0xab, 0xc3, 0x5a, 0xeb, 0x3f, 0xac, 0xd5, 0x77, 0xa1,
	0x63, 0xf2, 0xcb, 0x42, 0xd3, 0x1f, 0xeb, 0x0f, 0x76, 0xa0, 0xca, 0x58, 0x6f, 0x27, 0xae, 0x3d,
	0x66, 0x9b, 0x69, 0x5c, 0xbd, 0xa9, 0x2d, 0x1d, 0xf3, 0x71, 0x9c, 0xbf, 0x07, 0x00, 0x00, 0xff,
	0xff, 0x95, 0xab, 0x07, 0x5a, 0x7c, 0x02, 0x00, 0x00,
}
