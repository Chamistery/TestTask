package service

import (
	"context"
)

type AuthService interface {
	Create(ctx context.Context, create CreateModel) (string, error)
	Refresh(ctx context.Context, refr RefreshModel) (string, error)
}
