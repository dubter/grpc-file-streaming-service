package main

import (
	"context"
	"file_transfer/client/internal/config"
	"file_transfer/server/pkg/file_transfer_grpc"
	"flag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config_path", "", "path to yaml config for client settings")

	flag.Parse()

	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	components, err := config.ParseConfig(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	client := components.Client

	if components.HasGetFileList.Exec {
		ctx, cancel := context.WithTimeout(context.Background(), components.Timeout)
		defer cancel()

		r, err := client.GetFileList(ctx, &file_transfer_grpc.GetFileListRequest{})
		if err != nil {
			processError(err)
		} else {
			log.Printf("File list: %v", r.GetNames())
		}
	}

	if components.HasGetFile.Exec {
		ctx, cancel := context.WithTimeout(context.Background(), components.Timeout)
		defer cancel()

		r, err := client.GetFile(ctx, &file_transfer_grpc.GetFileRequest{
			Name: components.HasGetFile.Name,
		})
		if err != nil {
			processError(err)
		} else {
			for {
				msg, err := r.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("Error while receiving message: %v", err)
				}

				log.Printf("Received bytes: %s", msg.GetChunkData())
			}
		}
	}

	if components.HasGetFileInfo.Exec {
		ctx, cancel := context.WithTimeout(context.Background(), components.Timeout)
		defer cancel()

		r, err := client.GetFileInfo(ctx, &file_transfer_grpc.GetFileInfoRequest{
			Name: components.HasGetFileInfo.Name,
		})
		if err != nil {
			processError(err)
		} else {
			log.Printf("File name: %v", r.Name)
			log.Printf("File size: %v", r.Size)
			log.Printf("File type: %v", r.Type)
			log.Printf("File location: %v", r.Location)
			log.Printf("File access rights: %v", r.AccessRights)
			log.Printf("File creation date: %v", r.CreationDate)
			log.Printf("File last modified date: %v", r.LastModifiedDate)
		}
	}
}

func processError(err error) {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Internal:
			log.Printf("Internal server error: %v", e.Message())

		case codes.NotFound:
			log.Printf("Not found error: %v", e.Message())

		case codes.InvalidArgument:
			log.Printf("Not found error: %v", e.Message())

		default:
			log.Printf("Unknown error: %v", e.Message())
		}
	} else {
		log.Fatalf("Failed to get file list: %v", err)
	}
}
