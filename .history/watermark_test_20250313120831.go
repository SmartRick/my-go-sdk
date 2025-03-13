package gowatermark

import (
	"image/color"
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

func TestCreateTextWatermark(t *testing.T) {
	config := TextWatermarkConfig{
		OriginImagePath:    "./testdata/origin.jpg",
		CompositeImagePath: "./testdata/composite.jpg",
		FontPath:           "./testdata/font.ttf",
		TextInfos: []TextInfo{
			{
				Size: 100,
				Text: "hello world",
				X:    700,
				Y:    700,
				Color: color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 255,
				},
			},
		},
	}
	err := CreateTextWatermark(config)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTextTiledWatermark(t *testing.T) {
	config := TextTiledWatermarkConfig{
		OriginImagePath:    "./testdata/origin.jpg",
		CompositeImagePath: "./testdata/composite.jpg",
		FontPath:           "./testdata/font.ttf",
		Text:               "hello world",
		TiledRows:          3,
		TiledCols:          4,
		Color: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 240,
		},
	}
	err := CreateTextTiledWatermark(config)
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

	// var err error
	// 测试平铺水印
	// configTiled := TransparentTextWatermarkConfig{
	// 	OriginImagePath:    "./testdata/origin.jpg",
	// 	CompositeImagePath: "./testdata/composite_transparent_text_tiled.jpg",
	// 	FontPath:           "./testdata/font.ttf",
	// 	Text:               "透明文字水印",
	// 	Size:               48,
	// 	Color: color.RGBA{
	// 		R: 255,
	// 		G: 255,
	// 		B: 255,
	// 		A: 255,
	// 	},
	// 	WatermarkPos: Tiled,
	// 	Opacity:      0.6, // 设置70%透明度
	// 	TiledRows:    4,
	// 	TiledCols:    5,
	// 	Rotation:     45, // 45度旋转
	// }
	// err = CreateTransparentTextWatermark(configTiled)
	// if err != nil {
	// 	t.Error(err)
	// }
}
