package auth_http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/logger"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/Chamistery/TestTask/internal/auth/utils"
	"github.com/google/uuid"
)

func (i *Implementation) Create(ctx context.Context, req *model.CreateRequest, ip string) (*model.Response, error) {
	uuid := uuid.New().String()
	accessToken, err := utils.GenerateToken(uuid, ip,
		[]byte(i.tokenConfig.GetAccess()),
		i.tokenConfig.GetAccessTime(),
	)
	refreshToken, err := utils.GenerateToken(uuid, ip,
		[]byte(i.tokenConfig.GetRefr()),
		i.tokenConfig.GetRefreshTime(),
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	token, _ := utils.HashToken(refreshToken)
	create, err := i.authService.Create(ctx, model.CreateModel{uuid, token, req.Guid})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("Access and refresh (%s) tokens created", create))
	return &model.Response{accessToken, refreshToken}, nil
}
