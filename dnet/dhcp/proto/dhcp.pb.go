// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.6
// source: dhcp/proto/dhcp.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Resp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	// Will contain either an error message or success message
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Resp) Reset() {
	*x = Resp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dhcp_proto_dhcp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resp) ProtoMessage() {}

func (x *Resp) ProtoReflect() protoreflect.Message {
	mi := &file_dhcp_proto_dhcp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resp.ProtoReflect.Descriptor instead.
func (*Resp) Descriptor() ([]byte, []int) {
	return file_dhcp_proto_dhcp_proto_rawDescGZIP(), []int{0}
}

func (x *Resp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *Resp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type StartReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Networks    []*Network    `protobuf:"bytes,1,rep,name=networks,proto3" json:"networks,omitempty"`
	StaticHosts []*StaticHost `protobuf:"bytes,2,rep,name=staticHosts,proto3" json:"staticHosts,omitempty"`
}

func (x *StartReq) Reset() {
	*x = StartReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dhcp_proto_dhcp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartReq) ProtoMessage() {}

func (x *StartReq) ProtoReflect() protoreflect.Message {
	mi := &file_dhcp_proto_dhcp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartReq.ProtoReflect.Descriptor instead.
func (*StartReq) Descriptor() ([]byte, []int) {
	return file_dhcp_proto_dhcp_proto_rawDescGZIP(), []int{1}
}

func (x *StartReq) GetNetworks() []*Network {
	if x != nil {
		return x.Networks
	}
	return nil
}

func (x *StartReq) GetStaticHosts() []*StaticHost {
	if x != nil {
		return x.StaticHosts
	}
	return nil
}

type Network struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Router    string `protobuf:"bytes,1,opt,name=router,proto3" json:"router,omitempty"`
	Network   string `protobuf:"bytes,2,opt,name=network,proto3" json:"network,omitempty"`
	Min       string `protobuf:"bytes,3,opt,name=min,proto3" json:"min,omitempty"`
	Max       string `protobuf:"bytes,4,opt,name=max,proto3" json:"max,omitempty"`
	DnsServer string `protobuf:"bytes,5,opt,name=dns_server,json=dnsServer,proto3" json:"dns_server,omitempty"`
}

func (x *Network) Reset() {
	*x = Network{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dhcp_proto_dhcp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Network) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Network) ProtoMessage() {}

func (x *Network) ProtoReflect() protoreflect.Message {
	mi := &file_dhcp_proto_dhcp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Network.ProtoReflect.Descriptor instead.
func (*Network) Descriptor() ([]byte, []int) {
	return file_dhcp_proto_dhcp_proto_rawDescGZIP(), []int{2}
}

func (x *Network) GetRouter() string {
	if x != nil {
		return x.Router
	}
	return ""
}

func (x *Network) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Network) GetMin() string {
	if x != nil {
		return x.Min
	}
	return ""
}

func (x *Network) GetMax() string {
	if x != nil {
		return x.Max
	}
	return ""
}

func (x *Network) GetDnsServer() string {
	if x != nil {
		return x.DnsServer
	}
	return ""
}

type StaticHost struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address    string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	MacAddress string `protobuf:"bytes,3,opt,name=macAddress,proto3" json:"macAddress,omitempty"`
	DnsServer  string `protobuf:"bytes,4,opt,name=dns_server,json=dnsServer,proto3" json:"dns_server,omitempty"`
	DomainName string `protobuf:"bytes,5,opt,name=domain_name,json=domainName,proto3" json:"domain_name,omitempty"`
	Router     string `protobuf:"bytes,6,opt,name=router,proto3" json:"router,omitempty"`
}

func (x *StaticHost) Reset() {
	*x = StaticHost{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dhcp_proto_dhcp_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StaticHost) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticHost) ProtoMessage() {}

func (x *StaticHost) ProtoReflect() protoreflect.Message {
	mi := &file_dhcp_proto_dhcp_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticHost.ProtoReflect.Descriptor instead.
func (*StaticHost) Descriptor() ([]byte, []int) {
	return file_dhcp_proto_dhcp_proto_rawDescGZIP(), []int{3}
}

func (x *StaticHost) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StaticHost) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *StaticHost) GetMacAddress() string {
	if x != nil {
		return x.MacAddress
	}
	return ""
}

func (x *StaticHost) GetDnsServer() string {
	if x != nil {
		return x.DnsServer
	}
	return ""
}

func (x *StaticHost) GetDomainName() string {
	if x != nil {
		return x.DomainName
	}
	return ""
}

func (x *StaticHost) GetRouter() string {
	if x != nil {
		return x.Router
	}
	return ""
}

type StopReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopReq) Reset() {
	*x = StopReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dhcp_proto_dhcp_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopReq) ProtoMessage() {}

func (x *StopReq) ProtoReflect() protoreflect.Message {
	mi := &file_dhcp_proto_dhcp_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopReq.ProtoReflect.Descriptor instead.
func (*StopReq) Descriptor() ([]byte, []int) {
	return file_dhcp_proto_dhcp_proto_rawDescGZIP(), []int{4}
}

type UpdateReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Networks    []*Network    `protobuf:"bytes,1,rep,name=networks,proto3" json:"networks,omitempty"`
	StaticHosts []*StaticHost `protobuf:"bytes,2,rep,name=staticHosts,proto3" json:"staticHosts,omitempty"`
}

func (x *UpdateReq) Reset() {
	*x = UpdateReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dhcp_proto_dhcp_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateReq) ProtoMessage() {}

func (x *UpdateReq) ProtoReflect() protoreflect.Message {
	mi := &file_dhcp_proto_dhcp_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateReq.ProtoReflect.Descriptor instead.
func (*UpdateReq) Descriptor() ([]byte, []int) {
	return file_dhcp_proto_dhcp_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateReq) GetNetworks() []*Network {
	if x != nil {
		return x.Networks
	}
	return nil
}

func (x *UpdateReq) GetStaticHosts() []*StaticHost {
	if x != nil {
		return x.StaticHosts
	}
	return nil
}

var File_dhcp_proto_dhcp_proto protoreflect.FileDescriptor

var file_dhcp_proto_dhcp_proto_rawDesc = []byte{
	0x0a, 0x15, 0x64, 0x68, 0x63, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x68, 0x63,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x04, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x5f, 0x0a, 0x08, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x12,
	0x24, 0x0a, 0x08, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x08, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x73, 0x12, 0x2d, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x48,
	0x6f, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x48,
	0x6f, 0x73, 0x74, 0x73, 0x22, 0x7e, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x6e, 0x73, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x6e, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x22, 0xb2, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x48,
	0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x61, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x61, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x6e, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x6e, 0x73, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x22, 0x09, 0x0a, 0x07, 0x53, 0x74, 0x6f,
	0x70, 0x52, 0x65, 0x71, 0x22, 0x60, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x24, 0x0a, 0x08, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x08, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x12, 0x2d, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x69,
	0x63, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x69,
	0x63, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x32, 0x69, 0x0a, 0x04, 0x44, 0x48, 0x43, 0x50, 0x12, 0x1f,
	0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x44, 0x48, 0x43, 0x50, 0x12, 0x09, 0x2e, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x05, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12,
	0x1d, 0x0a, 0x08, 0x53, 0x74, 0x6f, 0x70, 0x44, 0x48, 0x43, 0x50, 0x12, 0x08, 0x2e, 0x53, 0x74,
	0x6f, 0x70, 0x52, 0x65, 0x71, 0x1a, 0x05, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x21,
	0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x48, 0x43, 0x50, 0x12, 0x0a, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x05, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x22,
	0x00, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x61, 0x75, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2d, 0x73, 0x65, 0x63, 0x75,
	0x72, 0x69, 0x74, 0x79, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x74, 0x74, 0x2f, 0x64, 0x6e, 0x65, 0x74,
	0x2f, 0x64, 0x68, 0x63, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_dhcp_proto_dhcp_proto_rawDescOnce sync.Once
	file_dhcp_proto_dhcp_proto_rawDescData = file_dhcp_proto_dhcp_proto_rawDesc
)

func file_dhcp_proto_dhcp_proto_rawDescGZIP() []byte {
	file_dhcp_proto_dhcp_proto_rawDescOnce.Do(func() {
		file_dhcp_proto_dhcp_proto_rawDescData = protoimpl.X.CompressGZIP(file_dhcp_proto_dhcp_proto_rawDescData)
	})
	return file_dhcp_proto_dhcp_proto_rawDescData
}

var file_dhcp_proto_dhcp_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_dhcp_proto_dhcp_proto_goTypes = []interface{}{
	(*Resp)(nil),       // 0: Resp
	(*StartReq)(nil),   // 1: StartReq
	(*Network)(nil),    // 2: Network
	(*StaticHost)(nil), // 3: StaticHost
	(*StopReq)(nil),    // 4: StopReq
	(*UpdateReq)(nil),  // 5: UpdateReq
}
var file_dhcp_proto_dhcp_proto_depIdxs = []int32{
	2, // 0: StartReq.networks:type_name -> Network
	3, // 1: StartReq.staticHosts:type_name -> StaticHost
	2, // 2: UpdateReq.networks:type_name -> Network
	3, // 3: UpdateReq.staticHosts:type_name -> StaticHost
	1, // 4: DHCP.StartDHCP:input_type -> StartReq
	4, // 5: DHCP.StopDHCP:input_type -> StopReq
	5, // 6: DHCP.UpdateDHCP:input_type -> UpdateReq
	0, // 7: DHCP.StartDHCP:output_type -> Resp
	0, // 8: DHCP.StopDHCP:output_type -> Resp
	0, // 9: DHCP.UpdateDHCP:output_type -> Resp
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_dhcp_proto_dhcp_proto_init() }
func file_dhcp_proto_dhcp_proto_init() {
	if File_dhcp_proto_dhcp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dhcp_proto_dhcp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resp); i {
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
		file_dhcp_proto_dhcp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartReq); i {
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
		file_dhcp_proto_dhcp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Network); i {
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
		file_dhcp_proto_dhcp_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StaticHost); i {
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
		file_dhcp_proto_dhcp_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopReq); i {
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
		file_dhcp_proto_dhcp_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateReq); i {
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
			RawDescriptor: file_dhcp_proto_dhcp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dhcp_proto_dhcp_proto_goTypes,
		DependencyIndexes: file_dhcp_proto_dhcp_proto_depIdxs,
		MessageInfos:      file_dhcp_proto_dhcp_proto_msgTypes,
	}.Build()
	File_dhcp_proto_dhcp_proto = out.File
	file_dhcp_proto_dhcp_proto_rawDesc = nil
	file_dhcp_proto_dhcp_proto_goTypes = nil
	file_dhcp_proto_dhcp_proto_depIdxs = nil
}
