package auth

import (
	"context"
)

func (s *serv) Create(ctx context.Context, create CreateModel) (string, error) {
	auth, err := s.authRepository.Create(ctx, create)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
