package gateway

import (
	"context"
	pb "file_transfer/server/pkg/file_transfer_grpc"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

const (
	defaultShutdownTimeout = 30 * time.Second
	defaultReadTimeout     = 30 * time.Second
	defaultWriteTimeout    = 30 * time.Second
	defaultGatewayAddr     = ":8080"
	defaultGprcAddr        = ":50051"
)

type Gateway struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

type HttpGatewayConfig struct {
	Port            string
	Host            string
	GrpcAddr        string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func NewGateway(ctx context.Context, config *HttpGatewayConfig) (*Gateway, error) {
	readTimeout := config.ReadTimeout
	if readTimeout <= time.Duration(0) {
		readTimeout = defaultReadTimeout
	}

	writeTimeout := config.WriteTimeout
	if writeTimeout < time.Duration(0) {
		writeTimeout = defaultWriteTimeout
	}

	shutdownTimeout := config.ShutdownTimeout
	if shutdownTimeout <= time.Duration(0) {
		shutdownTimeout = defaultShutdownTimeout
	}

	addr := fmt.Sprint(config.Host, ":", config.Port)
	if addr == "" {
		addr = defaultGatewayAddr
	}

	grpcAddr := config.GrpcAddr
	if grpcAddr == "" {
		grpcAddr = defaultGprcAddr
	}

	mux := runtime.NewServeMux()

	conn, err := grpc.DialContext(
		ctx,
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %v", err)
	}

	err = pb.RegisterFileTransferServiceHandler(ctx, mux, conn)
	if err != nil {
		return nil, fmt.Errorf("can not register gateway: %v", err)
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return &Gateway{
		server:          server,
		shutdownTimeout: shutdownTimeout,
	}, nil
}

func (gw *Gateway) Run() error {
	log.Printf("start gateway http server listening on %s", gw.server.Addr)
	defer log.Printf("close gateway http server listening on %s", gw.server.Addr)

	if err := gw.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (gw *Gateway) Shutdown() {
	shCtx, cancel := context.WithTimeout(context.Background(), gw.shutdownTimeout)
	defer cancel()

	if err := gw.server.Shutdown(shCtx); err != nil {
		log.Printf("can't close http gateway listening on %s: %s", gw.server.Addr, err.Error())
	}
}
