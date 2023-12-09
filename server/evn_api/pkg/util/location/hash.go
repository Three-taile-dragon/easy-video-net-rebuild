package location

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// GetFileHash 计算文件哈希
func GetFileHash(fileDir string) (string, error) {
	file, err := os.Open(fileDir)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}
