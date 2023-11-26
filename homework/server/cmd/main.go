package main

import (
	"context"
	"errors"
	"file_transfer/server/internal/config"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config_path", "", "path to yaml config for server settings")

	flag.Parse()

	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cx := context.Background()

	components, err := config.ParseConfig(cx, configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	grpcServer := components.GrpcServer
	httpGateway := components.HttpGateway

	eg, ctx := errgroup.WithContext(cx)
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(func() error {
		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			close(errCh)
		}()

		go func() {
			if err = grpcServer.Run(); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	// run http gateway
	eg.Go(func() error {
		errCh := make(chan error)

		defer func() {
			httpGateway.Shutdown()
			close(errCh)
		}()

		go func() {
			if err = httpGateway.Run(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
