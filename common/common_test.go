package common

import (
	"image"
	"testing"
)

func TestImageToBytes(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	bytes, err := ImageToBytes(img)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(bytes))
}
