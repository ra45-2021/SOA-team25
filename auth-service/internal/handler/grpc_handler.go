package handler

import (
	"context"
	"auth-service/internal/service"
	"auth-service/pb"
)

type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewAuthGRPCHandler(svc *service.AuthService) *AuthGRPCHandler {
	return &AuthGRPCHandler{svc: svc}
}

func (h *AuthGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	u, found, err := h.svc.GetByID(ctx, req.Id)
	if err != nil || !found {
		return nil, err
	}
	return &pb.UserResponse{Id: u.ID, Username: u.Username, Role: u.Role}, nil
}

func (h *AuthGRPCHandler) CheckAuthor(ctx context.Context, req *pb.CheckAuthorRequest) (*pb.CheckAuthorResponse, error) {
	u, found, err := h.svc.GetByID(ctx, req.Id)
	if err != nil || !found {
		return &pb.CheckAuthorResponse{Exists: false}, nil
	}
	return &pb.CheckAuthorResponse{Exists: true, Role: u.Role}, nil
}