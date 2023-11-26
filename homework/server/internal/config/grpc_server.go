package config

import (
	"file_transfer/server/internal/ports/grpc"
	"file_transfer/server/internal/service"
)

type YamlGrpcServer struct {
	Port      string `yaml:"port"`
	Host      string `yaml:"host"`
	ChunkSize int    `yaml:"chunk_size"`
}

func grpcServerFromYAML(yaml *YamlGrpcServer, app service.Service) (*grpc.Server, error) {
	return grpc.NewServer(&grpc.GrpcConfig{
		Port:      yaml.Port,
		Host:      yaml.Host,
		ChunkSize: yaml.ChunkSize,
		App:       app,
	})
}
