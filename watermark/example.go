package watermark

import (
	"fmt"
	"image/color"
	"os"
)

// WatermarkExamples 展示水印功能的使用示例
func WatermarkExamples() {
	fmt.Println("====== 水印功能使用示例 ======")

	// 示例1：使用默认选项添加单个水印
	fmt.Println("\n--- 示例1：使用默认选项添加单个水印 ---")
	imgBytes, err := os.ReadFile("testdata/test.jpg")
	if err != nil {
		fmt.Printf("读取图片失败: %v\n", err)
		return
	}

	// 使用默认选项添加右下角水印
	img, err := AddTextWatermark(imgBytes, "示例水印", RightBottom, nil)
	if err != nil {
		fmt.Printf("添加水印失败: %v\n", err)
		return
	}
	err = SaveImage(img, "testdata/output_single.jpg")
	if err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		return
	}
	fmt.Println("成功创建单个水印图片：testdata/output_single.jpg")

	// 示例2：自定义选项添加水印
	fmt.Println("\n--- 示例2：自定义选项添加水印 ---")
	opts := &WatermarkOptions{
		FontPath: "testdata/fonts/msyh.ttf",  // 使用微软雅黑字体
		Size:     72,                         // 较大字号
		Color:    color.RGBA{255, 0, 0, 255}, // 红色
		Opacity:  0.3,                        // 较高透明度
		Rotation: 45,                         // 45度旋转
		OffsetX:  50,                         // 较大偏移
		OffsetY:  50,
	}

	// 添加左上角水印
	img, err = AddTextWatermark(imgBytes, "自定义水印", LeftTop, opts)
	if err != nil {
		fmt.Printf("添加水印失败: %v\n", err)
		return
	}
	err = SaveImage(img, "testdata/output_custom.jpg")
	if err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		return
	}
	fmt.Println("成功创建自定义水印图片：testdata/output_custom.jpg")

	// 示例3：添加平铺水印
	fmt.Println("\n--- 示例3：添加平铺水印 ---")
	opts = &WatermarkOptions{
		FontPath: "testdata/fonts/msyh.ttf",
		Size:     36,
		Color:    color.RGBA{0, 0, 255, 255}, // 蓝色
		Opacity:  0.2,
		Rotation: 30,
	}

	// 添加3x3的平铺水印
	img, err = AddTiledTextWatermark(imgBytes, "平铺水印", 3, 3, opts)
	if err != nil {
		fmt.Printf("添加平铺水印失败: %v\n", err)
		return
	}
	err = SaveImage(img, "testdata/output_tiled.jpg")
	if err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		return
	}
	fmt.Println("成功创建平铺水印图片：testdata/output_tiled.jpg")

	// 示例4：直接处理文件
	fmt.Println("\n--- 示例4：直接处理文件 ---")
	img, err = AddTextWatermarkToFile("testdata/test.jpg", "文件水印", RightTop, nil)
	if err != nil {
		fmt.Printf("添加文件水印失败: %v\n", err)
		return
	}
	err = SaveImage(img, "testdata/output_file.jpg")
	if err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		return
	}
	fmt.Println("成功创建文件水印图片：testdata/output_file.jpg")

	fmt.Println("\n====== 水印功能示例结束 ======")
}
