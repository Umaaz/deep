// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: deep.proto

package deeppb

import (
	v1 "github.com/intergral/deep/pkg/deeppb/tracepoint/v1"
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

type PushSnapshotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Snapshot *v1.Snapshot `protobuf:"bytes,1,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
}

func (x *PushSnapshotRequest) Reset() {
	*x = PushSnapshotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deep_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushSnapshotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushSnapshotRequest) ProtoMessage() {}

func (x *PushSnapshotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_deep_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushSnapshotRequest.ProtoReflect.Descriptor instead.
func (*PushSnapshotRequest) Descriptor() ([]byte, []int) {
	return file_deep_proto_rawDescGZIP(), []int{0}
}

func (x *PushSnapshotRequest) GetSnapshot() *v1.Snapshot {
	if x != nil {
		return x.Snapshot
	}
	return nil
}

// Write
type PushSnapshotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushSnapshotResponse) Reset() {
	*x = PushSnapshotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deep_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushSnapshotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushSnapshotResponse) ProtoMessage() {}

func (x *PushSnapshotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_deep_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushSnapshotResponse.ProtoReflect.Descriptor instead.
func (*PushSnapshotResponse) Descriptor() ([]byte, []int) {
	return file_deep_proto_rawDescGZIP(), []int{1}
}

type PushBytesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Snapshot []byte `protobuf:"bytes,1,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
	Id       []byte `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PushBytesRequest) Reset() {
	*x = PushBytesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deep_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushBytesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushBytesRequest) ProtoMessage() {}

func (x *PushBytesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_deep_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushBytesRequest.ProtoReflect.Descriptor instead.
func (*PushBytesRequest) Descriptor() ([]byte, []int) {
	return file_deep_proto_rawDescGZIP(), []int{2}
}

func (x *PushBytesRequest) GetSnapshot() []byte {
	if x != nil {
		return x.Snapshot
	}
	return nil
}

func (x *PushBytesRequest) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

type PushBytesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushBytesResponse) Reset() {
	*x = PushBytesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deep_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushBytesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushBytesResponse) ProtoMessage() {}

func (x *PushBytesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_deep_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushBytesResponse.ProtoReflect.Descriptor instead.
func (*PushBytesResponse) Descriptor() ([]byte, []int) {
	return file_deep_proto_rawDescGZIP(), []int{3}
}

var File_deep_proto protoreflect.FileDescriptor

var file_deep_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x64, 0x65, 0x65, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x64, 0x65,
	0x65, 0x70, 0x70, 0x62, 0x1a, 0x1e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a, 0x13, 0x50, 0x75, 0x73, 0x68, 0x53, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x08, 0x73,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x64, 0x65, 0x65, 0x70, 0x70, 0x62, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x08, 0x73,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x22, 0x16, 0x0a, 0x14, 0x50, 0x75, 0x73, 0x68, 0x53,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x3e, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x13, 0x0a, 0x11, 0x50, 0x75, 0x73, 0x68, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x5f, 0x0a, 0x10, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x4b, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x1b, 0x2e, 0x64, 0x65, 0x65, 0x70, 0x70,
	0x62, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x64, 0x65, 0x65, 0x70, 0x70, 0x62, 0x2e, 0x50,
	0x75, 0x73, 0x68, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x55, 0x0a, 0x0f, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x50, 0x75, 0x73, 0x68,
	0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x18, 0x2e, 0x64, 0x65, 0x65, 0x70, 0x70, 0x62, 0x2e, 0x50,
	0x75, 0x73, 0x68, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x19, 0x2e, 0x64, 0x65, 0x65, 0x70, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x42, 0x79, 0x74,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x26, 0x5a, 0x24,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x67, 0x72, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x65, 0x70, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x64, 0x65,
	0x65, 0x70, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_deep_proto_rawDescOnce sync.Once
	file_deep_proto_rawDescData = file_deep_proto_rawDesc
)

func file_deep_proto_rawDescGZIP() []byte {
	file_deep_proto_rawDescOnce.Do(func() {
		file_deep_proto_rawDescData = protoimpl.X.CompressGZIP(file_deep_proto_rawDescData)
	})
	return file_deep_proto_rawDescData
}

var file_deep_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_deep_proto_goTypes = []interface{}{
	(*PushSnapshotRequest)(nil),  // 0: deeppb.PushSnapshotRequest
	(*PushSnapshotResponse)(nil), // 1: deeppb.PushSnapshotResponse
	(*PushBytesRequest)(nil),     // 2: deeppb.PushBytesRequest
	(*PushBytesResponse)(nil),    // 3: deeppb.PushBytesResponse
	(*v1.Snapshot)(nil),          // 4: deeppb.tracepoint.v1.Snapshot
}
var file_deep_proto_depIdxs = []int32{
	4, // 0: deeppb.PushSnapshotRequest.snapshot:type_name -> deeppb.tracepoint.v1.Snapshot
	0, // 1: deeppb.MetricsGenerator.PushSnapshot:input_type -> deeppb.PushSnapshotRequest
	2, // 2: deeppb.IngesterService.PushBytes:input_type -> deeppb.PushBytesRequest
	1, // 3: deeppb.MetricsGenerator.PushSnapshot:output_type -> deeppb.PushSnapshotResponse
	3, // 4: deeppb.IngesterService.PushBytes:output_type -> deeppb.PushBytesResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_deep_proto_init() }
func file_deep_proto_init() {
	if File_deep_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_deep_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushSnapshotRequest); i {
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
		file_deep_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushSnapshotResponse); i {
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
		file_deep_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushBytesRequest); i {
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
		file_deep_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushBytesResponse); i {
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
			RawDescriptor: file_deep_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_deep_proto_goTypes,
		DependencyIndexes: file_deep_proto_depIdxs,
		MessageInfos:      file_deep_proto_msgTypes,
	}.Build()
	File_deep_proto = out.File
	file_deep_proto_rawDesc = nil
	file_deep_proto_goTypes = nil
	file_deep_proto_depIdxs = nil
}
