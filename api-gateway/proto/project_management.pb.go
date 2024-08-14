// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: project_management.proto

package __

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

type ProjectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId string `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
}

func (x *ProjectRequest) Reset() {
	*x = ProjectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_project_management_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectRequest) ProtoMessage() {}

func (x *ProjectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_project_management_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectRequest.ProtoReflect.Descriptor instead.
func (*ProjectRequest) Descriptor() ([]byte, []int) {
	return file_project_management_proto_rawDescGZIP(), []int{0}
}

func (x *ProjectRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type ProjectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectName        string `protobuf:"bytes,1,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	ProjectDescription string `protobuf:"bytes,2,opt,name=project_description,json=projectDescription,proto3" json:"project_description,omitempty"`
}

func (x *ProjectResponse) Reset() {
	*x = ProjectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_project_management_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectResponse) ProtoMessage() {}

func (x *ProjectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_project_management_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectResponse.ProtoReflect.Descriptor instead.
func (*ProjectResponse) Descriptor() ([]byte, []int) {
	return file_project_management_proto_rawDescGZIP(), []int{1}
}

func (x *ProjectResponse) GetProjectName() string {
	if x != nil {
		return x.ProjectName
	}
	return ""
}

func (x *ProjectResponse) GetProjectDescription() string {
	if x != nil {
		return x.ProjectDescription
	}
	return ""
}

var File_project_management_proto protoreflect.FileDescriptor

var file_project_management_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x2f,
	0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x22,
	0x65, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0x71, 0x0a, 0x18, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_project_management_proto_rawDescOnce sync.Once
	file_project_management_proto_rawDescData = file_project_management_proto_rawDesc
)

func file_project_management_proto_rawDescGZIP() []byte {
	file_project_management_proto_rawDescOnce.Do(func() {
		file_project_management_proto_rawDescData = protoimpl.X.CompressGZIP(file_project_management_proto_rawDescData)
	})
	return file_project_management_proto_rawDescData
}

var file_project_management_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_project_management_proto_goTypes = []any{
	(*ProjectRequest)(nil),  // 0: project_management.ProjectRequest
	(*ProjectResponse)(nil), // 1: project_management.ProjectResponse
}
var file_project_management_proto_depIdxs = []int32{
	0, // 0: project_management.ProjectManagementService.GetProject:input_type -> project_management.ProjectRequest
	1, // 1: project_management.ProjectManagementService.GetProject:output_type -> project_management.ProjectResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_project_management_proto_init() }
func file_project_management_proto_init() {
	if File_project_management_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_project_management_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectRequest); i {
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
		file_project_management_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectResponse); i {
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
			RawDescriptor: file_project_management_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_project_management_proto_goTypes,
		DependencyIndexes: file_project_management_proto_depIdxs,
		MessageInfos:      file_project_management_proto_msgTypes,
	}.Build()
	File_project_management_proto = out.File
	file_project_management_proto_rawDesc = nil
	file_project_management_proto_goTypes = nil
	file_project_management_proto_depIdxs = nil
}
