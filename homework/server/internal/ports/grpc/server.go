package grpc

import (
	"file_transfer/server/internal/service"
	pb "file_transfer/server/pkg/file_transfer_grpc"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	defaultGrpcAddr  = ":50051"
	defaultChunkSize = 64 * 1024 // 64Kb
)

type Server struct {
	pb.FileTransferServiceServer
	app service.Service

	Addr       string
	grpcServer *grpc.Server
	listener   net.Listener
	chunkSize  int
}

type GrpcConfig struct {
	Host      string
	Port      string
	ChunkSize int
	App       service.Service
}

func NewServer(config *GrpcConfig) (*Server, error) {
	addr := fmt.Sprint(config.Host, ":", config.Port)
	if addr == "" {
		addr = defaultGrpcAddr
	}

	chunkSize := config.ChunkSize
	if chunkSize < 1 {
		chunkSize = defaultChunkSize
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryServerInterceptor),
		grpc.ChainStreamInterceptor(StreamServerInterceptor))

	server := &Server{
		Addr:       addr,
		grpcServer: grpcServer,
		listener:   lis,
		app:        config.App,
		chunkSize:  chunkSize,
	}

	pb.RegisterFileTransferServiceServer(grpcServer, server)

	return server, nil
}

func (s *Server) Run() error {
	log.Printf("start grpc server listening on %v", s.listener.Addr())
	defer log.Printf("close grpc server listening on %v", s.listener.Addr())
	defer s.listener.Close()

	if err := s.grpcServer.Serve(s.listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
