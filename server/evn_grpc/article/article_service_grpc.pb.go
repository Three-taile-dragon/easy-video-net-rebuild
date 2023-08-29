// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: article_service.proto

package article

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

// ArticleServiceClient is the client API for ArticleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArticleServiceClient interface {
	GetArticleContributionList(ctx context.Context, in *CommonPageInfo, opts ...grpc.CallOption) (*CommonDataResponse, error)
	GetArticleContributionListByUser(ctx context.Context, in *CommonIDRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	GetArticleComment(ctx context.Context, in *GetArticleCommentRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	GetArticleClassificationList(ctx context.Context, in *CommonZeroRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	GetArticleTotalInfo(ctx context.Context, in *CommonZeroRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	GetArticleContributionByID(ctx context.Context, in *CommonIDAndUIDRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	CreateArticleContribution(ctx context.Context, in *CreateArticleContributionRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	UpdateArticleContribution(ctx context.Context, in *UpdateArticleContributionRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	DeleteArticleByID(ctx context.Context, in *CommonIDAndUIDRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	ArticlePostComment(ctx context.Context, in *ArticlePostCommentRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
	GetArticleManagementList(ctx context.Context, in *GetArticleManagementListRequest, opts ...grpc.CallOption) (*CommonDataResponse, error)
}

type articleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArticleServiceClient(cc grpc.ClientConnInterface) ArticleServiceClient {
	return &articleServiceClient{cc}
}

func (c *articleServiceClient) GetArticleContributionList(ctx context.Context, in *CommonPageInfo, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleContributionList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetArticleContributionListByUser(ctx context.Context, in *CommonIDRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleContributionListByUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetArticleComment(ctx context.Context, in *GetArticleCommentRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetArticleClassificationList(ctx context.Context, in *CommonZeroRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleClassificationList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetArticleTotalInfo(ctx context.Context, in *CommonZeroRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleTotalInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetArticleContributionByID(ctx context.Context, in *CommonIDAndUIDRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleContributionByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) CreateArticleContribution(ctx context.Context, in *CreateArticleContributionRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/CreateArticleContribution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) UpdateArticleContribution(ctx context.Context, in *UpdateArticleContributionRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/UpdateArticleContribution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) DeleteArticleByID(ctx context.Context, in *CommonIDAndUIDRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/DeleteArticleByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) ArticlePostComment(ctx context.Context, in *ArticlePostCommentRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/ArticlePostComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetArticleManagementList(ctx context.Context, in *GetArticleManagementListRequest, opts ...grpc.CallOption) (*CommonDataResponse, error) {
	out := new(CommonDataResponse)
	err := c.cc.Invoke(ctx, "/article.service.v1.ArticleService/GetArticleManagementList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArticleServiceServer is the server API for ArticleService service.
// All implementations must embed UnimplementedArticleServiceServer
// for forward compatibility
type ArticleServiceServer interface {
	GetArticleContributionList(context.Context, *CommonPageInfo) (*CommonDataResponse, error)
	GetArticleContributionListByUser(context.Context, *CommonIDRequest) (*CommonDataResponse, error)
	GetArticleComment(context.Context, *GetArticleCommentRequest) (*CommonDataResponse, error)
	GetArticleClassificationList(context.Context, *CommonZeroRequest) (*CommonDataResponse, error)
	GetArticleTotalInfo(context.Context, *CommonZeroRequest) (*CommonDataResponse, error)
	GetArticleContributionByID(context.Context, *CommonIDAndUIDRequest) (*CommonDataResponse, error)
	CreateArticleContribution(context.Context, *CreateArticleContributionRequest) (*CommonDataResponse, error)
	UpdateArticleContribution(context.Context, *UpdateArticleContributionRequest) (*CommonDataResponse, error)
	DeleteArticleByID(context.Context, *CommonIDAndUIDRequest) (*CommonDataResponse, error)
	ArticlePostComment(context.Context, *ArticlePostCommentRequest) (*CommonDataResponse, error)
	GetArticleManagementList(context.Context, *GetArticleManagementListRequest) (*CommonDataResponse, error)
	mustEmbedUnimplementedArticleServiceServer()
}

// UnimplementedArticleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedArticleServiceServer struct {
}

func (UnimplementedArticleServiceServer) GetArticleContributionList(context.Context, *CommonPageInfo) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleContributionList not implemented")
}
func (UnimplementedArticleServiceServer) GetArticleContributionListByUser(context.Context, *CommonIDRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleContributionListByUser not implemented")
}
func (UnimplementedArticleServiceServer) GetArticleComment(context.Context, *GetArticleCommentRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleComment not implemented")
}
func (UnimplementedArticleServiceServer) GetArticleClassificationList(context.Context, *CommonZeroRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleClassificationList not implemented")
}
func (UnimplementedArticleServiceServer) GetArticleTotalInfo(context.Context, *CommonZeroRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleTotalInfo not implemented")
}
func (UnimplementedArticleServiceServer) GetArticleContributionByID(context.Context, *CommonIDAndUIDRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleContributionByID not implemented")
}
func (UnimplementedArticleServiceServer) CreateArticleContribution(context.Context, *CreateArticleContributionRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArticleContribution not implemented")
}
func (UnimplementedArticleServiceServer) UpdateArticleContribution(context.Context, *UpdateArticleContributionRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateArticleContribution not implemented")
}
func (UnimplementedArticleServiceServer) DeleteArticleByID(context.Context, *CommonIDAndUIDRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArticleByID not implemented")
}
func (UnimplementedArticleServiceServer) ArticlePostComment(context.Context, *ArticlePostCommentRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArticlePostComment not implemented")
}
func (UnimplementedArticleServiceServer) GetArticleManagementList(context.Context, *GetArticleManagementListRequest) (*CommonDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleManagementList not implemented")
}
func (UnimplementedArticleServiceServer) mustEmbedUnimplementedArticleServiceServer() {}

// UnsafeArticleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArticleServiceServer will
// result in compilation errors.
type UnsafeArticleServiceServer interface {
	mustEmbedUnimplementedArticleServiceServer()
}

func RegisterArticleServiceServer(s grpc.ServiceRegistrar, srv ArticleServiceServer) {
	s.RegisterService(&ArticleService_ServiceDesc, srv)
}

func _ArticleService_GetArticleContributionList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonPageInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleContributionList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleContributionList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleContributionList(ctx, req.(*CommonPageInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetArticleContributionListByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleContributionListByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleContributionListByUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleContributionListByUser(ctx, req.(*CommonIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetArticleComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArticleCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleComment(ctx, req.(*GetArticleCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetArticleClassificationList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonZeroRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleClassificationList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleClassificationList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleClassificationList(ctx, req.(*CommonZeroRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetArticleTotalInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonZeroRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleTotalInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleTotalInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleTotalInfo(ctx, req.(*CommonZeroRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetArticleContributionByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonIDAndUIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleContributionByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleContributionByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleContributionByID(ctx, req.(*CommonIDAndUIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_CreateArticleContribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateArticleContributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).CreateArticleContribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/CreateArticleContribution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).CreateArticleContribution(ctx, req.(*CreateArticleContributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_UpdateArticleContribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateArticleContributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).UpdateArticleContribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/UpdateArticleContribution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).UpdateArticleContribution(ctx, req.(*UpdateArticleContributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_DeleteArticleByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonIDAndUIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).DeleteArticleByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/DeleteArticleByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).DeleteArticleByID(ctx, req.(*CommonIDAndUIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_ArticlePostComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticlePostCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).ArticlePostComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/ArticlePostComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).ArticlePostComment(ctx, req.(*ArticlePostCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetArticleManagementList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArticleManagementListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetArticleManagementList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/article.service.v1.ArticleService/GetArticleManagementList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetArticleManagementList(ctx, req.(*GetArticleManagementListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ArticleService_ServiceDesc is the grpc.ServiceDesc for ArticleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArticleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "article.service.v1.ArticleService",
	HandlerType: (*ArticleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetArticleContributionList",
			Handler:    _ArticleService_GetArticleContributionList_Handler,
		},
		{
			MethodName: "GetArticleContributionListByUser",
			Handler:    _ArticleService_GetArticleContributionListByUser_Handler,
		},
		{
			MethodName: "GetArticleComment",
			Handler:    _ArticleService_GetArticleComment_Handler,
		},
		{
			MethodName: "GetArticleClassificationList",
			Handler:    _ArticleService_GetArticleClassificationList_Handler,
		},
		{
			MethodName: "GetArticleTotalInfo",
			Handler:    _ArticleService_GetArticleTotalInfo_Handler,
		},
		{
			MethodName: "GetArticleContributionByID",
			Handler:    _ArticleService_GetArticleContributionByID_Handler,
		},
		{
			MethodName: "CreateArticleContribution",
			Handler:    _ArticleService_CreateArticleContribution_Handler,
		},
		{
			MethodName: "UpdateArticleContribution",
			Handler:    _ArticleService_UpdateArticleContribution_Handler,
		},
		{
			MethodName: "DeleteArticleByID",
			Handler:    _ArticleService_DeleteArticleByID_Handler,
		},
		{
			MethodName: "ArticlePostComment",
			Handler:    _ArticleService_ArticlePostComment_Handler,
		},
		{
			MethodName: "GetArticleManagementList",
			Handler:    _ArticleService_GetArticleManagementList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "article_service.proto",
}
