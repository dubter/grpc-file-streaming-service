package file_system

import (
	"fmt"
	"os"
)

func formatFileSize(size int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	case size < GB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size < TB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size < PB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	default:
		return fmt.Sprintf("%.2f PB", float64(size)/float64(PB))
	}
}

func GetFileType(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()

	if mode.IsRegular() {
		return "Regular File"
	} else if mode.IsDir() {
		return "Directory"
	} else if mode&os.ModeSymlink != 0 {
		return "Symbolic Link"
	} else if mode&os.ModeNamedPipe != 0 {
		return "Named Pipe"
	} else if mode&os.ModeSocket != 0 {
		return "Socket"
	} else if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice != 0 {
			return "Character Device"
		} else {
			return "Block Device"
		}
	}

	return "Unknown"
}

// expandTilde change '~/' to home dir
func expandTilde(path string) (string, error) {
	if path[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("can not parse home dir: %v", err)
		}
		path = homeDir + path[1:]
	}
	return path, nil
}
