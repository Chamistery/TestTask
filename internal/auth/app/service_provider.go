package app

import (
	"context"
	"flag"

	"github.com/Chamistery/TestTask/internal/auth/client/db"
	"github.com/Chamistery/TestTask/internal/auth/client/db/pg"
	"github.com/Chamistery/TestTask/internal/auth/client/db/transaction"
	"github.com/Chamistery/TestTask/internal/auth/closer"
	config2 "github.com/Chamistery/TestTask/internal/auth/config"
	"github.com/Chamistery/TestTask/internal/auth/logger"
	"github.com/Chamistery/TestTask/internal/auth/repository"
	authRepository "github.com/Chamistery/TestTask/internal/auth/repository/auth"
	"github.com/Chamistery/TestTask/internal/auth/service"
	authService "github.com/Chamistery/TestTask/internal/auth/service/auth"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var logLevel = flag.String("l", "info", "log level")

type serviceProvider struct {
	pgConfig   config2.PGConfig
	httpConfig config2.HTTPConfig

	dbClient       db.Client
	txManager      db.TxManager
	authRepository repository.AuthRepository

	authService service.AuthService

	authhttpImpl *auth_http.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config2.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config2.NewPGConfig()
		if err != nil {
			logger.Error("failed to get pg config: %s", zap.Error(err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) HTTPAuthConfig() config2.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config2.NewHTTPAuthConfig()
		if err != nil {
			logger.Error("failed to get http config: %s", zap.Error(err))
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Error("failed to create db client: %v", zap.Error(err))
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Error("ping error: %s", zap.Error(err))
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) HTTPAuthImpl(ctx context.Context) *auth_http.Implementation {
	if s.authhttpImpl == nil {
		s.authhttpImpl = auth_http.NewImplementation(s.AuthService(ctx))
	}

	return s.authhttpImpl
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
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
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
