# gowatermark

使用 Go 语言开发图片水印工具，可以添加图片和文字水印。

## 安装

```
go get -u github.com/xing-zr/gowatermark
```

## 示例

### 添加图片水印

```go
config := gowatermark.ImageWatermarkConfig{
    OriginImagePath:    "./testdata/origin.jpg",
    WatermarkImagePath: "./testdata/watermark.png",
    WatermarkPos:       LeftTop,
    CompositeImagePath: "./testdata/composite.jpg",
}
gowatermark.CreateImageWatermark(config)
```

### 添加文字水印

```go
config := gowatermark.TextWatermarkConfig{
    OriginImagePath:    "./testdata/origin.jpg",
    CompositeImagePath: "./testdata/composite.jpg",
    FontPath:           "./testdata/font.ttf",
    TextInfos: []TextInfo{
        {
           Size: 100,
           Text: "hello world",
		   X:    700,
           Y:    700,
        },
    },
}
gowatermark.CreateTextWatermark(config)
```