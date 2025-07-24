package util

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrFileTooLarge         = errors.New("file too large, maximum size is 2MB")
	ErrFailedToOpenFile     = errors.New("failed to open file")
	ErrFailedToCheckTypeFile = errors.New("failed to check file type")
	ErrFileTypeNotAllowed   = errors.New("file type not allowed")
)

var allowedTypes = []string{"image/jpeg", "image/png", "image/webp"}

func ValidateFile(file *multipart.FileHeader) error {
	if file.Size > 1024*1024*2 {
		return ErrFileTooLarge
	}

	open, err := file.Open()
	if err != nil {
		return ErrFailedToOpenFile
	}
	defer open.Close()

	ok, fileType, err := isAllowedFileType(open)
	if err != nil {
		return ErrFailedToCheckTypeFile
	}

	if !ok {
		return fmt.Errorf("%w: type file not allowed: %s", ErrFileTypeNotAllowed, fileType)
	}

	return nil
}

func SanitizeFileName(original string) string {
	safeName := strings.ReplaceAll(original, " ", "-")

	ext := filepath.Ext(safeName)
	name := strings.TrimSuffix(safeName, ext)

	rand.Seed(time.Now().UnixNano())
	randomStr := fmt.Sprintf("%06d", rand.Intn(1000000))

	newName := fmt.Sprintf("%s-%s%s", name, randomStr, ext)
	return newName
}

func isAllowedFileType(file multipart.File) (bool, string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, "", err
	}

	mimeType := http.DetectContentType(buffer)

	for _, t := range allowedTypes {
		if mimeType == t {
			return true, mimeType, nil
		}
	}

	return false, mimeType, nil
}
