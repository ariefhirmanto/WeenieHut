package imagecompressor

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp" // For decoding
)

var (
	MaxConcurrentCompress, _ = strconv.Atoi(os.Getenv("MAX_CONCURRENT_COMPRESS"))
)

type ImageCompressor struct {
	semaphore chan struct{}
}

func New(
	maxConcurrentCompress int,
) *ImageCompressor {
	return &ImageCompressor{
		semaphore: make(chan struct{}, maxConcurrentCompress),
	}
}

func (cmp *ImageCompressor) compressPNG(ctx context.Context, img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	if err := encoder.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("error encoding image: %w", err)
	}

	return buf.Bytes(), nil
}

func (cmp *ImageCompressor) compressJPEG(ctx context.Context, img image.Image) ([]byte, error) {
	var buf bytes.Buffer

	err := jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: 10,
	})

	if err != nil {
		return nil, fmt.Errorf("error encoding image: %w", err)
	}

	return buf.Bytes(), nil
}

func (cmp *ImageCompressor) thumbnail(ctx context.Context, img image.Image, sizeInPixels int) image.Image {
	thumbnail := resize.Thumbnail(uint(sizeInPixels), uint(sizeInPixels), img, resize.Lanczos2)
	return thumbnail
}

func (cmp *ImageCompressor) Compress(ctx context.Context, src string) (string, error) {
	select {
	case cmp.semaphore <- struct{}{}:
		defer func() { <-cmp.semaphore }()
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("compression queue timeout")
	}

	file, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding file: %w", err)
	}

	thumbnail := cmp.thumbnail(ctx, img, 300)
	var result []byte
	switch format {
	case "jpeg":
		result, err = cmp.compressJPEG(ctx, thumbnail)
	case "png":
		result, err = cmp.compressJPEG(ctx, thumbnail)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}

	if err != nil {
		return "", fmt.Errorf("error in compressing image: %w", err)
	}

	ext := filepath.Ext(src)
	nameWithoutExt := strings.TrimSuffix(src, ext)
	resultFilename := fmt.Sprintf("%s_compressed.%s", nameWithoutExt, "jpeg")

	resultFile, err := os.Create(resultFilename)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer resultFile.Close()

	_, err = resultFile.Write(result)
	if err != nil {
		return "", fmt.Errorf("error writing data: %w", err)
	}
	return resultFilename, nil
}
