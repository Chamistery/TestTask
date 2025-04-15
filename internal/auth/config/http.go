package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	httpAuthHostEnvName = "HTTP_AUTH_HOST"
	httpAuthPortEnvName = "HTTP_AUTH_PORT"
)

type HTTPConfig interface {
	Address() string
}

type httpAuthConfig struct {
	host string
	port string
}

func NewHTTPAuthConfig() (HTTPConfig, error) {
	host := os.Getenv(httpAuthHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpAuthPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}
	return &httpAuthConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpAuthConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
