package files

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

// Get meta data from image format of []byte
func GetMetaFromImageByteSlice(image []byte) (width int, height int, fileExtension string, fileSize int, err error) {
	imageReader := bytes.NewReader(image)

	imgNew, err := imaging.Decode(imageReader)
	if err != nil {
		return 0, 0, "", 0, err

	}

	////////////////////////////////
	// Width and height
	bounds := imgNew.Bounds()
	width = bounds.Max.X - bounds.Min.X
	height = bounds.Max.Y - bounds.Min.Y

	/////////////////
	// Size
	fileSize = len(image)

	/////////////////
	// Extension
	// Detect content type
	contentType := http.DetectContentType(image)
	// Extract file extension from content type
	fileExtension = strings.TrimPrefix(contentType, "image/")

	// Result
	return width, height, fileExtension, fileSize, nil
}
