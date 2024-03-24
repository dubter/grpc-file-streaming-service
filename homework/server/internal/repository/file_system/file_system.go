package file_system

import (
	"context"
	"file_transfer/server/internal/repository"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type fileSystem struct {
	folderPath string
}

func NewFileSystem(folderPath string) (repository.Repository, error) {
	var err error

	folderPath, err = expandTilde(folderPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(folderPath); err != nil {
		err = os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can not create folder: %v", err)
		}
	}

	return &fileSystem{
		folderPath: folderPath,
	}, nil
}

func (fs *fileSystem) GetAllFilesNames(_ context.Context) ([]string, error) {
	var fileNames []string

	err := filepath.Walk(fs.folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileNames = append(fileNames, info.Name())
		}

		return nil
	})

	return fileNames, err
}

func (fs *fileSystem) GetFile(_ context.Context, name string) (*os.File, error) {
	filePath := filepath.Join(fs.folderPath, name)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can not open file: %v", err)
	}

	return file, nil
}

func (fs *fileSystem) GetFileInfo(_ context.Context, name string) (*repository.FileInfo, error) {
	filePath := filepath.Join(fs.folderPath, name)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("can not see stat: %v", err)
	}

	return &repository.FileInfo{
		Name:             fileInfo.Name(),
		Size:             formatFileSize(fileInfo.Size()),
		Type:             GetFileType(fileInfo),
		CreationDate:     fileInfo.ModTime().Format(time.RFC3339),
		LastModifiedDate: fileInfo.ModTime().Format(time.RFC3339),
		AccessRights:     fileInfo.Mode().String(),
		Location:         filePath,
	}, nil
}
