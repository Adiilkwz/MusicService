package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	user, err := s.profileUC.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	return &auth.GetProfileResponse{
		UserId:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Role:        user.Role,
	}, nil
}

func (s *AuthServer) UpdateProfile(ctx context.Context, req *auth.UpdateProfileRequest) (*auth.SuccessResponse, error) {
	err := s.profileUC.UpdateProfile(ctx, req.UserId, req.DisplayName, req.AvatarUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SuccessResponse{Success: true, Message: "Profile updated"}, nil
}

func (s *AuthServer) DeleteAccount(ctx context.Context, req *auth.DeleteAccountRequest) (*auth.SuccessResponse, error) {
	err := s.profileUC.DeleteAccount(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SuccessResponse{Success: true, Message: "Account deleted"}, nil
}
