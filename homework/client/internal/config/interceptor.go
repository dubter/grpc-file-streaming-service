package config

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor measures the duration of RPC calls
func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
		method,
		time.Since(start),
		err)
	return err
}

// StreamClientInterceptor logs the duration of streaming RPC calls
func StreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	start := time.Now()
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	log.Printf("Stream Request - Method:%s\tDuration:%s\tError:%v\n",
		method,
		time.Since(start),
		err)
	return clientStream, err
}
