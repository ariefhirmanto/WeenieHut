package imagecompressor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	testdataDir := "testdata"
	compressor := New(5)

	t.Run("PNG", func(t *testing.T) {
		input := fmt.Sprintf("%s/sample.png", testdataDir)
		result, err := compressor.Compress(context.TODO(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("JPEG", func(t *testing.T) {
		input := fmt.Sprintf("%s/sample.jpeg", testdataDir)
		result, err := compressor.Compress(context.TODO(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	// t.Run("JPG", func(t *testing.T) {
	// 	input := fmt.Sprintf("%s/sample.jpg", testdataDir)
	// 	result, err := compressor.Compress(context.TODO(), input)
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, result)
	// })
}
