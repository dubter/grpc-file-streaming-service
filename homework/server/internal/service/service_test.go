package service

import (
	"context"
	mocks "file_transfer/server/internal/mocks"
	"file_transfer/server/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	testFilePath := "testfile.txt"
	ctx := context.Background()
	mockRepo.EXPECT().GetFile(ctx, testFilePath).Return(&os.File{}, nil)

	service := NewFileTransferService(mockRepo)

	file, err := service.GetFile(ctx, testFilePath)

	assert.Nil(t, err)
	assert.NotNil(t, file)
}

func TestGetFileInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	testFilePath := "testfile.txt"
	ctx := context.Background()
	mockRepo.EXPECT().GetFileInfo(ctx, testFilePath).Return(&repository.FileInfo{}, nil)

	service := NewFileTransferService(mockRepo)

	fileInfo, err := service.GetFileInfo(ctx, testFilePath)

	assert.Nil(t, err)
	assert.NotNil(t, fileInfo)
}

func TestGetAllFilesNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	testFiles := []string{"file1.txt", "file2.txt"}
	ctx := context.Background()
	mockRepo.EXPECT().GetAllFilesNames(ctx).Return(testFiles, nil)

	service := NewFileTransferService(mockRepo)

	fileNames, err := service.GetAllFilesNames(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, testFiles, fileNames)
}

func TestGetAllFilesNames_EmptyRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)

	var testFiles []string
	ctx := context.Background()
	mockRepo.EXPECT().GetAllFilesNames(ctx).Return(testFiles, nil)

	service := NewFileTransferService(mockRepo)

	fileNames, err := service.GetAllFilesNames(context.Background())

	assert.NotNil(t, err)
	assert.Equal(t, ErrEmptyRepository, err)
	assert.Nil(t, fileNames)
}
