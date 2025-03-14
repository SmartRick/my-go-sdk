# gowatermark

使用 Go 语言开发图片水印工具，可以添加图片和透明文字水印。

## 安装

```
go get -u github.com/xing-zr/gowatermark
```

## 示例

### 添加图片水印

```golang
// 相关配置
config := gowatermark.ImageWatermarkConfig{
	OriginImagePath:    "./origin.jpg",        // 水印底图图片路径
	WatermarkImagePath: "./watermark.png",     // 水印图图片路径
	WatermarkPos:       gowatermark.LeftTop,   // 水印位置 有 左上、左下、右上、右下、平铺几个选项
	CompositeImagePath: "./composite.jpg",     // 合成后图片路径
    Opacity:   0.5,                            // 水印透明度 0-1 之间
    TiledRows: 3,                              // 水印平铺行数（水印类型为平铺有效）
    TiledCols: 4,                              // 水印平铺列数（水印类型为平铺有效）
}
// 生成图片水印图
err := gowatermark.CreateImageWatermark(config)
if err != nil {
	fmt.Println(err)
}
```

### 添加透明文字水印

```golang
// 相关配置 - 固定位置的透明文字水印
config := gowatermark.TransparentTextWatermarkConfig{
    OriginImagePath:    "./origin.jpg",           // 水印底图图片路径
    CompositeImagePath: "./composite.jpg",        // 合成后图片路径
    FontPath:           "./font.ttf",             // 字体文件路径
    Text:               "透明文字水印",             // 文字内容
    Size:               72,                       // 字体大小
    Color:              gowatermark.White,        // 文字颜色，提供预设颜色：White, Black, Red, Green, Blue
    WatermarkPos:       gowatermark.RightBottom, // 水印位置，支持左上、右上、左下、右下、平铺
    Opacity:            0.5,                     // 水印透明度 0-1 之间
    OffsetX:            20,                      // X轴偏移量
    OffsetY:            20,                      // Y轴偏移量
    Rotation:           0,                       // 文字旋转角度
}
// 生成透明文字水印图
err := gowatermark.CreateTransparentTextWatermark(config)
if err != nil {
    fmt.Println(err)
}

// 相关配置 - 平铺透明文字水印
configTiled := gowatermark.TransparentTextWatermarkConfig{
    OriginImagePath:    "./origin.jpg",           // 水印底图图片路径
    CompositeImagePath: "./composite_tiled.jpg",  // 合成后图片路径
    FontPath:           "./font.ttf",             // 字体文件路径
    Text:               "透明文字水印",             // 文字内容
    Size:               48,                       // 字体大小
    Color:              gowatermark.White,        // 文字颜色，提供预设颜色
    WatermarkPos:       gowatermark.Tiled,       // 水印位置设为平铺
    Opacity:            0.3,                     // 水印透明度 0-1 之间
    TiledRows:          4,                       // 平铺行数
    TiledCols:          5,                       // 平铺列数
    Rotation:           45,                      // 文字旋转角度
}
// 生成平铺透明文字水印图
err = gowatermark.CreateTransparentTextWatermark(configTiled)
if err != nil {
    fmt.Println(err)
}
```

### 字体支持

该库支持以下两种方式指定字体：

1. 显式指定字体文件路径
```golang
config := gowatermark.TransparentTextWatermarkConfig{
    // ... other config ...
    FontPath: "./my-custom-font.ttf",
}
```

2. 使用系统默认字体（将FontPath留空）
```golang
config := gowatermark.TransparentTextWatermarkConfig{
    // ... other config ...
    FontPath: "", // 将使用系统默认字体
}
```

系统默认字体会按以下顺序查找：
- Linux: /usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf
- Windows: C:/Windows/Fonts/simhei.ttf
- macOS: /System/Library/Fonts/PingFang.ttc

如果以上路径都无法找到可用字体，则会返回错误。

## 可用水印位置

- `gowatermark.LeftTop` - 左上角
- `gowatermark.RightTop` - 右上角
- `gowatermark.LeftBottom` - 左下角
- `gowatermark.RightBottom` - 右下角 
- `gowatermark.Tiled` - 平铺模式

## 预设颜色

- `gowatermark.White` - 白色
- `gowatermark.Black` - 黑色
- `gowatermark.Red` - 红色
- `gowatermark.Green` - 绿色
- `gowatermark.Blue` - 蓝色