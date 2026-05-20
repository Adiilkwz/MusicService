package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"api_gateway/internal/middleware"

	"github.com/Adiilkwz/music-grpc-go/auth"
)

type RegisterInput struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateRoleInput struct {
	NewRole string `json:"new_role" binding:"required"`
}

func RegisterPublicAuthRoutes(rg *gin.RouterGroup, client auth.AuthServiceClient) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", registerHandler(client))
		authGroup.POST("/login", loginHandler(client))
	}
}

func RegisterProtectedAuthRoutes(rg *gin.RouterGroup, client auth.AuthServiceClient) {
	profileGroup := rg.Group("/profile")
	{
		profileGroup.GET("/me", getProfileHandler(client))
		profileGroup.GET("/", getProfileHandler(client))
		profileGroup.PUT("/", updateProfileHandler(client))
		profileGroup.DELETE("/", deleteAccountHandler(client))
	}

	adminGroup := rg.Group("/admin/users")
	{
		adminGroup.GET("/", listUsersHandler(client))
		adminGroup.PUT("/:id/role", updateUserRoleHandler(client))
	}
}

func registerHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.Register(c.Request.Context(), &auth.RegisterRequest{
			Email:       input.Email,
			Password:    input.Password,
			DisplayName: input.DisplayName,
		})
		if err != nil {
			st, _ := status.FromError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user_id": resp.UserId, "message": "User registered successfully"})
	}
}

func loginHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}

		resp, err := client.Login(c.Request.Context(), &auth.LoginRequest{
			Email:    input.Email,
			Password: input.Password,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		})
	}
}

func getProfileHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := middleware.GetGrpcContext(c)

		resp, err := client.GetProfile(ctx, &auth.GetProfileRequest{})
		if err != nil {
			st, _ := status.FromError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func listUsersHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := middleware.GetGrpcContext(c)

		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
		offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

		resp, err := client.ListUsers(ctx, &auth.ListUsersRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() == codes.PermissionDenied {
				c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusOK, resp.Users)
	}
}

func updateUserRoleHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := middleware.GetGrpcContext(c)

		targetUserIDStr := c.Param("id")
		targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var input UpdateRoleInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New role is required"})
			return
		}

		_, err = client.UpdateUserRole(ctx, &auth.UpdateUserRoleRequest{
			TargetUserId: targetUserID,
			NewRole:      input.NewRole,
		})
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() == codes.PermissionDenied {
				c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
	}
}

func updateProfileHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			DisplayName string `json:"display_name"`
			AvatarUrl   string `json:"avatar_url"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx := middleware.GetGrpcContext(c)
		resp, err := client.UpdateProfile(ctx, &auth.UpdateProfileRequest{
			DisplayName: input.DisplayName,
			AvatarUrl:   input.AvatarUrl,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": status.Convert(err).Message()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func deleteAccountHandler(client auth.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := middleware.GetGrpcContext(c)
		resp, err := client.DeleteAccount(ctx, &auth.DeleteAccountRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": status.Convert(err).Message()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
