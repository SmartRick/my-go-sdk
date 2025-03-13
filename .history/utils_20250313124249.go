package gowatermark

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	// DefaultFont 默认字体，程序启动时自动加载系统默认字体
	DefaultFont *truetype.Font
)

// 初始化默认字体
func init() {
	fontPaths := []string{
		"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf", // Linux
		"C:/Windows/Fonts/simhei.ttf",                               // Windows
		"/System/Library/Fonts/PingFang.ttc",                        // macOS
	}

	for _, fontPath := range fontPaths {
		if _, err := os.Stat(fontPath); err == nil {
			fontData, err := os.ReadFile(fontPath)
			if err != nil {
				continue
			}
			DefaultFont, err = freetype.ParseFont(fontData)
			if err == nil {
				break
			}
		}
	}
}

// LoadFont 加载字体文件，如果fontPath为空则使用默认字体
func LoadFont(fontPath string) (*truetype.Font, error) {
	if fontPath == "" {
		if DefaultFont == nil {
			return nil, fmt.Errorf("no default font available and no font path provided")
		}
		return DefaultFont, nil
	}

	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	font, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return font, nil
}

// PathExists 检查文件路径是否存在
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
func CreateTextImage(text string, fontSize float64, fontData []byte, textColor color.RGBA, rotation float64) (*image.NRGBA, error) {
	// 解析字体
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

// PrepareOutputPath 准备输出文件路径，删除已存在的文件并创建必要的目录
func PrepareOutputPath(path string) error {
	// 如果合成图片存在则删除重新生成
	isExists, _ := PathExists(path)
	if isExists {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	// 判断文件夹是否存在，不存在创建
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
