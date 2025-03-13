package gowatermark

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ImageWatermarkConfig struct {
	OriginImagePath    string       // 原图地址
	WatermarkImagePath string       // 水印图地址
	WatermarkPos       watermarkPos // 水印位置
	CompositeImagePath string       // 合成图地址
	OffsetX            int          // 水印位置偏移量X
	OffsetY            int          // 水印位置偏移量Y
	Opacity            float64      // 水印透明度
	TiledRows          int          // 水印图横向平铺行数
	TiledCols          int          // 水印图横向平铺列数
}

type watermarkPos string

const (
	LeftTop     watermarkPos = "left_top"
	RightTop    watermarkPos = "right_top"
	LeftBottom  watermarkPos = "left_bottom"
	RightBottom watermarkPos = "right_bottom"
	Tiled       watermarkPos = "tiled"
)

func CreateImageWatermark(config ImageWatermarkConfig) error {
	watermarkFile, err := os.Open(config.WatermarkImagePath)
	if err != nil {
		return errors.New("open watermark image file error:" + err.Error())
	}
	defer watermarkFile.Close()

	originFile, err := os.Open(config.OriginImagePath)
	if err != nil {
		return errors.New("open origin image file error:" + err.Error())
	}
	defer originFile.Close()

	// 准备输出路径
	if err := PrepareOutputPath(config.CompositeImagePath); err != nil {
		return errors.New("prepare output path error:" + err.Error())
	}

	// 水印透明度判断
	if config.Opacity < 0 || config.Opacity > 1 {
		return errors.New("watermark opacity error:Ensure 0.0 <= opacity <= 1.0")
	}
	if config.Opacity == 0 {
		config.Opacity = 1
	}
	// 获取原图大小
	originImg, _ := imaging.Decode(originFile)
	watermarkImg, _ := imaging.Decode(watermarkFile)
	originImgWidth := originImg.Bounds().Dx()
	originImgHeight := originImg.Bounds().Dy()
	// 对水印图进行缩放(对比原图)
	targetWatermarkImgWidth := uint(originImgWidth / 5)
	destwatermarkImg := imaging.Resize(watermarkImg, int(targetWatermarkImgWidth), 0, imaging.Lanczos)

	// 根据水印位置合成图片
	var destImg image.Image
	switch config.WatermarkPos {
	case LeftTop:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(config.OffsetX, config.OffsetY), config.Opacity)
	case RightTop:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(originImgWidth-int(targetWatermarkImgWidth)-config.OffsetX, config.OffsetY), config.Opacity)
	case LeftBottom:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(config.OffsetX, originImgHeight-destwatermarkImg.Bounds().Dy()-config.OffsetY), config.Opacity)
	case RightBottom:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(originImgWidth-int(targetWatermarkImgWidth)-config.OffsetX, originImgHeight-destwatermarkImg.Bounds().Dy()-config.OffsetY), config.Opacity)
	case Tiled:
		if config.TiledCols == 0 || config.TiledRows == 0 {
			return errors.New("watermark position tiled need tiled_cols and tiled_rows")
		}
		mainBounds := originImg.Bounds()
		watermarkBounds := destwatermarkImg.Bounds()

		// 创建一个与主图相同尺寸的新图像作为结果图像
		result := image.NewNRGBA(mainBounds)
		draw.Draw(result, mainBounds, originImg, image.Point{}, draw.Src)

		// 计算水印在主图上平铺所需的行数和列数
		rows := config.TiledRows
		cols := config.TiledCols
		// 计算行间距和列间距
		totalWidth := cols * watermarkBounds.Dx()
		totalHeight := rows * watermarkBounds.Dy()
		extraWidth := mainBounds.Dx() - totalWidth
		extraHeight := mainBounds.Dy() - totalHeight
		rowSpacing := extraHeight / (rows + 1)
		colSpacing := extraWidth / (cols + 1)
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				// 计算当前水印在主图上的位置
				x := c*(watermarkBounds.Dx()+colSpacing) + colSpacing/2
				y := r*(watermarkBounds.Dy()+rowSpacing) + rowSpacing/2

				// 将水印粘贴到结果图像的相应位置
				draw.DrawMask(result, image.Rect(x, y, x+watermarkBounds.Dx(), y+watermarkBounds.Dy()), destwatermarkImg, destwatermarkImg.Bounds().Min, destwatermarkImg, destwatermarkImg.Bounds().Min, draw.Over)
			}
		}
		destImg = result
	default:
		return errors.New("watermark position error")
	}
	if err = imaging.Save(destImg, config.CompositeImagePath); err != nil {
		return errors.New("create composite image error:" + err.Error())
	}
	return nil
}

type TextWatermarkConfig struct {
	OriginImagePath    string // 原图地址
	CompositeImagePath string // 合成图地址
	FontPath           string // 字体文件地址
	TextInfos          []TextInfo
}

type TextInfo struct {
	Text  string     // 文字内容
	Size  float64    // 文字大小
	Color color.RGBA // 文字颜色透明度
	X     int        // 位置信息
	Y     int        // 位置信息
}

func CreateTextWatermark(config TextWatermarkConfig) error {
	originFile, err := os.Open(config.OriginImagePath)
	if err != nil {
		return errors.New("open origin image file error:" + err.Error())
	}
	defer originFile.Close()
	// 如果合成图片存在则删除重新生成
	isExists, _ := pathExists(config.CompositeImagePath)
	if isExists {
		err = os.Remove(config.CompositeImagePath)
		if err != nil {
			return errors.New("old composite image remove error:" + err.Error())
		}
	}
	// 判断文件夹是否存在，不存在创建
	dirPath := filepath.Dir(config.CompositeImagePath)
	isExist, _ := pathExists(dirPath)
	if !isExist {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	img, err := imaging.Decode(originFile)
	if err != nil {
		return err
	}
	dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Over)
	// load font file
	fontBytes, err := ioutil.ReadFile(config.FontPath)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}
	for _, v := range config.TextInfos {
		f := freetype.NewContext()
		f.SetDPI(72)
		f.SetFont(font)       // 加载字体
		f.SetFontSize(v.Size) // 设置字体尺寸
		f.SetClip(dst.Bounds())
		f.SetDst(dst)
		f.SetSrc(image.NewUniform(v.Color)) // 设置字体颜色
		// 位置信息
		pt := freetype.Pt(v.X, v.Y)
		_, err = f.DrawString(v.Text, pt)
		if err != nil {
			return err
		}
	}
	if err = imaging.Save(dst, config.CompositeImagePath); err != nil {
		return errors.New("create composite image error:" + err.Error())
	}
	return nil
}

type TextTiledWatermarkConfig struct {
	OriginImagePath    string     // 原图地址
	CompositeImagePath string     // 合成图地址
	FontPath           string     // 字体文件地址
	Text               string     // 文字内容
	Color              color.RGBA // 文字颜色透明度
	TiledRows          int        // 水印图横向平铺行数
	TiledCols          int        // 水印图横向平铺列数
}

// TransparentTextWatermarkConfig 透明文字水印配置
type TransparentTextWatermarkConfig struct {
	OriginImagePath    string       // 原图地址
	CompositeImagePath string       // 合成图地址
	FontPath           string       // 字体文件地址
	Text               string       // 文字内容
	Size               float64      // 文字大小
	Color              color.RGBA   // 文字颜色
	WatermarkPos       watermarkPos // 水印位置
	Opacity            float64      // 水印透明度
	OffsetX            int          // 水印位置偏移量X
	OffsetY            int          // 水印位置偏移量Y
	Rotation           float64      // 文字旋转角度
	TiledRows          int          // 水印图横向平铺行数(仅Tiled位置时使用)
	TiledCols          int          // 水印图横向平铺列数(仅Tiled位置时使用)
}

// CreateTransparentTextWatermark 创建透明文字水印
// 先将文字渲染到透明图层，然后作为图片叠加到目标图片上
func CreateTransparentTextWatermark(config TransparentTextWatermarkConfig) error {
	// 输入参数验证
	if config.Opacity < 0 || config.Opacity > 1 {
		return errors.New("watermark opacity error: Ensure 0.0 <= opacity <= 1.0")
	}
	if config.Opacity == 0 {
		config.Opacity = 1
	}
	if config.WatermarkPos == Tiled && (config.TiledCols == 0 || config.TiledRows == 0) {
		return errors.New("watermark position tiled need tiled_cols and tiled_rows")
	}

	// 打开原始图片
	originFile, err := os.Open(config.OriginImagePath)
	if err != nil {
		return errors.New("open origin image file error:" + err.Error())
	}
	defer originFile.Close()

	// 准备输出路径
	if err := PrepareOutputPath(config.CompositeImagePath); err != nil {
		return errors.New("prepare output path error:" + err.Error())
	}

	// 解码原始图片
	originImg, err := imaging.Decode(originFile)
	if err != nil {
		return errors.New("decode origin image error: " + err.Error())
	}

	// 创建文字水印图像
	textWatermarkImg, err := CreateTextImage(config.Text, config.FontPath, config.Size, config.Color, config.Rotation)
	if err != nil {
		return errors.New("create text image error: " + err.Error())
	}

	// 根据水印位置合成图片
	var destImg image.Image
	originImgWidth := originImg.Bounds().Dx()
	originImgHeight := originImg.Bounds().Dy()
	textImgWidth := textWatermarkImg.Bounds().Dx()
	textImgHeight := textWatermarkImg.Bounds().Dy()

	switch config.WatermarkPos {
	case LeftTop:
		destImg = imaging.Overlay(originImg, textWatermarkImg, image.Pt(config.OffsetX, config.OffsetY), config.Opacity)
	case RightTop:
		destImg = imaging.Overlay(originImg, textWatermarkImg, image.Pt(originImgWidth-textImgWidth-config.OffsetX, config.OffsetY), config.Opacity)
	case LeftBottom:
		destImg = imaging.Overlay(originImg, textWatermarkImg, image.Pt(config.OffsetX, originImgHeight-textImgHeight-config.OffsetY), config.Opacity)
	case RightBottom:
		destImg = imaging.Overlay(originImg, textWatermarkImg, image.Pt(originImgWidth-textImgWidth-config.OffsetX, originImgHeight-textImgHeight-config.OffsetY), config.Opacity)
	case Tiled:
		mainBounds := originImg.Bounds()
		watermarkBounds := textWatermarkImg.Bounds()

		// 创建一个与主图相同尺寸的新图像作为结果图像
		result := image.NewNRGBA(mainBounds)
		draw.Draw(result, mainBounds, originImg, image.Point{}, draw.Src)

		// 计算水印在主图上平铺所需的行数和列数
		rows := config.TiledRows
		cols := config.TiledCols

		// 计算行间距和列间距
		totalWidth := cols * watermarkBounds.Dx()
		totalHeight := rows * watermarkBounds.Dy()
		extraWidth := mainBounds.Dx() - totalWidth
		extraHeight := mainBounds.Dy() - totalHeight
		rowSpacing := extraHeight / (rows + 1)
		colSpacing := extraWidth / (cols + 1)

		// 创建一个临时的画布用于叠加水印图像
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				// 计算当前水印在主图上的位置
				x := c*(watermarkBounds.Dx()+colSpacing) + colSpacing
				y := r*(watermarkBounds.Dy()+rowSpacing) + rowSpacing

				// 透明度混合处理
				opacity := config.Opacity
				if opacity == 1.0 {
					draw.Draw(result, image.Rect(x, y, x+watermarkBounds.Dx(), y+watermarkBounds.Dy()),
						textWatermarkImg, image.Point{}, draw.Over)
				} else {
					// 创建临时画布并设置不透明度
					tmp := imaging.New(watermarkBounds.Dx(), watermarkBounds.Dy(), color.Transparent)
					tmp = imaging.Overlay(tmp, textWatermarkImg, image.Point{}, opacity)
					draw.Draw(result, image.Rect(x, y, x+watermarkBounds.Dx(), y+watermarkBounds.Dy()),
						tmp, image.Point{}, draw.Over)
				}
			}
		}
		destImg = result
	default:
		return errors.New("watermark position error")
	}

	// 保存结果图片
	if err = imaging.Save(destImg, config.CompositeImagePath); err != nil {
		return errors.New("create composite image error:" + err.Error())
	}
	return nil
}

// CreateTextImage 创建文字图像
func CreateTextImage(text, fontPath string, size float64, color color.RGBA, rotation float64) (*image.NRGBA, error) {
	// 加载字体文件
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, errors.New("failed to read font file:" + err.Error())
	}

	fontFace, err := truetype.Parse(fontData)
	if err != nil {
		return nil, errors.New("failed to parse font:" + err.Error())
	}

	// 设置字体大小
	var fontSize float64 = 50
	face := truetype.NewFace(fontFace, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	// 计算文字的宽度和高度
	textWidth, textHeight := measureText(face, text)
	// 创建一个新的图像
	img := image.NewRGBA(image.Rect(0, 0, textWidth, textHeight))
	// 绘制背景颜色透明
	draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)
	// 设置文字颜色并添加透明度
	//textColor := color.RGBA{R: 0, G: 0, B: 0, A: 200} // 黑色，半透明
	textColor := color.RGBA{R: color.R, G: color.G, B: color.B, A: 255}
	// 绘制文字
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}
	// 设置绘制起点
	d.Dot = fixed.P(0, textHeight-5)
	// 绘制文字
	d.DrawString(text)
	// 图片旋转
	dst := imaging.Rotate(img, rotation, color.Transparent)

	return dst, nil
}

// measureText 计算给定文字的宽度和高度
func measureText(face font.Face, text string) (int, int) {
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
