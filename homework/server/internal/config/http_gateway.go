package config

import (
	"context"
	"file_transfer/server/internal/ports/gateway"
	"time"
)

type YamlHttpGatewayServer struct {
	Port            string        `yaml:"port"`
	Host            string        `yaml:"host"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
}

func httpGatewayFromYAML(ctx context.Context, yaml *YamlHttpGatewayServer, grpcAddr string) (*gateway.Gateway, error) {
	return gateway.NewGateway(ctx, &gateway.HttpGatewayConfig{
		Port:            yaml.Port,
		Host:            yaml.Host,
		GrpcAddr:        grpcAddr,
		ShutdownTimeout: yaml.ShutdownTimeout,
		ReadTimeout:     yaml.ReadTimeout,
		WriteTimeout:    yaml.WriteTimeout,
	})
}
