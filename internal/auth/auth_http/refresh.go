package auth_http

import (
	"context"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/Chamistery/TestTask/internal/auth/utils"
	"github.com/Chamistery/TestTask/internal/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Refresh(ctx context.Context, req *model.RefreshRequest, ip string) (*model.Response, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(i.tokenConfig.GetRefr()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.UserInfo{
		Ip: claims.ip,
	},
		[]byte(i.tokenConfig.GetAccess()),
		i.tokenConfig.GetAccessTime(),
	)
	refreshToken, err := utils.GenerateToken(ip,
		[]byte(i.tokenConfig.GetRefr()),
		i.tokenConfig.GetRefreshTime(),
	)
	currAuth, err := i.authService.Refresh(ctx, req.GetRefreshToken(), refreshToken)
	if err != nil {
		return nil, err
	}

	if ip != claims.ip {
		logger.Warn(fmt.Sprintf("[MOCK EMAIL] Sent warning to user %s: Token was refreshed from new IP %s (was %s)\n",
			currAuth, ip, claims.Ip))
	}

	if err != nil {
		return nil, err
	}
	return &Response{access_token: accessToken, refresh_token: refreshToken}, nil
}
