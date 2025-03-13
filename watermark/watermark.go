package watermark

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"

	"github.com/SmartRick/my-go-sdk/common"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ImageWatermarkConfig struct {
	OriginImagePath    string       // 原图地址
	WatermarkImagePath string       // 水印图地址
	WatermarkPos       WatermarkPos // 水印位置
	CompositeImagePath string       // 合成图地址
	OffsetX            int          // 水印位置偏移量X
	OffsetY            int          // 水印位置偏移量Y
	Opacity            float64      // 水印透明度
	TiledRows          int          // 水印图横向平铺行数
	TiledCols          int          // 水印图横向平铺列数
}

type WatermarkPos string

const (
	LeftTop     WatermarkPos = "left_top"
	RightTop    WatermarkPos = "right_top"
	LeftBottom  WatermarkPos = "left_bottom"
	RightBottom WatermarkPos = "right_bottom"
	Tiled       WatermarkPos = "tiled"
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
	// 如果合成图片存在则删除重新生成
	isExists, _ := common.PathExists(config.CompositeImagePath)
	if isExists {
		err = os.Remove(config.CompositeImagePath)
		if err != nil {
			return errors.New("old composite image remove error:" + err.Error())
		}
	}
	// 判断文件夹是否存在，不存在创建
	dirPath := filepath.Dir(config.CompositeImagePath)
	isExist, _ := common.PathExists(dirPath)
	if !isExist {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
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

// TransparentTextWatermarkConfig 透明文字水印配置
type TransparentTextWatermarkConfig struct {
	OriginImagePath    string       // 原图地址
	CompositeImagePath string       // 合成图地址
	FontPath           string       // 字体文件地址（可选，为空时使用系统默认字体）
	Text               string       // 文字内容
	Size               float64      // 文字大小
	Color              color.RGBA   // 文字颜色
	WatermarkPos       WatermarkPos // 水印位置
	Opacity            float64      // 水印透明度
	OffsetX            int          // 水印位置偏移量X
	OffsetY            int          // 水印位置偏移量Y
	Rotation           float64      // 文字旋转角度
	TiledRows          int          // 水印图横向平铺行数(仅Tiled位置时使用)
	TiledCols          int          // 水印图横向平铺列数(仅Tiled位置时使用)
}

// 创建几个预选颜色
var (
	Black = color.RGBA{0, 0, 0, 255}
	White = color.RGBA{255, 255, 255, 255}
	Red   = color.RGBA{255, 0, 0, 255}
	Green = color.RGBA{0, 255, 0, 255}
	Blue  = color.RGBA{0, 0, 255, 255}
)

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

// CreateTransparentTextWatermark 创建透明文字水印
// 先将文字渲染到透明图层，然后作为图片叠加到目标图片上
func CreateTransparentTextWatermark(config TransparentTextWatermarkConfig) (image.Image, error) {
	// 输入参数验证
	if config.Opacity < 0 || config.Opacity > 1 {
		return nil, errors.New("watermark opacity error: Ensure 0.0 <= opacity <= 1.0")
	}
	if config.Opacity == 0 {
		config.Opacity = 1
	}
	if config.WatermarkPos == Tiled && (config.TiledCols == 0 || config.TiledRows == 0) {
		return nil, errors.New("watermark position tiled need tiled_cols and tiled_rows")
	}

	// 打开原始图片
	originFile, err := os.Open(config.OriginImagePath)
	if err != nil {
		return nil, errors.New("open origin image file error:" + err.Error())
	}
	defer originFile.Close()

	// 解码原始图片
	originImg, err := imaging.Decode(originFile)
	if err != nil {
		return nil, errors.New("decode origin image error: " + err.Error())
	}

	// 创建文字水印图像
	textWatermarkImg, err := createTextImage(config)
	if err != nil {
		return nil, err
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
		return nil, errors.New("watermark position error")
	}

	// 如果指定了输出路径，则保存结果图片
	if config.CompositeImagePath != "" {
		// 处理输出路径
		isExists, _ := common.PathExists(config.CompositeImagePath)
		if isExists {
			err = os.Remove(config.CompositeImagePath)
			if err != nil {
				return nil, errors.New("old composite image remove error:" + err.Error())
			}
		}

		dirPath := filepath.Dir(config.CompositeImagePath)
		isExist, _ := common.PathExists(dirPath)
		if !isExist {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				return nil, err
			}
		}

		if err = imaging.Save(destImg, config.CompositeImagePath); err != nil {
			return nil, errors.New("create composite image error:" + err.Error())
		}
	}
	return destImg, nil
}

// createTextImage 创建文字图像
func createTextImage(config TransparentTextWatermarkConfig) (*image.NRGBA, error) {
	// 加载字体文件
	fontData, err := os.ReadFile(config.FontPath)
	if err != nil {
		return nil, errors.New("failed to read font file:" + err.Error())
	}

	fontFace, err := truetype.Parse(fontData)
	if err != nil {
		return nil, errors.New("failed to parse font:" + err.Error())
	}

	// 设置字体大小和选项
	face := truetype.NewFace(fontFace, &truetype.Options{
		Size:    config.Size,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// 计算文字的宽度和高度
	textWidth, textHeight := measureText(face, config.Text)

	// 文字需要旋转时，确保最终图像足够大以容纳旋转后的文本
	padding := 0
	if config.Rotation != 0 {
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
		Src:  image.NewUniform(config.Color),
		Face: face,
	}

	// 设置绘制起点（添加边距）
	d.Dot = fixed.P(padding, textHeight+padding-5)

	// 绘制文字
	d.DrawString(config.Text)

	// 如果需要旋转
	var dst *image.NRGBA
	if config.Rotation != 0 {
		dst = imaging.Rotate(img, config.Rotation, color.Transparent)
	} else {
		dst = imaging.Clone(img)
	}

	return dst, nil
}
