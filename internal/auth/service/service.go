package service

import (
	"context"
	"github.com/Chamistery/TestTask/internal/auth/model"
)

type AuthService interface {
	Create(ctx context.Context, create model.CreateModel) (string, error)
	Refresh(ctx context.Context, refr model.RefreshModel) (string, error)
}
