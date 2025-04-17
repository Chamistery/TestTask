package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/closer"
	"github.com/Chamistery/TestTask/internal/auth/service/mocks"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Chamistery/TestTask/internal/auth/auth_http"
	"github.com/Chamistery/TestTask/internal/auth/config"
	"github.com/Chamistery/TestTask/internal/auth/logger"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	path, err := config.FindProjectRoot()
	if err != nil {
		panic(err)
	}
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(path, "logs/app_test.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	logLevel := "info"
	if err := level.Set(logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}

func initConfig(_ context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	envFilePath, err := config.GetEnvFilePath()
	if err != nil {
		return fmt.Errorf("ошибка получения пути к .env: %v", err)
	}

	err = config.Load(envFilePath)
	if err != nil {
		return err
	}

	return nil
}

func TestCreate(t *testing.T) {
	initConfig(context.Background())
	logger.Init(getCore(getAtomicLevel()))
	mc := minimock.NewController(t)
	defer mc.Finish()

	mockAuthService := mocks.NewAuthServiceMock(mc)

	tests := []struct {
		name           string
		createReq      model.CreateRequest
		ip             string
		mockCreateFunc func(ctx context.Context, m model.CreateModel) (string, error)
		expectError    bool
	}{
		{
			name:      "Success",
			createReq: model.CreateRequest{Guid: "user-guid-123"},
			ip:        "127.0.0.1",
			mockCreateFunc: func(ctx context.Context, m model.CreateModel) (string, error) {
				if m.Uuid == "" || m.RefreshToken == "" || m.Guid != "user-guid-123" {
					return "", errors.New("invalid model values")
				}
				return "expected-guid", nil
			},
			expectError: false,
		},
		{
			name:      "AuthService Create fails",
			createReq: model.CreateRequest{Guid: "user-guid-456"},
			ip:        "127.0.0.1",
			mockCreateFunc: func(ctx context.Context, m model.CreateModel) (string, error) {
				return "", errors.New("create failed")
			},
			expectError: true,
		},
	}

	impl := auth_http.NewImplementation(mockAuthService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthService.CreateMock.Set(func(ctx context.Context, m model.CreateModel) (string, error) {
				return tc.mockCreateFunc(ctx, m)
			})

			resp, err := impl.Create(context.Background(), &tc.createReq, tc.ip)
			if tc.expectError {
				assert.Error(t, err, "expected error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp, "expected non-nil response")
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
			}
		})
	}
}
