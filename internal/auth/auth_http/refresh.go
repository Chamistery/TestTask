package auth_http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/logger"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/Chamistery/TestTask/internal/auth/utils"
)

func (i *Implementation) Refresh(ctx context.Context, req *model.RefreshRequest, ip string) (*model.Response, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(i.tokenConfig.GetRefr()))
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(claims.Uuid, ip,
		[]byte(i.tokenConfig.GetAccess()),
		i.tokenConfig.GetAccessTime(),
	)
	refreshToken, err := utils.GenerateToken(claims.Uuid, ip,
		[]byte(i.tokenConfig.GetRefr()),
		i.tokenConfig.GetRefreshTime(),
	)
	token_old := req.GetRefreshToken()
	token_new, _ := utils.HashToken(refreshToken)
	currAuth, err := i.authService.Refresh(ctx, model.RefreshModel{claims.Uuid, token_old, token_new})
	if err != nil {
		return nil, err
	}

	if ip != claims.Ip {
		logger.Warn(fmt.Sprintf("[MOCK EMAIL] Sent warning to user %s: Token was refreshed from new IP %s (was %s)\n",
			currAuth, ip, claims.Ip))
		return nil, errors.New("Ip changed")
	}

	return &model.Response{accessToken, refreshToken}, nil
}
