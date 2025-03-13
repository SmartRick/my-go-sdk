package gowatermark

import (
	"testing"
)

func TestCreateImageWatermark(t *testing.T) {
	config := ImageWatermarkConfig{
		OriginImagePath:    "./testdata/origin.jpg",
		WatermarkImagePath: "./testdata/watermark.png",
		WatermarkPos:       LeftTop,
		CompositeImagePath: "./testdata/composite.jpg",
		Opacity:            0.2,
	}
	err := CreateImageWatermark(config)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTransparentTextWatermark(t *testing.T) {
	// 测试固定位置的文字水印
	config := TransparentTextWatermarkConfig{
		OriginImagePath:    "./testdata/origin.jpg",
		CompositeImagePath: "./testdata/composite_transparent_text.jpg",
		FontPath:           "./testdata/font.ttf",
		Text:               "透明文字水印测试",
		Size:               72,
		Color:              White,
		WatermarkPos:       LeftTop,
		Opacity:            0.5, // 设置50%透明度
		OffsetX:            20,
		OffsetY:            20,
		Rotation:           0, // 不旋转
	}
	err := CreateTransparentTextWatermark(config)
	if err != nil {
		t.Error(err)
	}

	// 测试平铺水印
	configTiled := TransparentTextWatermarkConfig{
		OriginImagePath:    "./testdata/origin.jpg",
		CompositeImagePath: "./testdata/composite_transparent_text_tiled.jpg",
		FontPath:           "./testdata/font.ttf",
		Text:               "透明文字水印",
		Size:               48,
		Color:              White,
		WatermarkPos:       Tiled,
		Opacity:            0.6, // 设置60%透明度
		TiledRows:          4,
		TiledCols:          5,
		Rotation:           -45, // -45度旋转
	}
	err = CreateTransparentTextWatermark(configTiled)
	if err != nil {
		t.Error(err)
	}
}
