package service

import (
	"context"
	"errors"
	"file_transfer/server/internal/repository"
	"os"
)

var (
	ErrEmptyRepository = errors.New("have not any files")
)

//go:generate mockgen -package internal -destination ../mocks/service.go . Service
type Service interface {
	GetFile(context.Context, string) (*os.File, error)
	GetAllFilesNames(context.Context) ([]string, error)
	GetFileInfo(context.Context, string) (*repository.FileInfo, error)
}

type fileTransferService struct {
	repo repository.Repository
}

func NewFileTransferService(repo repository.Repository) Service {
	return &fileTransferService{
		repo: repo,
	}
}

func (s *fileTransferService) GetFile(ctx context.Context, name string) (*os.File, error) {
	return s.repo.GetFile(ctx, name)
}

func (s *fileTransferService) GetFileInfo(ctx context.Context, name string) (*repository.FileInfo, error) {
	return s.repo.GetFileInfo(ctx, name)
}

func (s *fileTransferService) GetAllFilesNames(ctx context.Context) ([]string, error) {
	names, err := s.repo.GetAllFilesNames(ctx)
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return nil, ErrEmptyRepository
	}

	return names, nil
}
