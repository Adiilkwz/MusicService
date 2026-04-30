package grpc

import (
	"auth_service/internal/domain"

	"github.com/Adiilkwz/music-grpc-go/auth"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer

	accessUC  domain.AccessUsecase
	profileUC domain.ProfileUsecase
	adminUC   domain.AdminUsecase
}

func NewAuthServer(accessUC domain.AccessUsecase, profileUC domain.ProfileUsecase, adminUC domain.AdminUsecase) *AuthServer {
	return &AuthServer{
		accessUC:  accessUC,
		profileUC: profileUC,
		adminUC:   adminUC,
	}
}
