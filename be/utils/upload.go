package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFile(file multipart.File, filename string) (string, error) {
	if _, err := os.Stat(UploadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(UploadPath, os.ModePerm); err != nil {
			log.Println("Error creating upload directory:", err)
			return "", err
		}
	}

	filePath := filepath.Join(UploadPath, filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file) // Use io.Copy
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s", filePath), nil
}
