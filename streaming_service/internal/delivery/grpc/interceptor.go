package grpc

import (
	"context"
	"strings"

	"streaming_service/internal/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(authClient *clients.AuthClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		token := strings.TrimPrefix(values[0], "Bearer ")

		resp, err := authClient.ValidateToken(ctx, token)
		if err != nil || !resp.IsValid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}
