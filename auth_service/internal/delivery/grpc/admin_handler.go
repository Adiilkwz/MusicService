package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) ListUsers(ctx context.Context, req *auth.ListUsersRequest) (*auth.ListUsersResponse, error) {
	users, err := s.adminUC.ListUsers(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var grpcUsers []*auth.UserItem
	for _, u := range users {
		grpcUsers = append(grpcUsers, &auth.UserItem{
			Id:          u.ID,
			Email:       u.Email,
			DisplayName: u.DisplayName,
			Role:        u.Role,
		})
	}

	return &auth.ListUsersResponse{Users: grpcUsers}, nil
}

func (s *AuthServer) UpdateUserRole(ctx context.Context, req *auth.UpdateUserRoleRequest) (*auth.SuccessResponse, error) {
	err := s.adminUC.UpdateUserRole(ctx, req.AdminId, req.TargetUserId, req.NewRole)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return &auth.SuccessResponse{Success: true, Message: "User role updated"}, nil
}
