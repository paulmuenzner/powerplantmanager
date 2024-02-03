package files

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

// Convert image format of []byte to image.Image
func ByteSliceToGolangImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Convert image format of image.Image to []byte
func GolangImageToByteSlice(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	switch format {
	case "image/png":
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	case "image/jpeg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
	return buf.Bytes(), nil
}
