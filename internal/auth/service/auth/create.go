package auth

import (
	"context"
	"github.com/Chamistery/TestTask/internal/auth/model"
)

func (s *serv) Create(ctx context.Context, create model.CreateModel) (string, error) {
	auth, err := s.authRepository.Create(ctx, create)
	if err != nil {
		return "", err
	}

	return auth, nil
}
