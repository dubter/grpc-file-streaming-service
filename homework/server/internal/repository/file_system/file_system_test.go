package file_system

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileSystem(t *testing.T) {
	folderPath := "test_folder"
	repo, err := NewFileSystem(folderPath)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if repo == nil {
		t.Fatal("Expected a valid file system, got nil")
	}

	if _, err := os.Stat(folderPath); err != nil {
		t.Fatalf("Expected folder to be created, got error: %v", err)
	}

	defer os.RemoveAll(folderPath)
}

func TestGetAllFilesNames(t *testing.T) {
	folderPath := "test_folder"
	repo, _ := NewFileSystem(folderPath)

	fileNames := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, name := range fileNames {
		file, _ := os.Create(filepath.Join(folderPath, name))
		defer file.Close()
	}

	names, err := repo.GetAllFilesNames(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, expectedName := range fileNames {
		found := false
		for _, actualName := range names {
			if expectedName == actualName {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected file %s in the list, but not found", expectedName)
		}
	}
}

func TestGetFileInfo(t *testing.T) {
	folderPath := "test_folder"
	repo, _ := NewFileSystem(folderPath)

	fileName := "file.txt"
	file, _ := os.Create(filepath.Join(folderPath, fileName))
	defer file.Close()

	fileInfo, err := repo.GetFileInfo(context.Background(), fileName)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fileInfo.Name != fileName {
		t.Fatalf("Expected file name to be %s, got %s", fileName, fileInfo.Name)
	}

	defer os.RemoveAll(folderPath)
}
