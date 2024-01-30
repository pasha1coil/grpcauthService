package auth

import (
	"context"
	auth1 "github.com/pasha1coil/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string,
		password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, phone string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverApi struct {
	auth1.UnimplementedAuthServer
	auth Auth
}

const emptyvalue = 0

func RegisterServerApi(gRPC *grpc.Server, auth Auth) {
	auth1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) Login(ctx context.Context, req *auth1.LoginRequest) (*auth1.LoginResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyvalue {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(ctx context.Context, req *auth1.RegisterRequest) (*auth1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {

		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetPhone() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword(), req.GetPhone())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *auth1.IsAdminRequest) (*auth1.IsAdminResponse, error) {
	if req.GetUserId() == emptyvalue {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
