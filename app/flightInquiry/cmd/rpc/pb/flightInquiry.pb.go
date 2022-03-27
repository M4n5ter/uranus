// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
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

//航班信息
type FlightInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//航班号，例如MU5735
	FlightNumber string `protobuf:"bytes,1,opt,name=FlightNumber,proto3" json:"FlightNumber,omitempty"`
	//出发日期
	SetOutDate *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=SetOutDate,proto3" json:"SetOutDate,omitempty"`
	//是否为头等舱/商务舱
	IsFirstClass bool `protobuf:"varint,3,opt,name=IsFirstClass,proto3" json:"IsFirstClass,omitempty"`
	//票价(￥)
	Price uint64 `protobuf:"varint,4,opt,name=Price,proto3" json:"Price,omitempty"`
	//折扣(-n%)
	Discount int64 `protobuf:"varint,5,opt,name=Discount,proto3" json:"Discount,omitempty"`
	//剩余量(由于有超卖可能性，可能为负)
	Surplus int64 `protobuf:"varint,6,opt,name=Surplus,proto3" json:"Surplus,omitempty"`
	//准点率(例如97，表示97%)
	Punctuality uint32 `protobuf:"varint,7,opt,name=Punctuality,proto3" json:"Punctuality,omitempty"`
	//起飞地点
	DepartPosition string `protobuf:"bytes,8,opt,name=DepartPosition,proto3" json:"DepartPosition,omitempty"`
	//起飞时间
	DepartTime *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=DepartTime,proto3" json:"DepartTime,omitempty"`
	//降落地点
	ArrivePosition string `protobuf:"bytes,10,opt,name=ArrivePosition,proto3" json:"ArrivePosition,omitempty"`
	//降落时间
	ArriveTime *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=ArriveTime,proto3" json:"ArriveTime,omitempty"`
	//退票信息
	RefundInfo *RefundInfo `protobuf:"bytes,12,opt,name=RefundInfo,proto3" json:"RefundInfo,omitempty"`
	//改票信息
	ChangeInfo *ChangeInfo `protobuf:"bytes,13,opt,name=ChangeInfo,proto3" json:"ChangeInfo,omitempty"`
	//托运行李额(KG)
	Cba int64 `protobuf:"varint,14,opt,name=Cba,proto3" json:"Cba,omitempty"`
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

func (x *FlightInfo) GetSetOutDate() *timestamppb.Timestamp {
	if x != nil {
		return x.SetOutDate
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

func (x *FlightInfo) GetDiscount() int64 {
	if x != nil {
		return x.Discount
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

func (x *FlightInfo) GetDepartPosition() string {
	if x != nil {
		return x.DepartPosition
	}
	return ""
}

func (x *FlightInfo) GetDepartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DepartTime
	}
	return nil
}

func (x *FlightInfo) GetArrivePosition() string {
	if x != nil {
		return x.ArrivePosition
	}
	return ""
}

func (x *FlightInfo) GetArriveTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ArriveTime
	}
	return nil
}

func (x *FlightInfo) GetRefundInfo() *RefundInfo {
	if x != nil {
		return x.RefundInfo
	}
	return nil
}

func (x *FlightInfo) GetChangeInfo() *ChangeInfo {
	if x != nil {
		return x.ChangeInfo
	}
	return nil
}

func (x *FlightInfo) GetCba() int64 {
	if x != nil {
		return x.Cba
	}
	return 0
}

//时间-费用
type TimeFee struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=Time,proto3" json:"Time,omitempty"`
	Fee  uint64                 `protobuf:"varint,2,opt,name=Fee,proto3" json:"Fee,omitempty"`
}

func (x *TimeFee) Reset() {
	*x = TimeFee{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeFee) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeFee) ProtoMessage() {}

func (x *TimeFee) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TimeFee.ProtoReflect.Descriptor instead.
func (*TimeFee) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{1}
}

func (x *TimeFee) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *TimeFee) GetFee() uint64 {
	if x != nil {
		return x.Fee
	}
	return 0
}

//退票信息
type RefundInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeFees []*TimeFee `protobuf:"bytes,1,rep,name=TimeFees,proto3" json:"TimeFees,omitempty"`
}

func (x *RefundInfo) Reset() {
	*x = RefundInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefundInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefundInfo) ProtoMessage() {}

func (x *RefundInfo) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use RefundInfo.ProtoReflect.Descriptor instead.
func (*RefundInfo) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{2}
}

func (x *RefundInfo) GetTimeFees() []*TimeFee {
	if x != nil {
		return x.TimeFees
	}
	return nil
}

//改票信息
type ChangeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeFees []*TimeFee `protobuf:"bytes,1,rep,name=TimeFees,proto3" json:"TimeFees,omitempty"`
}

func (x *ChangeInfo) Reset() {
	*x = ChangeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeInfo) ProtoMessage() {}

func (x *ChangeInfo) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ChangeInfo.ProtoReflect.Descriptor instead.
func (*ChangeInfo) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{3}
}

func (x *ChangeInfo) GetTimeFees() []*TimeFee {
	if x != nil {
		return x.TimeFees
	}
	return nil
}

//QuireBySetOutDateAndFlightNumberReq 通过给定日期、航班号进行航班查询请求
type QuireBySetOutDateAndFlightNumberReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 航班号，例如MU5735
	FlightNumber string `protobuf:"bytes,1,opt,name=FlightNumber,proto3" json:"FlightNumber,omitempty"`
	// 出发日期
	SetOutDate *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=SetOutDate,proto3" json:"SetOutDate,omitempty"`
}

func (x *QuireBySetOutDateAndFlightNumberReq) Reset() {
	*x = QuireBySetOutDateAndFlightNumberReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutDateAndFlightNumberReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutDateAndFlightNumberReq) ProtoMessage() {}

func (x *QuireBySetOutDateAndFlightNumberReq) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use QuireBySetOutDateAndFlightNumberReq.ProtoReflect.Descriptor instead.
func (*QuireBySetOutDateAndFlightNumberReq) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{4}
}

func (x *QuireBySetOutDateAndFlightNumberReq) GetFlightNumber() string {
	if x != nil {
		return x.FlightNumber
	}
	return ""
}

func (x *QuireBySetOutDateAndFlightNumberReq) GetSetOutDate() *timestamppb.Timestamp {
	if x != nil {
		return x.SetOutDate
	}
	return nil
}

//QuireBySetOutDateAndFlightNumberResp 为 QuireBySetOutDateAndFlightNumberReq 的响应
type QuireBySetOutDateAndFlightNumberResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 查询结果
	FlightInfos []*FlightInfo `protobuf:"bytes,1,rep,name=FlightInfos,proto3" json:"FlightInfos,omitempty"`
}

func (x *QuireBySetOutDateAndFlightNumberResp) Reset() {
	*x = QuireBySetOutDateAndFlightNumberResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutDateAndFlightNumberResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutDateAndFlightNumberResp) ProtoMessage() {}

func (x *QuireBySetOutDateAndFlightNumberResp) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutDateAndFlightNumberResp.ProtoReflect.Descriptor instead.
func (*QuireBySetOutDateAndFlightNumberResp) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{5}
}

func (x *QuireBySetOutDateAndFlightNumberResp) GetFlightInfos() []*FlightInfo {
	if x != nil {
		return x.FlightInfos
	}
	return nil
}

//QuireBySetOutDateStartPositionEndPositionReq 通过给定日期、出发地、目的地进行航班查询请求
type QuireBySetOutDateStartPositionEndPositionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 出发日期
	SetOutDate *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=SetOutDate,proto3" json:"SetOutDate,omitempty"`
	//起飞地点
	DepartPosition string `protobuf:"bytes,2,opt,name=DepartPosition,proto3" json:"DepartPosition,omitempty"`
	//降落地点
	ArrivePosition string `protobuf:"bytes,3,opt,name=ArrivePosition,proto3" json:"ArrivePosition,omitempty"`
}

func (x *QuireBySetOutDateStartPositionEndPositionReq) Reset() {
	*x = QuireBySetOutDateStartPositionEndPositionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutDateStartPositionEndPositionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutDateStartPositionEndPositionReq) ProtoMessage() {}

func (x *QuireBySetOutDateStartPositionEndPositionReq) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutDateStartPositionEndPositionReq.ProtoReflect.Descriptor instead.
func (*QuireBySetOutDateStartPositionEndPositionReq) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{6}
}

func (x *QuireBySetOutDateStartPositionEndPositionReq) GetSetOutDate() *timestamppb.Timestamp {
	if x != nil {
		return x.SetOutDate
	}
	return nil
}

func (x *QuireBySetOutDateStartPositionEndPositionReq) GetDepartPosition() string {
	if x != nil {
		return x.DepartPosition
	}
	return ""
}

func (x *QuireBySetOutDateStartPositionEndPositionReq) GetArrivePosition() string {
	if x != nil {
		return x.ArrivePosition
	}
	return ""
}

//QuireBySetOutDateStartPositionEndPositionResp 为 QuireBySetOutTimeStartPositionEndPositionReq 的响应
type QuireBySetOutDateStartPositionEndPositionResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 查询结果
	FlightInfos []*FlightInfo `protobuf:"bytes,1,rep,name=FlightInfos,proto3" json:"FlightInfos,omitempty"`
}

func (x *QuireBySetOutDateStartPositionEndPositionResp) Reset() {
	*x = QuireBySetOutDateStartPositionEndPositionResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuireBySetOutDateStartPositionEndPositionResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuireBySetOutDateStartPositionEndPositionResp) ProtoMessage() {}

func (x *QuireBySetOutDateStartPositionEndPositionResp) ProtoReflect() protoreflect.Message {
	mi := &file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuireBySetOutDateStartPositionEndPositionResp.ProtoReflect.Descriptor instead.
func (*QuireBySetOutDateStartPositionEndPositionResp) Descriptor() ([]byte, []int) {
	return file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_rawDescGZIP(), []int{7}
}

func (x *QuireBySetOutDateStartPositionEndPositionResp) GetFlightInfos() []*FlightInfo {
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
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8, 0x04, 0x0a, 0x0a, 0x46, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x46, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x0a, 0x53, 0x65,
	0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x53, 0x65, 0x74, 0x4f,
	0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x49, 0x73, 0x46, 0x69, 0x72, 0x73,
	0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x49, 0x73,
	0x46, 0x69, 0x72, 0x73, 0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x53, 0x75, 0x72, 0x70, 0x6c, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x53,
	0x75, 0x72, 0x70, 0x6c, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x75, 0x6e, 0x63, 0x74, 0x75,
	0x61, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x50, 0x75, 0x6e,
	0x63, 0x74, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x44, 0x65, 0x70, 0x61,
	0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x3a, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0a, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e,
	0x41, 0x72, 0x72, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x41, 0x72, 0x72, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x0a, 0x41, 0x72, 0x72, 0x69, 0x76, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x41, 0x72, 0x72, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x2e, 0x0a, 0x0a, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x2e, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x10, 0x0a, 0x03, 0x43, 0x62, 0x61, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x43,
	0x62, 0x61, 0x22, 0x4b, 0x0a, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x65, 0x65, 0x12, 0x2e, 0x0a,
	0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x46, 0x65, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x46, 0x65, 0x65, 0x22,
	0x35, 0x0a, 0x0a, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a,
	0x08, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x65, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x65, 0x65, 0x52, 0x08, 0x54, 0x69,
	0x6d, 0x65, 0x46, 0x65, 0x65, 0x73, 0x22, 0x35, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a, 0x08, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x65, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x46, 0x65, 0x65, 0x52, 0x08, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x65, 0x65, 0x73, 0x22, 0x85, 0x01,
	0x0a, 0x23, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44,
	0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x22, 0x0a, 0x0c, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x46, 0x6c, 0x69,
	0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x0a, 0x53, 0x65, 0x74,
	0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75,
	0x74, 0x44, 0x61, 0x74, 0x65, 0x22, 0x58, 0x0a, 0x24, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79,
	0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69,
	0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x30, 0x0a,
	0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x22,
	0xba, 0x01, 0x0a, 0x2c, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75,
	0x74, 0x44, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x12, 0x3a, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x0e,
	0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x41, 0x72, 0x72, 0x69, 0x76, 0x65, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x41, 0x72,
	0x72, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x61, 0x0a, 0x2d,
	0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e,
	0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x30, 0x0a,
	0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x0b, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x32,
	0x99, 0x02, 0x0a, 0x0d, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x71, 0x75, 0x69, 0x72,
	0x79, 0x12, 0x75, 0x0a, 0x20, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f,
	0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x27, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x69, 0x72, 0x65,
	0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x46,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x28,
	0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75,
	0x74, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x90, 0x01, 0x0a, 0x29, 0x51, 0x75, 0x69,
	0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x69, 0x72,
	0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x31, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75,
	0x69, 0x72, 0x65, 0x42, 0x79, 0x53, 0x65, 0x74, 0x4f, 0x75, 0x74, 0x44, 0x61, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x42, 0x06, 0x5a, 0x04, 0x2e,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_goTypes = []interface{}{
	(*FlightInfo)(nil), // 0: pb.FlightInfo
	(*TimeFee)(nil),    // 1: pb.TimeFee
	(*RefundInfo)(nil), // 2: pb.RefundInfo
	(*ChangeInfo)(nil), // 3: pb.ChangeInfo
	(*QuireBySetOutDateAndFlightNumberReq)(nil),           // 4: pb.QuireBySetOutDateAndFlightNumberReq
	(*QuireBySetOutDateAndFlightNumberResp)(nil),          // 5: pb.QuireBySetOutDateAndFlightNumberResp
	(*QuireBySetOutDateStartPositionEndPositionReq)(nil),  // 6: pb.QuireBySetOutDateStartPositionEndPositionReq
	(*QuireBySetOutDateStartPositionEndPositionResp)(nil), // 7: pb.QuireBySetOutDateStartPositionEndPositionResp
	(*timestamppb.Timestamp)(nil),                         // 8: google.protobuf.Timestamp
}
var file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_depIdxs = []int32{
	8,  // 0: pb.FlightInfo.SetOutDate:type_name -> google.protobuf.Timestamp
	8,  // 1: pb.FlightInfo.DepartTime:type_name -> google.protobuf.Timestamp
	8,  // 2: pb.FlightInfo.ArriveTime:type_name -> google.protobuf.Timestamp
	2,  // 3: pb.FlightInfo.RefundInfo:type_name -> pb.RefundInfo
	3,  // 4: pb.FlightInfo.ChangeInfo:type_name -> pb.ChangeInfo
	8,  // 5: pb.TimeFee.Time:type_name -> google.protobuf.Timestamp
	1,  // 6: pb.RefundInfo.TimeFees:type_name -> pb.TimeFee
	1,  // 7: pb.ChangeInfo.TimeFees:type_name -> pb.TimeFee
	8,  // 8: pb.QuireBySetOutDateAndFlightNumberReq.SetOutDate:type_name -> google.protobuf.Timestamp
	0,  // 9: pb.QuireBySetOutDateAndFlightNumberResp.FlightInfos:type_name -> pb.FlightInfo
	8,  // 10: pb.QuireBySetOutDateStartPositionEndPositionReq.SetOutDate:type_name -> google.protobuf.Timestamp
	0,  // 11: pb.QuireBySetOutDateStartPositionEndPositionResp.FlightInfos:type_name -> pb.FlightInfo
	4,  // 12: pb.flightInquiry.QuireBySetOutDateAndFlightNumber:input_type -> pb.QuireBySetOutDateAndFlightNumberReq
	6,  // 13: pb.flightInquiry.QuireBySetOutDateStartPositionEndPosition:input_type -> pb.QuireBySetOutDateStartPositionEndPositionReq
	5,  // 14: pb.flightInquiry.QuireBySetOutDateAndFlightNumber:output_type -> pb.QuireBySetOutDateAndFlightNumberResp
	7,  // 15: pb.flightInquiry.QuireBySetOutDateStartPositionEndPosition:output_type -> pb.QuireBySetOutDateStartPositionEndPositionResp
	14, // [14:16] is the sub-list for method output_type
	12, // [12:14] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
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
			switch v := v.(*TimeFee); i {
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
			switch v := v.(*RefundInfo); i {
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
			switch v := v.(*ChangeInfo); i {
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
			switch v := v.(*QuireBySetOutDateAndFlightNumberReq); i {
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
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutDateAndFlightNumberResp); i {
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
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutDateStartPositionEndPositionReq); i {
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
		file_app_flightInquiry_cmd_rpc_pb_flightInquiry_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuireBySetOutDateStartPositionEndPositionResp); i {
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
			NumMessages:   8,
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