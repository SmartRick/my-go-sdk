package common

import (
	"bytes"
	"image"
	"image/png"
)

// ImageToBytes 将image转换成bytes
func ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// BytesToImage 将bytes转换成image
func BytesToImage(data []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}
