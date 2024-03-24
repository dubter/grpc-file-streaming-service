package grpc

import (
	"context"
	"errors"
	mocks "file_transfer/server/internal/mocks"
	"file_transfer/server/internal/repository"
	"file_transfer/server/internal/service"
	pb "file_transfer/server/pkg/file_transfer_grpc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer_GetFileByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := mocks.NewMockService(ctrl)

	server, err := NewServer(&GrpcConfig{
		App:       mockApp,
		ChunkSize: 1024,
	})

	require.NoError(t, err)

	ctx := context.Background()
	fileName := "test.txt"
	request := &pb.GetFileRequest{Name: fileName}

	mockApp.EXPECT().GetFile(ctx, fileName).Return(nil, errors.New("file not found"))
	err = server.GetFile(request, nil)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestServer_GetFileInfoByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := mocks.NewMockService(ctrl)

	server, err := NewServer(&GrpcConfig{
		App:       mockApp,
		ChunkSize: 1024,
	})

	require.NoError(t, err)

	ctx := context.Background()
	fileName := "test.txt"
	request := &pb.GetFileInfoRequest{Name: fileName}

	mockApp.EXPECT().GetFileInfo(ctx, fileName).Return(nil, errors.New("file not found"))
	_, err = server.GetFileInfo(ctx, request)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	mockFileInfo := &repository.FileInfo{
		Name:         "test.txt",
		Size:         "10KB",
		Type:         "text/plain",
		AccessRights: "rw-r--r--",
		Location:     "/path/to/file",
	}

	mockApp.EXPECT().GetFileInfo(ctx, fileName).Return(mockFileInfo, nil)
	response, err := server.GetFileInfo(ctx, request)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if response.Name != mockFileInfo.Name || response.Size != mockFileInfo.Size {
		t.Errorf("Response does not match expected file info")
	}
}

func TestServer_GetAllFilesNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := mocks.NewMockService(ctrl)

	server, err := NewServer(&GrpcConfig{
		App:       mockApp,
		ChunkSize: 1024,
	})

	require.NoError(t, err)

	ctx := context.Background()
	request := &pb.GetFileListRequest{}

	mockApp.EXPECT().GetAllFilesNames(ctx).Return(nil, service.ErrEmptyRepository)
	_, err = server.GetFileList(ctx, request)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	fileNames := []string{"file1.txt", "file2.txt", "file3.txt"}
	mockApp.EXPECT().GetAllFilesNames(ctx).Return(fileNames, nil)
	response, err := server.GetFileList(ctx, request)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(response.Names) != len(fileNames) {
		t.Errorf("Response does not match expected file names")
	}
}
