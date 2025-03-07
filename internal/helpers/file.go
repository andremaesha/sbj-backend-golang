package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveTempFile(fileHeader *multipart.FileHeader, destDir, filePrefix, fileSuffix string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err = src.Close()
		if err != nil {
			panic(err)
		}
	}(src)

	if destDir == "" {
		destDir = os.TempDir()
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(fileHeader.Filename)

	fileName := fmt.Sprintf("%s%d%s%s", filePrefix, time.Now().UnixNano(), fileSuffix, ext)
	tempFilePath := filepath.Join(destDir, fileName)

	dst, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			panic(err)
		}
	}(dst)

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}
