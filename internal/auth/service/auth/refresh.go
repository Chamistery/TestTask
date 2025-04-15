package auth

import (
	"context"
)

func (s *serv) Refresh(ctx context.Context, refr RefreshModel) (string, error) {
	auth, err := s.authRepository.Refresh(ctx, refr)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
