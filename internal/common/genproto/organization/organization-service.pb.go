// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: api/protobuf/organization-service.proto

package organization

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *ByIDRequest) Reset() {
	*x = ByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_organization_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ByIDRequest) ProtoMessage() {}

func (x *ByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_organization_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ByIDRequest.ProtoReflect.Descriptor instead.
func (*ByIDRequest) Descriptor() ([]byte, []int) {
	return file_api_protobuf_organization_service_proto_rawDescGZIP(), []int{0}
}

func (x *ByIDRequest) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

var File_api_protobuf_organization_service_proto protoreflect.FileDescriptor

var file_api_protobuf_organization_service_proto_rawDesc = []byte{
	0x0a, 0x27, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6f,
	0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6f, 0x72, 0x67, 0x73, 0x76,
	0x63, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x1d, 0x0a, 0x0b, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49,
	0x44, 0x32, 0x4d, 0x0a, 0x13, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42,
	0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x13, 0x2e, 0x6f, 0x72, 0x67, 0x73, 0x76, 0x63, 0x2e, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68,
	0x42, 0x27, 0x5a, 0x25, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_protobuf_organization_service_proto_rawDescOnce sync.Once
	file_api_protobuf_organization_service_proto_rawDescData = file_api_protobuf_organization_service_proto_rawDesc
)

func file_api_protobuf_organization_service_proto_rawDescGZIP() []byte {
	file_api_protobuf_organization_service_proto_rawDescOnce.Do(func() {
		file_api_protobuf_organization_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_protobuf_organization_service_proto_rawDescData)
	})
	return file_api_protobuf_organization_service_proto_rawDescData
}

var file_api_protobuf_organization_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_protobuf_organization_service_proto_goTypes = []interface{}{
	(*ByIDRequest)(nil), // 0: orgsvc.ByIDRequest
	(*Branch)(nil),      // 1: organization.Branch
}
var file_api_protobuf_organization_service_proto_depIdxs = []int32{
	0, // 0: orgsvc.OrganizationService.GetBranch:input_type -> orgsvc.ByIDRequest
	1, // 1: orgsvc.OrganizationService.GetBranch:output_type -> organization.Branch
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_protobuf_organization_service_proto_init() }
func file_api_protobuf_organization_service_proto_init() {
	if File_api_protobuf_organization_service_proto != nil {
		return
	}
	file_api_protobuf_organization_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_protobuf_organization_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ByIDRequest); i {
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
			RawDescriptor: file_api_protobuf_organization_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_protobuf_organization_service_proto_goTypes,
		DependencyIndexes: file_api_protobuf_organization_service_proto_depIdxs,
		MessageInfos:      file_api_protobuf_organization_service_proto_msgTypes,
	}.Build()
	File_api_protobuf_organization_service_proto = out.File
	file_api_protobuf_organization_service_proto_rawDesc = nil
	file_api_protobuf_organization_service_proto_goTypes = nil
	file_api_protobuf_organization_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// OrganizationServiceClient is the client API for OrganizationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrganizationServiceClient interface {
	GetBranch(ctx context.Context, in *ByIDRequest, opts ...grpc.CallOption) (*Branch, error)
}

type organizationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrganizationServiceClient(cc grpc.ClientConnInterface) OrganizationServiceClient {
	return &organizationServiceClient{cc}
}

func (c *organizationServiceClient) GetBranch(ctx context.Context, in *ByIDRequest, opts ...grpc.CallOption) (*Branch, error) {
	out := new(Branch)
	err := c.cc.Invoke(ctx, "/orgsvc.OrganizationService/GetBranch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrganizationServiceServer is the server API for OrganizationService service.
type OrganizationServiceServer interface {
	GetBranch(context.Context, *ByIDRequest) (*Branch, error)
}

// UnimplementedOrganizationServiceServer can be embedded to have forward compatible implementations.
type UnimplementedOrganizationServiceServer struct {
}

func (*UnimplementedOrganizationServiceServer) GetBranch(context.Context, *ByIDRequest) (*Branch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBranch not implemented")
}

func RegisterOrganizationServiceServer(s *grpc.Server, srv OrganizationServiceServer) {
	s.RegisterService(&_OrganizationService_serviceDesc, srv)
}

func _OrganizationService_GetBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationServiceServer).GetBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/orgsvc.OrganizationService/GetBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationServiceServer).GetBranch(ctx, req.(*ByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrganizationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orgsvc.OrganizationService",
	HandlerType: (*OrganizationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBranch",
			Handler:    _OrganizationService_GetBranch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/protobuf/organization-service.proto",
}
