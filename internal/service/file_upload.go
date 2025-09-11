package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/storage"
	"WeenieHut/internal/utils"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (s *Service) UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (string, error) {
	if sizeInBytes > constants.MaxUploadSizeInBytes {
		return "", constants.ErrMaximumFileSize
	}
	if err := utils.ValidateFileExtensions(filename, constants.AllowedExtensions); err != nil {
		return "", constants.ErrInvalidFileType
	}

	bucket := storage.S3Bucket
	identifier := uuid.NewString()

	filepath := filepath.Join("/tmp", fmt.Sprintf("%s_%s", identifier, filename))
	tempFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	n, err := io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}

	log.Printf("written size: %d filename: %s", n, tempFile.Name())

	// Upload to object storage
	remotePath := fmt.Sprintf("%s_%s", identifier, filename)
	uri, err := s.storage.UploadFile(ctx, bucket, tempFile.Name(), remotePath)
	if err != nil {
		return "", err
	}

	return uri, nil
}
