package grpc

import (
	"context"
	"strconv"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func extractUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return 0, status.Error(codes.Unauthenticated, "metadata is not found")
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return 0, status.Error(codes.Unauthenticated, "user_id is missing in metadata")
	}

	userID, err := strconv.ParseInt(userIDs[0], 10, 64)
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "invalid user_id format in metadata")
	}

	return userID, nil
}

func (s *AuthServer) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.profileUC.GetProfile(ctx, userID)
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
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	err = s.profileUC.UpdateProfile(ctx, userID, req.DisplayName, req.AvatarUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.SuccessResponse{Success: true, Message: "Profile updated"}, nil
}

func (s *AuthServer) DeleteAccount(ctx context.Context, req *auth.DeleteAccountRequest) (*auth.SuccessResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	err = s.profileUC.DeleteAccount(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.SuccessResponse{Success: true, Message: "Account deleted"}, nil
}
