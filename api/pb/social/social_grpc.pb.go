// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.0--rc3
// source: social.proto

package social

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SocialService_FollowAction_FullMethodName    = "/social.SocialService/FollowAction"
	SocialService_GetFollowList_FullMethodName   = "/social.SocialService/GetFollowList"
	SocialService_GetFollowerList_FullMethodName = "/social.SocialService/GetFollowerList"
	SocialService_GetFriendList_FullMethodName   = "/social.SocialService/GetFriendList"
	SocialService_GetFollowInfo_FullMethodName   = "/social.SocialService/GetFollowInfo"
	SocialService_PostMessage_FullMethodName     = "/social.SocialService/PostMessage"
	SocialService_GetMessage_FullMethodName      = "/social.SocialService/GetMessage"
)

// SocialServiceClient is the client API for SocialService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SocialServiceClient interface {
	FollowAction(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowResponse, error)
	GetFollowList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error)
	GetFollowerList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error)
	GetFriendList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error)
	GetFollowInfo(ctx context.Context, in *FollowInfoRequest, opts ...grpc.CallOption) (*FollowInfoResponse, error)
	PostMessage(ctx context.Context, in *PostMessageRequest, opts ...grpc.CallOption) (*PostMessageResponse, error)
	GetMessage(ctx context.Context, in *GetMessageRequest, opts ...grpc.CallOption) (*GetMessageResponse, error)
}

type socialServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSocialServiceClient(cc grpc.ClientConnInterface) SocialServiceClient {
	return &socialServiceClient{cc}
}

func (c *socialServiceClient) FollowAction(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowResponse, error) {
	out := new(FollowResponse)
	err := c.cc.Invoke(ctx, SocialService_FollowAction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServiceClient) GetFollowList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error) {
	out := new(FollowListResponse)
	err := c.cc.Invoke(ctx, SocialService_GetFollowList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServiceClient) GetFollowerList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error) {
	out := new(FollowListResponse)
	err := c.cc.Invoke(ctx, SocialService_GetFollowerList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServiceClient) GetFriendList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error) {
	out := new(FollowListResponse)
	err := c.cc.Invoke(ctx, SocialService_GetFriendList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServiceClient) GetFollowInfo(ctx context.Context, in *FollowInfoRequest, opts ...grpc.CallOption) (*FollowInfoResponse, error) {
	out := new(FollowInfoResponse)
	err := c.cc.Invoke(ctx, SocialService_GetFollowInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServiceClient) PostMessage(ctx context.Context, in *PostMessageRequest, opts ...grpc.CallOption) (*PostMessageResponse, error) {
	out := new(PostMessageResponse)
	err := c.cc.Invoke(ctx, SocialService_PostMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServiceClient) GetMessage(ctx context.Context, in *GetMessageRequest, opts ...grpc.CallOption) (*GetMessageResponse, error) {
	out := new(GetMessageResponse)
	err := c.cc.Invoke(ctx, SocialService_GetMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SocialServiceServer is the server API for SocialService service.
// All implementations must embed UnimplementedSocialServiceServer
// for forward compatibility
type SocialServiceServer interface {
	FollowAction(context.Context, *FollowRequest) (*FollowResponse, error)
	GetFollowList(context.Context, *FollowListRequest) (*FollowListResponse, error)
	GetFollowerList(context.Context, *FollowListRequest) (*FollowListResponse, error)
	GetFriendList(context.Context, *FollowListRequest) (*FollowListResponse, error)
	GetFollowInfo(context.Context, *FollowInfoRequest) (*FollowInfoResponse, error)
	PostMessage(context.Context, *PostMessageRequest) (*PostMessageResponse, error)
	GetMessage(context.Context, *GetMessageRequest) (*GetMessageResponse, error)
	mustEmbedUnimplementedSocialServiceServer()
}

// UnimplementedSocialServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSocialServiceServer struct {
}

func (UnimplementedSocialServiceServer) FollowAction(context.Context, *FollowRequest) (*FollowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowAction not implemented")
}
func (UnimplementedSocialServiceServer) GetFollowList(context.Context, *FollowListRequest) (*FollowListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowList not implemented")
}
func (UnimplementedSocialServiceServer) GetFollowerList(context.Context, *FollowListRequest) (*FollowListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowerList not implemented")
}
func (UnimplementedSocialServiceServer) GetFriendList(context.Context, *FollowListRequest) (*FollowListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendList not implemented")
}
func (UnimplementedSocialServiceServer) GetFollowInfo(context.Context, *FollowInfoRequest) (*FollowInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowInfo not implemented")
}
func (UnimplementedSocialServiceServer) PostMessage(context.Context, *PostMessageRequest) (*PostMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostMessage not implemented")
}
func (UnimplementedSocialServiceServer) GetMessage(context.Context, *GetMessageRequest) (*GetMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedSocialServiceServer) mustEmbedUnimplementedSocialServiceServer() {}

// UnsafeSocialServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SocialServiceServer will
// result in compilation errors.
type UnsafeSocialServiceServer interface {
	mustEmbedUnimplementedSocialServiceServer()
}

func RegisterSocialServiceServer(s grpc.ServiceRegistrar, srv SocialServiceServer) {
	s.RegisterService(&SocialService_ServiceDesc, srv)
}

func _SocialService_FollowAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).FollowAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_FollowAction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).FollowAction(ctx, req.(*FollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SocialService_GetFollowList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).GetFollowList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_GetFollowList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).GetFollowList(ctx, req.(*FollowListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SocialService_GetFollowerList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).GetFollowerList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_GetFollowerList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).GetFollowerList(ctx, req.(*FollowListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SocialService_GetFriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).GetFriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_GetFriendList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).GetFriendList(ctx, req.(*FollowListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SocialService_GetFollowInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).GetFollowInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_GetFollowInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).GetFollowInfo(ctx, req.(*FollowInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SocialService_PostMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).PostMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_PostMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).PostMessage(ctx, req.(*PostMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SocialService_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServiceServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SocialService_GetMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServiceServer).GetMessage(ctx, req.(*GetMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SocialService_ServiceDesc is the grpc.ServiceDesc for SocialService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SocialService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "social.SocialService",
	HandlerType: (*SocialServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FollowAction",
			Handler:    _SocialService_FollowAction_Handler,
		},
		{
			MethodName: "GetFollowList",
			Handler:    _SocialService_GetFollowList_Handler,
		},
		{
			MethodName: "GetFollowerList",
			Handler:    _SocialService_GetFollowerList_Handler,
		},
		{
			MethodName: "GetFriendList",
			Handler:    _SocialService_GetFriendList_Handler,
		},
		{
			MethodName: "GetFollowInfo",
			Handler:    _SocialService_GetFollowInfo_Handler,
		},
		{
			MethodName: "PostMessage",
			Handler:    _SocialService_PostMessage_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _SocialService_GetMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "social.proto",
}