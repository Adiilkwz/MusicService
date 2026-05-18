package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	"github.com/Adiilkwz/music-grpc-go/auth"
)

func AuthMiddleware(authClient auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization heading"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect format of token (expected Bearer <token>)"})
			return
		}
		token := parts[1]

		req := &auth.ValidateTokenRequest{AccessToken: token}
		resp, err := authClient.ValidateToken(context.Background(), req)

		if err != nil || !resp.IsValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("user_id", resp.UserId)
		c.Set("user_role", resp.Role)

		c.Next()
	}
}

func GetGrpcContext(c *gin.Context) context.Context {
	md := metadata.Pairs()

	if userID, exists := c.Get("user_id"); exists {
		md.Set("user_id", fmt.Sprintf("%v", userID))
		fmt.Println("[GATEWAY] Success: Attached user_id to gRPC metadata:", userID)
	} else {
		fmt.Println("[GATEWAY] FATAL: user_id was completely missing from Gin context!")
	}

	if role, exists := c.Get("user_role"); exists {
		md.Set("role", fmt.Sprintf("%v", role))
	}

	return metadata.NewOutgoingContext(context.Background(), md)
}
