package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor measures the duration of RPC calls
func UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)
	return h, err
}

// StreamServerInterceptor logs the duration of streaming RPC calls
func StreamServerInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	err := handler(srv, ss)
	log.Printf("Stream Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)
	return err
}
