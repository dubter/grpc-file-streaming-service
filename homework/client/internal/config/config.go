package config

import (
	"errors"
	"file_transfer/server/pkg/file_transfer_grpc"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

const (
	defaultServerAddr = ":50051"
	defaultTimeout    = time.Second * 5
)

type YamlConfig struct {
	GetFile     *MethodParams `yaml:"get_file"`
	GetFileList *MethodParams `yaml:"get_file_list"`
	GetFileInfo *MethodParams `yaml:"get_file_info"`
	Timeout     time.Duration `yaml:"timeout"`
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
}

type MethodParams struct {
	Exec bool   `yaml:"exec"`
	Name string `yaml:"name"`
}

type Components struct {
	HasGetFile     *MethodParams
	HasGetFileList *MethodParams
	HasGetFileInfo *MethodParams
	Timeout        time.Duration
	Client         file_transfer_grpc.FileTransferServiceClient
}

func ParseConfig(configPath string) (*Components, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.New("Failed to read config file: " + err.Error())
	}

	var config YamlConfig
	err = yaml.UnmarshalStrict(data, &config)
	if err != nil {
		return nil, errors.New("Failed to unmarshal config file: " + err.Error())
	}

	addr := fmt.Sprint(config.Host, ":", config.Port)
	if addr == "" {
		addr = defaultServerAddr
	}

	timeout := config.Timeout
	if timeout <= time.Duration(0) {
		timeout = defaultTimeout
	}

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithChainStreamInterceptor(StreamClientInterceptor))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	client := file_transfer_grpc.NewFileTransferServiceClient(conn)

	return &Components{
		Client:         client,
		Timeout:        timeout,
		HasGetFile:     config.GetFile,
		HasGetFileList: config.GetFileList,
		HasGetFileInfo: config.GetFileInfo,
	}, nil
}
