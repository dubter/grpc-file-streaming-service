package repository

import (
	"context"
	"os"
)

//go:generate mockgen -package internal -destination ../mocks/repository.go . Repository
type Repository interface {
	GetFile(context.Context, string) (*os.File, error)
	GetAllFilesNames(context.Context) ([]string, error)
	GetFileInfo(context.Context, string) (*FileInfo, error)
}

type FileInfo struct {
	Name             string
	Size             string
	Type             string
	CreationDate     string
	LastModifiedDate string
	AccessRights     string
	Location         string
}
