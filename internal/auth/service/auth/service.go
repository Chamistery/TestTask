package auth

import (
	"github.com/Chamistery/TestTask/internal/auth/client/db"
	"github.com/Chamistery/TestTask/internal/auth/repository"
	"github.com/Chamistery/TestTask/internal/auth/service"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

func NewService(
	authRepository repository.AuthRepository,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
	}
}
