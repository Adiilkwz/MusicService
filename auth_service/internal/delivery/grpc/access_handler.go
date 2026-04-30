package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	id, err := s.accessUC.Register(ctx, req.Email, req.Password, req.DisplayName)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	return &auth.RegisterResponse{UserId: id}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	session, err := s.accessUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &auth.LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		UserId:       session.UserID,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	uid, role, err := s.accessUC.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &auth.ValidateTokenResponse{IsValid: false}, nil
	}
	return &auth.ValidateTokenResponse{IsValid: true, UserId: uid, Role: role}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	session, err := s.accessUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &auth.RefreshTokenResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (s *AuthServer) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.SuccessResponse, error) {
	err := s.accessUC.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &auth.SuccessResponse{Success: true, Message: "Password changed successfully"}, nil
}

func (s *AuthServer) SendPasswordReset(ctx context.Context, req *auth.SendResetRequest) (*auth.SuccessResponse, error) {
	err := s.accessUC.SendPasswordReset(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SuccessResponse{Success: true, Message: "Reset code sent"}, nil
}

func (s *AuthServer) ConfirmPasswordReset(ctx context.Context, req *auth.ConfirmResetRequest) (*auth.SuccessResponse, error) {
	err := s.accessUC.ConfirmPasswordReset(ctx, req.Email, req.ResetCode, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &auth.SuccessResponse{Success: true, Message: "Password reset confirmed"}, nil
}
