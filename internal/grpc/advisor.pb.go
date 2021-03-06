// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: advisor.proto

package grpc

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ReschedulingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time             string                         `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	CurrentPlacement *ReschedulingRequest_Placement `protobuf:"bytes,2,opt,name=currentPlacement,proto3" json:"currentPlacement,omitempty"`
	NewPlacement     *ReschedulingRequest_Placement `protobuf:"bytes,3,opt,name=newPlacement,proto3" json:"newPlacement,omitempty"`
}

func (x *ReschedulingRequest) Reset() {
	*x = ReschedulingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_advisor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReschedulingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReschedulingRequest) ProtoMessage() {}

func (x *ReschedulingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_advisor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReschedulingRequest.ProtoReflect.Descriptor instead.
func (*ReschedulingRequest) Descriptor() ([]byte, []int) {
	return file_advisor_proto_rawDescGZIP(), []int{0}
}

func (x *ReschedulingRequest) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

func (x *ReschedulingRequest) GetCurrentPlacement() *ReschedulingRequest_Placement {
	if x != nil {
		return x.CurrentPlacement
	}
	return nil
}

func (x *ReschedulingRequest) GetNewPlacement() *ReschedulingRequest_Placement {
	if x != nil {
		return x.NewPlacement
	}
	return nil
}

type ReconfigurationCost struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentCost int64 `protobuf:"varint,1,opt,name=currentCost,proto3" json:"currentCost,omitempty"`
	OptimalCost int64 `protobuf:"varint,2,opt,name=optimalCost,proto3" json:"optimalCost,omitempty"`
}

func (x *ReconfigurationCost) Reset() {
	*x = ReconfigurationCost{}
	if protoimpl.UnsafeEnabled {
		mi := &file_advisor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconfigurationCost) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconfigurationCost) ProtoMessage() {}

func (x *ReconfigurationCost) ProtoReflect() protoreflect.Message {
	mi := &file_advisor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconfigurationCost.ProtoReflect.Descriptor instead.
func (*ReconfigurationCost) Descriptor() ([]byte, []int) {
	return file_advisor_proto_rawDescGZIP(), []int{1}
}

func (x *ReconfigurationCost) GetCurrentCost() int64 {
	if x != nil {
		return x.CurrentCost
	}
	return 0
}

func (x *ReconfigurationCost) GetOptimalCost() int64 {
	if x != nil {
		return x.OptimalCost
	}
	return 0
}

type ReschedulingRequest_Placement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Assignments map[string]string `protobuf:"bytes,1,rep,name=assignments,proto3" json:"assignments,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ReschedulingRequest_Placement) Reset() {
	*x = ReschedulingRequest_Placement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_advisor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReschedulingRequest_Placement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReschedulingRequest_Placement) ProtoMessage() {}

func (x *ReschedulingRequest_Placement) ProtoReflect() protoreflect.Message {
	mi := &file_advisor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReschedulingRequest_Placement.ProtoReflect.Descriptor instead.
func (*ReschedulingRequest_Placement) Descriptor() ([]byte, []int) {
	return file_advisor_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ReschedulingRequest_Placement) GetAssignments() map[string]string {
	if x != nil {
		return x.Assignments
	}
	return nil
}

var File_advisor_proto protoreflect.FileDescriptor

var file_advisor_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x64, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x22, 0xf5, 0x02, 0x0a, 0x13, 0x52, 0x65,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x53, 0x0a, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50,
	0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x4b, 0x0a, 0x0c, 0x6e, 0x65,
	0x77, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0c, 0x6e, 0x65, 0x77, 0x50, 0x6c,
	0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0xa7, 0x01, 0x0a, 0x09, 0x50, 0x6c, 0x61, 0x63,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x5a, 0x0a, 0x0b, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x1a, 0x3e, 0x0a, 0x10, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x59, 0x0a, 0x13, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x70,
	0x74, 0x69, 0x6d, 0x61, 0x6c, 0x43, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x6f, 0x70, 0x74, 0x69, 0x6d, 0x61, 0x6c, 0x43, 0x6f, 0x73, 0x74, 0x32, 0x5e, 0x0a, 0x10,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x41, 0x64, 0x76, 0x69, 0x73, 0x6f, 0x72,
	0x12, 0x4a, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x1d, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x73, 0x74, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_advisor_proto_rawDescOnce sync.Once
	file_advisor_proto_rawDescData = file_advisor_proto_rawDesc
)

func file_advisor_proto_rawDescGZIP() []byte {
	file_advisor_proto_rawDescOnce.Do(func() {
		file_advisor_proto_rawDescData = protoimpl.X.CompressGZIP(file_advisor_proto_rawDescData)
	})
	return file_advisor_proto_rawDescData
}

var file_advisor_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_advisor_proto_goTypes = []interface{}{
	(*ReschedulingRequest)(nil),           // 0: internal.ReschedulingRequest
	(*ReconfigurationCost)(nil),           // 1: internal.ReconfigurationCost
	(*ReschedulingRequest_Placement)(nil), // 2: internal.ReschedulingRequest.Placement
	nil,                                   // 3: internal.ReschedulingRequest.Placement.AssignmentsEntry
}
var file_advisor_proto_depIdxs = []int32{
	2, // 0: internal.ReschedulingRequest.currentPlacement:type_name -> internal.ReschedulingRequest.Placement
	2, // 1: internal.ReschedulingRequest.newPlacement:type_name -> internal.ReschedulingRequest.Placement
	3, // 2: internal.ReschedulingRequest.Placement.assignments:type_name -> internal.ReschedulingRequest.Placement.AssignmentsEntry
	0, // 3: internal.SchedulerAdvisor.GetCosts:input_type -> internal.ReschedulingRequest
	1, // 4: internal.SchedulerAdvisor.GetCosts:output_type -> internal.ReconfigurationCost
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_advisor_proto_init() }
func file_advisor_proto_init() {
	if File_advisor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_advisor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReschedulingRequest); i {
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
		file_advisor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconfigurationCost); i {
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
		file_advisor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReschedulingRequest_Placement); i {
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
			RawDescriptor: file_advisor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_advisor_proto_goTypes,
		DependencyIndexes: file_advisor_proto_depIdxs,
		MessageInfos:      file_advisor_proto_msgTypes,
	}.Build()
	File_advisor_proto = out.File
	file_advisor_proto_rawDesc = nil
	file_advisor_proto_goTypes = nil
	file_advisor_proto_depIdxs = nil
}
