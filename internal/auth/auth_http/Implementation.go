package auth_grpc

import (
	config2 "github.com/Chamistery/TestTask/internal/auth/config"
	"github.com/Chamistery/TestTask/internal/auth/logger"
	"github.com/Chamistery/TestTask/internal/auth/service"
	"go.uber.org/zap"
)

type Implementation struct {
	authService service.AuthService

	tokenConfig config2.TokenConfig
}

func NewImplementation(authService service.AuthService) *Implementation {
	token, err := config2.NewTokenConfig()
	if err != nil {
		logger.Error("No config for auth and access", zap.Error(err))
	}
	return &Implementation{
		authService: authService,
		tokenConfig: token,
	}
}
