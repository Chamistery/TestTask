package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	refreshKey = "REFRESH_TOKEN_SEC"
	accessKey  = "ACCESS_TOKEN_SEC"
	refreshDur = "REFRESH_TOKEN_DUR"
	accessDur  = "ACCESS_TOKEN_DUR"
)

type TokenConfig interface {
	GetRefr() string
	GetAccess() string
	GetRefreshTime() time.Duration
	GetAccessTime() time.Duration
}

type tokenConfig struct {
	refresh     string
	access      string
	refreshTime time.Duration
	accessTime  time.Duration
}

func NewTokenConfig() (TokenConfig, error) {
	CreateSecrets()
	refreshK := os.Getenv(refreshKey)
	accessK := os.Getenv(accessKey)
	accT, err1 := strconv.Atoi(os.Getenv(accessDur))
	refrT, err2 := strconv.Atoi(os.Getenv(refreshDur))
	if err1 != nil || err2 != nil {
		return nil, errors.New("refresh time or access time is invalid")
	}
	accessT := time.Minute * time.Duration(accT)
	refreshT := time.Minute * time.Duration(refrT)
	if len(refreshK) == 0 {
		return nil, errors.New("RefreshKey not found")
	}

	if len(accessK) == 0 {
		return nil, errors.New("AccessKey not found")
	}

	return &tokenConfig{
		refresh:     refreshK,
		access:      accessK,
		refreshTime: refreshT,
		accessTime:  accessT,
	}, nil
}

func (cfg *tokenConfig) GetRefr() string {
	return cfg.refresh
}

func (cfg *tokenConfig) GetAccess() string {
	return cfg.access
}

func (cfg *tokenConfig) GetRefreshTime() time.Duration {
	return cfg.refreshTime
}

func (cfg *tokenConfig) GetAccessTime() time.Duration {
	return cfg.accessTime
}
