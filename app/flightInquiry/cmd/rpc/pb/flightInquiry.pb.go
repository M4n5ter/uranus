// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.2
// source: app/flightInquiry/cmd/rpc/pb/flightInquiry.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FlightInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 航班号，例如MU5735
	FlightNumber string `protobuf:"bytes,1,opt,name=FlightNumber,proto3" json:"FlightNumber,omitempty"`
	// 出发日期
	SetOutTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=SetOutTime,proto3" json:"SetOutTime,omitempty"`
	// 是否为头等舱/商务舱
	IsFirstClass bool `protobuf:"varint,3,opt,name=IsFirstClass,proto3" json:"IsFirstClass,omitempty"`
	// 票价
	Price uint64 `protobuf:"varint,4,opt,name=Price,proto3" json:"Price,omitempty"`
	// 剩余量(由于有超卖可能性，可能为负)
	Surplus int64 `protobuf:"varint,5,opt,name=Surplus,proto3" json:"Surplus,omitempty"`
	// 准点率(例如97，表示97%)
	Punctuality uint32 `protobuf:"varint,6,opt,name=Punctuality,proto3" json:"Punctuality,omitempty"`
	//起飞地点
	StartPosition string `protobuf:"bytes,7,opt,name=StartPosition,proto3" json:"StartPosition,omitempty"`
	//起飞时间
	StartTime *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=StartTime,proto3" json:"StartTime,omitempty"`
	//降落地点
	EndPosition string `protobuf:"bytes,9,opt,name=EndPosition,proto3" json:"EndPosition,omitempty"`
	//降落时间
	EndTime *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=EndTime,proto3" json:"EndTime,omitempty"`
}

func (x *FlightInfo) Reset() {
	*x = FlightInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlightInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlightInfo) ProtoMessage() {}

func (x *FlightInfo) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlightInfo.ProtoReflect.Descriptor instead.
func (*FlightInfo) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{0}
}

func (x *FlightInfo) GetFlightNumber() string {
	if x != nil {
		return x.FlightNumber
	}
	return ""
}

func (x *FlightInfo) GetSetOutTime() *timestamppb.Timestamp {
	if x != nil {
		return x.SetOutTime
	}
	return nil
}

func (x *FlightInfo) GetIsFirstClass() bool {
	if x != nil {
		return x.IsFirstClass
	}
	return false
}

func (x *FlightInfo) GetPrice() uint64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *FlightInfo) GetSurplus() int64 {
	if x != nil {
		return x.Surplus
	}
	return 0
}

func (x *FlightInfo) GetPunctuality() uint32 {
	if x != nil {
		return x.Punctuality
	}
	return 0
}

func (x *FlightInfo) GetStartPosition() string {
	if x != nil {
		return x.StartPosition
	}
	return ""
}

func (x *FlightInfo) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *FlightInfo) GetEndPosition() string {
	if x != nil {
		return x.EndPosition
	}
	return ""
}

func (x *FlightInfo) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

// QuireBySetOutTimeAndFlightNumberReq 通过给定日期、航班号进行航班查询请求
type QuireBySetOutTimeAndFlightNumberReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 航班号，例如MU5735
	FlightNumber string `protobuf:"bytes,1,opt,name=FlightNumber,proto3" json:"FlightNumber,omitempty"`
	// 出发日期
	SetOutTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=SetOutTime,proto3" json:"SetOutTime,omitempty"`
}

func (x *QuireBySetOutTimeAndFlightNumberReq) Reset() {
	*x = QuireBySetOutTimeAndFlightNumberReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutTimeAndFlightNumberReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutTimeAndFlightNumberReq) ProtoMessage() {}

func (x *QuireBySetOutTimeAndFlightNumberReq) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutTimeAndFlightNumberReq.ProtoReflect.Descriptor instead.
func (*QuireBySetOutTimeAndFlightNumberReq) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{1}
}

func (x *QuireBySetOutTimeAndFlightNumberReq) GetFlightNumber() string {
	if x != nil {
		return x.FlightNumber
	}
	return ""
}

func (x *QuireBySetOutTimeAndFlightNumberReq) GetSetOutTime() *timestamppb.Timestamp {
	if x != nil {
		return x.SetOutTime
	}
	return nil
}

// QuireBySetOutTimeAndFlightNumberResp 为 QuireBySetOutTimeAndFlightNumberReq 的响应
type QuireBySetOutTimeAndFlightNumberResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 查询结果
	FlightInfos []*FlightInfo `protobuf:"bytes,1,rep,name=FlightInfos,proto3" json:"FlightInfos,omitempty"`
}

func (x *QuireBySetOutTimeAndFlightNumberResp) Reset() {
	*x = QuireBySetOutTimeAndFlightNumberResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutTimeAndFlightNumberResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutTimeAndFlightNumberResp) ProtoMessage() {}

func (x *QuireBySetOutTimeAndFlightNumberResp) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutTimeAndFlightNumberResp.ProtoReflect.Descriptor instead.
func (*QuireBySetOutTimeAndFlightNumberResp) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{2}
}

func (x *QuireBySetOutTimeAndFlightNumberResp) GetFlightInfos() []*FlightInfo {
	if x != nil {
		return x.FlightInfos
	}
	return nil
}

// QuireBySetOutTimeStartPositionEndPositionReq 通过给定日期、出发地、目的地进行航班查询请求
type QuireBySetOutTimeStartPositionEndPositionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 出发日期
	SetOutTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=SetOutTime,proto3" json:"SetOutTime,omitempty"`
	//起飞地点
	StartPosition string `protobuf:"bytes,2,opt,name=StartPosition,proto3" json:"StartPosition,omitempty"`
	//降落地点
	EndPosition string `protobuf:"bytes,3,opt,name=EndPosition,proto3" json:"EndPosition,omitempty"`
}

func (x *QuireBySetOutTimeStartPositionEndPositionReq) Reset() {
	*x = QuireBySetOutTimeStartPositionEndPositionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutTimeStartPositionEndPositionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutTimeStartPositionEndPositionReq) ProtoMessage() {}

func (x *QuireBySetOutTimeStartPositionEndPositionReq) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutTimeStartPositionEndPositionReq.ProtoReflect.Descriptor instead.
func (*QuireBySetOutTimeStartPositionEndPositionReq) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{3}
}

func (x *QuireBySetOutTimeStartPositionEndPositionReq) GetSetOutTime() *timestamppb.Timestamp {
	if x != nil {
		return x.SetOutTime
	}
	return nil
}

func (x *QuireBySetOutTimeStartPositionEndPositionReq) GetStartPosition() string {
	if x != nil {
		return x.StartPosition
	}
	return ""
}

func (x *QuireBySetOutTimeStartPositionEndPositionReq) GetEndPosition() string {
	if x != nil {
		return x.EndPosition
	}
	return ""
}

// QuireBySetOutTimeStartPositionEndPositionResp 为 QuireBySetOutTimeStartPositionEndPositionReq 的响应
type QuireBySetOutTimeStartPositionEndPositionResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 查询结果
	FlightInfos []*FlightInfo `protobuf:"bytes,1,rep,name=FlightInfos,proto3" json:"FlightInfos,omitempty"`
}

func (x *QuireBySetOutTimeStartPositionEndPositionResp) Reset() {
	*x = QuireBySetOutTimeStartPositionEndPositionResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutTimeStartPositionEndPositionResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutTimeStartPositionEndPositionResp) ProtoMessage() {}

func (x *QuireBySetOutTimeStartPositionEndPositionResp) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutTimeStartPositionEndPositionResp.ProtoReflect.Descriptor instead.
func (*QuireBySetOutTimeStartPositionEndPositionResp) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{4}
}

func (x *QuireBySetOutTimeStartPositionEndPositionResp) GetFlightInfos() []*FlightInfo {
	if x != nil {
		return x.FlightInfos
	}
	return nil
}

var File_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto protoreflect.FileDescriptor

var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDesc = []byte{
	0x0a, 0x30, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x71, 0x75,
	0x69, 0x72, 0x79, 0x2f, 0x63, 0x6d, 0x64, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x66,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x71, 0x75, 0x69, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x03, 0x0a, 0x0a, 0x46, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x46, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x0a, 0x53, 0x65,
	0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x53, 0x65, 0x74, 0x4f,
	0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x49, 0x73, 0x46, 0x69, 0x72, 0x73,
	0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x49, 0x73,
	0x46, 0x69, 0x72, 0x73, 0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x72, 0x70, 0x6c, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x07, 0x53, 0x75, 0x72, 0x70, 0x6c, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x75,
	0x6e, 0x63, 0x74, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0b, 0x50, 0x75, 0x6e, 0x63, 0x74, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x24, 0x0a, 0x0d,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x34,
	0x0a, 0x07, 0x45, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x45, 0x6e, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0x85, 0x01, 0x0a, 0x23, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79,
	0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69,
	0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x22, 0x0a, 0x0c,
	0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x3a, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x58, 0x0a, 0x24,
	0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x30, 0x0a, 0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x46,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x22, 0xb2, 0x01, 0x0a, 0x2c, 0x51, 0x75, 0x69, 0x72, 0x65,
	0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x3a, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x45, 0x6e, 0x64,
	0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x61, 0x0a, 0x2d, 0x51,
	0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64,
	0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x30, 0x0a, 0x0b,
	0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x32, 0x99,
	0x02, 0x0a, 0x0d, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x71, 0x75, 0x69, 0x72, 0x79,
	0x12, 0x75, 0x0a, 0x20, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x27, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42,
	0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x28, 0x2e,
	0x70, 0x62, 0x2e, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x90, 0x01, 0x0a, 0x29, 0x51, 0x75, 0x69, 0x72,
	0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x69, 0x72, 0x65,
	0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x31, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x69,
	0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescOnce sync.Once
	file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescData = file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDesc
)

func file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP() []byte {
	file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescOnce.Do(func() {
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescData)
	})
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescData
}

var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_goTypes = []interface{}{
	(*FlightInfo)(nil),                                    // 0: pb.FlightInfo
	(*QuireBySetOutTimeAndFlightNumberReq)(nil),           // 1: pb.QuireBySetOutTimeAndFlightNumberReq
	(*QuireBySetOutTimeAndFlightNumberResp)(nil),          // 2: pb.QuireBySetOutTimeAndFlightNumberResp
	(*QuireBySetOutTimeStartPositionEndPositionReq)(nil),  // 3: pb.QuireBySetOutTimeStartPositionEndPositionReq
	(*QuireBySetOutTimeStartPositionEndPositionResp)(nil), // 4: pb.QuireBySetOutTimeStartPositionEndPositionResp
	(*timestamppb.Timestamp)(nil),                         // 5: google.protobuf.Timestamp
}
var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_depIdxs = []int32{
	5, // 0: pb.FlightInfo.SetOutTime:type_name -> google.protobuf.Timestamp
	5, // 1: pb.FlightInfo.StartTime:type_name -> google.protobuf.Timestamp
	5, // 2: pb.FlightInfo.EndTime:type_name -> google.protobuf.Timestamp
	5, // 3: pb.QuireBySetOutTimeAndFlightNumberReq.SetOutTime:type_name -> google.protobuf.Timestamp
	0, // 4: pb.QuireBySetOutTimeAndFlightNumberResp.FlightInfos:type_name -> pb.FlightInfo
	5, // 5: pb.QuireBySetOutTimeStartPositionEndPositionReq.SetOutTime:type_name -> google.protobuf.Timestamp
	0, // 6: pb.QuireBySetOutTimeStartPositionEndPositionResp.FlightInfos:type_name -> pb.FlightInfo
	1, // 7: pb.flightInquiry.QuireBySetOutTimeAndFlightNumber:input_type -> pb.QuireBySetOutTimeAndFlightNumberReq
	3, // 8: pb.flightInquiry.QuireBySetOutTimeStartPositionEndPosition:input_type -> pb.QuireBySetOutTimeStartPositionEndPositionReq
	2, // 9: pb.flightInquiry.QuireBySetOutTimeAndFlightNumber:output_type -> pb.QuireBySetOutTimeAndFlightNumberResp
	4, // 10: pb.flightInquiry.QuireBySetOutTimeStartPositionEndPosition:output_type -> pb.QuireBySetOutTimeStartPositionEndPositionResp
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_init() }
func file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_init() {
	if File_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlightInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutTimeAndFlightNumberReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutTimeAndFlightNumberResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutTimeStartPositionEndPositionReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutTimeStartPositionEndPositionResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_goTypes,
		DependencyIndexes: file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_depIdxs,
		MessageInfos:      file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes,
	}.Build()
	File_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto = out.File
	file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDesc = nil
	file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_goTypes = nil
	file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_depIdxs = nil
}
