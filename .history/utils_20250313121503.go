package gowatermark

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// 常用颜色预定义
var (
	Black = color.RGBA{0, 0, 0, 255}
	White = color.RGBA{255, 255, 255, 255}
	Red   = color.RGBA{255, 0, 0, 255}
	Green = color.RGBA{0, 255, 0, 255}
	Blue  = color.RGBA{0, 0, 255, 255}
)

// PathExists 检查路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// PrepareOutputPath 准备输出路径
// 如果输出文件已存在则删除，如果目录不存在则创建
func PrepareOutputPath(path string) error {
	// 检查文件是否存在
	isExists, _ := PathExists(path)
	if isExists {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	// 检查目录是否存在
	dirPath := filepath.Dir(path)
	isExist, _ := PathExists(dirPath)
	if !isExist {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// MeasureText 计算给定文字的宽度和高度
func MeasureText(face font.Face, text string) (int, int) {
	var (
		width  int
		height int
	)
	// 获取字体的度量信息
	metrics := face.Metrics()

	// 遍历每个字符，计算总宽度
	for _, textRune := range text {
		// 获取字符的水平间距
		advance, _ := face.GlyphAdvance(textRune)
		width += int(advance >> 6)
	}

	// 计算高度
	height = int(metrics.Height >> 6)

	return width, height
}

// CreateTextImage 创建文字图像
func CreateTextImage(text string, fontPath string, fontSize float64, textColor color.RGBA, rotation float64) (*image.NRGBA, error) {
	// 加载字体文件
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	fontFace, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	// 设置字体大小和选项
	face := truetype.NewFace(fontFace, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// 计算文字的宽度和高度
	textWidth, textHeight := MeasureText(face, text)

	// 文字需要旋转时，确保最终图像足够大以容纳旋转后的文本
	padding := 0
	if rotation != 0 {
		// 当旋转时，需要更大的画布以确保文本在旋转后不会被裁剪
		diagonal := int(math.Sqrt(float64(textWidth*textWidth + textHeight*textHeight)))
		padding = (diagonal - textWidth) / 2
	}

	// 创建一个完全透明的新图像
	img := image.NewRGBA(image.Rect(0, 0, textWidth+padding*2, textHeight+padding*2))
	draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)

	// 绘制文字
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	// 设置绘制起点（添加边距）
	d.Dot = fixed.P(padding, textHeight+padding-5)

	// 绘制文字
	d.DrawString(text)

	// 如果需要旋转
	var dst *image.NRGBA
	if rotation != 0 {
		dst = imaging.Rotate(img, rotation, color.Transparent)
	} else {
		dst = imaging.Clone(img)
	}

	return dst, nil
}
