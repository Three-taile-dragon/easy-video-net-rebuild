package common

import "os"

func GetFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -1
	}
	return fileInfo.Size()
}
