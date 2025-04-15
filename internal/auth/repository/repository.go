package repository

import (
	"context"
)

type AuthRepository interface {
	Create(ctx context.Context, create CreateModel) (string, error)
	Refresh(ctx context.Context, refr RefreshModel) (string, error)
}
