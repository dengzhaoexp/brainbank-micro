// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: userService.proto

package userService

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for UserService service

func NewUserServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for UserService service

type UserService interface {
	// 注册
	UserRegister(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*UserServiceResponse, error)
	// 登录
	UserLogin(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*UserLoginResponse, error)
	// 用户空间
	UserStorage(ctx context.Context, in *StorageRequest, opts ...client.CallOption) (*StorageResponse, error)
	// 更新密码
	UpdatePassword(ctx context.Context, in *UpdatePwdRequest, opts ...client.CallOption) (*UserServiceResponse, error)
	// 重设密码
	ResetPassword(ctx context.Context, in *ResetPwdRequest, opts ...client.CallOption) (*UserServiceResponse, error)
	// 发送邮件
	SendEmail(ctx context.Context, in *SendEmailRequest, opts ...client.CallOption) (*UserServiceResponse, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) UserRegister(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*UserServiceResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.UserRegister", in)
	out := new(UserServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UserLogin(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*UserLoginResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.UserLogin", in)
	out := new(UserLoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UserStorage(ctx context.Context, in *StorageRequest, opts ...client.CallOption) (*StorageResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.UserStorage", in)
	out := new(StorageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UpdatePassword(ctx context.Context, in *UpdatePwdRequest, opts ...client.CallOption) (*UserServiceResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.UpdatePassword", in)
	out := new(UserServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) ResetPassword(ctx context.Context, in *ResetPwdRequest, opts ...client.CallOption) (*UserServiceResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.ResetPassword", in)
	out := new(UserServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) SendEmail(ctx context.Context, in *SendEmailRequest, opts ...client.CallOption) (*UserServiceResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.SendEmail", in)
	out := new(UserServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	// 注册
	UserRegister(context.Context, *RegisterRequest, *UserServiceResponse) error
	// 登录
	UserLogin(context.Context, *LoginRequest, *UserLoginResponse) error
	// 用户空间
	UserStorage(context.Context, *StorageRequest, *StorageResponse) error
	// 更新密码
	UpdatePassword(context.Context, *UpdatePwdRequest, *UserServiceResponse) error
	// 重设密码
	ResetPassword(context.Context, *ResetPwdRequest, *UserServiceResponse) error
	// 发送邮件
	SendEmail(context.Context, *SendEmailRequest, *UserServiceResponse) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		UserRegister(ctx context.Context, in *RegisterRequest, out *UserServiceResponse) error
		UserLogin(ctx context.Context, in *LoginRequest, out *UserLoginResponse) error
		UserStorage(ctx context.Context, in *StorageRequest, out *StorageResponse) error
		UpdatePassword(ctx context.Context, in *UpdatePwdRequest, out *UserServiceResponse) error
		ResetPassword(ctx context.Context, in *ResetPwdRequest, out *UserServiceResponse) error
		SendEmail(ctx context.Context, in *SendEmailRequest, out *UserServiceResponse) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) UserRegister(ctx context.Context, in *RegisterRequest, out *UserServiceResponse) error {
	return h.UserServiceHandler.UserRegister(ctx, in, out)
}

func (h *userServiceHandler) UserLogin(ctx context.Context, in *LoginRequest, out *UserLoginResponse) error {
	return h.UserServiceHandler.UserLogin(ctx, in, out)
}

func (h *userServiceHandler) UserStorage(ctx context.Context, in *StorageRequest, out *StorageResponse) error {
	return h.UserServiceHandler.UserStorage(ctx, in, out)
}

func (h *userServiceHandler) UpdatePassword(ctx context.Context, in *UpdatePwdRequest, out *UserServiceResponse) error {
	return h.UserServiceHandler.UpdatePassword(ctx, in, out)
}

func (h *userServiceHandler) ResetPassword(ctx context.Context, in *ResetPwdRequest, out *UserServiceResponse) error {
	return h.UserServiceHandler.ResetPassword(ctx, in, out)
}

func (h *userServiceHandler) SendEmail(ctx context.Context, in *SendEmailRequest, out *UserServiceResponse) error {
	return h.UserServiceHandler.SendEmail(ctx, in, out)
}
