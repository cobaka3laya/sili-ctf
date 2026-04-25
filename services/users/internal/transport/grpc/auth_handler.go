package grpc

import (
	"context"

	pb "github.com/cobaka3laya/sili-ctf/services/users/gen/proto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/service"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)

func (h *AuthHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error)

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LoginResponse, error)
