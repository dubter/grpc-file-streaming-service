package config

import (
	"context"
	"errors"
	"file_transfer/server/internal/ports/gateway"
	"file_transfer/server/internal/ports/grpc"
	"file_transfer/server/internal/repository/file_system"
	"file_transfer/server/internal/service"
	"gopkg.in/yaml.v2"
	"os"
)

type YamlConfig struct {
	GrpcServer        *YamlGrpcServer        `yaml:"grpc_server"`
	HttpGatewayServer *YamlHttpGatewayServer `yaml:"http_gateway"`
	FolderPath        string                 `yaml:"folder_path"`
}

type Components struct {
	GrpcServer  *grpc.Server
	HttpGateway *gateway.Gateway
}

func ParseConfig(ctx context.Context, configPath string) (*Components, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.New("Failed to read config file: " + err.Error())
	}

	var config YamlConfig
	err = yaml.UnmarshalStrict(data, &config)
	if err != nil {
		return nil, errors.New("Failed to unmarshal config file: " + err.Error())
	}

	return createService(ctx, &config)
}

func createService(ctx context.Context, yaml *YamlConfig) (*Components, error) {
	repo, err := file_system.NewFileSystem(yaml.FolderPath)
	if err != nil {
		return nil, err
	}

	service := service.NewFileTransferService(repo)

	grpcServer, err := grpcServerFromYAML(yaml.GrpcServer, service)
	if err != nil {
		return nil, err
	}

	httpGateway, err := httpGatewayFromYAML(ctx, yaml.HttpGatewayServer, grpcServer.Addr)
	if err != nil {
		return nil, err
	}

	return &Components{
		GrpcServer:  grpcServer,
		HttpGateway: httpGateway,
	}, nil
}
