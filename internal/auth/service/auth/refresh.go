package auth

import (
	"context"
	"github.com/Chamistery/TestTask/internal/auth/model"
)

func (s *serv) Refresh(ctx context.Context, refr model.RefreshModel) (string, error) {
	auth, err := s.authRepository.Refresh(ctx, refr)
	if err != nil {
		return "", err
	}

	return auth, nil
}
