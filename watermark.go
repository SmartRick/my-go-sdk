package go_watermark

import (
	"errors"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
	"os"
	"path/filepath"
)

type ImageWatermarkConfig struct {
	OriginImagePath    string       // 原图地址
	WatermarkImagePath string       // 水印图地址
	WatermarkPos       watermarkPos // 水印位置
	CompositeImagePath string       // 合成图地址
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
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(10, 10), 1)
	case RightTop:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(originImgWidth-int(targetWatermarkImgWidth)-10, 10), 1)
	case LeftBottom:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(10, originImgHeight-destwatermarkImg.Bounds().Dy()-10), 1)
	case RightBottom:
		destImg = imaging.Overlay(originImg, destwatermarkImg, image.Pt(originImgWidth-int(targetWatermarkImgWidth)-10, originImgHeight-destwatermarkImg.Bounds().Dy()-10), 1)
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
