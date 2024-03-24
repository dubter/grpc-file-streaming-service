package grpc

import (
	"context"
	"errors"
	"file_transfer/server/internal/service"
	pb "file_transfer/server/pkg/file_transfer_grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

func (s *Server) GetFile(in *pb.GetFileRequest, stream pb.FileTransferService_GetFileServer) error {
	ctx := context.Background()

	if in.Name == "" {
		return status.Errorf(codes.InvalidArgument, "Empty file name")
	}

	file, err := s.app.GetFile(ctx, in.Name)
	if err != nil {
		return status.Errorf(codes.NotFound, "File not found: %v", err)
	}

	for {
		chunk := make([]byte, s.chunkSize)
		n, err := file.Read(chunk)
		if err == io.EOF {
			break
		}

		if err != nil {
			return status.Errorf(codes.Internal, "Can not read file: %v", err)
		}

		err = stream.Send(&pb.FileChunk{
			ChunkData: chunk[:n],
		})
		if err != nil {
			return status.Errorf(codes.Internal, "Can not send file: %v", err)
		}
	}

	return nil
}

func (s *Server) GetFileInfo(ctx context.Context, in *pb.GetFileInfoRequest) (*pb.GetFileInfoResponse, error) {
	if in.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty file name")
	}

	info, err := s.app.GetFileInfo(ctx, in.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "File not found: %v", err)
	}

	return &pb.GetFileInfoResponse{
		Name:             info.Name,
		Size:             info.Size,
		Type:             info.Type,
		CreationDate:     info.CreationDate,
		LastModifiedDate: info.LastModifiedDate,
		AccessRights:     info.AccessRights,
		Location:         info.Location,
	}, nil
}

func (s *Server) GetFileList(ctx context.Context, _ *pb.GetFileListRequest) (*pb.GetFileListResponse, error) {
	names, err := s.app.GetAllFilesNames(ctx)
	if err != nil {
		if errors.Is(err, service.ErrEmptyRepository) {
			return nil, status.Errorf(codes.NotFound, "Files not found: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return &pb.GetFileListResponse{
		Names: names,
	}, nil
}
