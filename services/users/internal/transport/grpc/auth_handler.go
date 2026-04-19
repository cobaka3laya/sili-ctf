package grpc

import (
	"context"

	pb "github.com/cobaka3laya/sili-ctf/services/users/gen/proto"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)

func (h *AuthHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error)
