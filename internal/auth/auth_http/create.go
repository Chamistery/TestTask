package auth_http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/Chamistery/TestTask/internal/auth/utils"
	"github.com/Chamistery/TestTask/internal/logger"
)

func (i *Implementation) Create(ctx context.Context, req *model.CreateRequest, ip string) (*model.Response, error) {
	accessToken, err := utils.GenerateToken(ip,
		[]byte(i.tokenConfig.GetAccess()),
		i.tokenConfig.GetAccessTime(),
	)
	refreshToken, err := utils.GenerateToken(ip,
		[]byte(i.tokenConfig.GetRefr()),
		i.tokenConfig.GetRefreshTime(),
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	create, err := i.authService.Create(ctx, refreshToken)
	if err != nil {
		nil, err
	}
	logger.Info(fmt.Sprintf("Access and refresh %s tokens created", create))
	return &Response{access_token: accessToken, refresh_token: refreshToken}, nil
}
