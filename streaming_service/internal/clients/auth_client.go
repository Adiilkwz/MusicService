package clients

import (
	"context"
	"fmt"

	authpb "github.com/Adiilkwz/music-grpc-go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

func NewAuthClient(authServiceAddress string) (*AuthClient, error) {
	conn, err := grpc.Dial(authServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		client: authpb.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*authpb.ValidateTokenResponse, error) {
	req := &authpb.ValidateTokenRequest{AccessToken: token}
	return c.client.ValidateToken(ctx, req)
}
