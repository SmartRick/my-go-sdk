# my-go-sdk v1.0.0

一个实用的Go语言工具集合，提供丰富的功能模块，包括图片水印处理、文件操作、加解密、HTTP请求和并发控制等工具。本项目致力于简化常见的开发任务，提高开发效率。

## 安装

```bash
go get -u github.com/SmartRick/my-go-sdk
```

## 功能模块

### 1. 图片水印工具 (watermark)

提供图片和文字水印功能，支持多种水印位置和透明度设置。

详细功能：
- 图片水印：支持在图片上添加图片水印
- 文字水印：支持添加透明文字水印
- 位置控制：支持左上、右上、左下、右下、平铺等多种位置
- 透明度调节：可调节水印透明度
- 批量处理：支持批量给多张图片添加水印

[详细文档和使用示例](./watermark/README.md)

### 2. 文件工具 (common/file)

提供丰富的文件和目录操作函数，简化日常文件处理任务。

主要功能：
- 文件检测：检查文件/目录是否存在、是否为文件/目录
- 文件创建：创建文件、目录
- 文件删除：删除文件、目录
- 文件复制/移动：复制或移动文件
- 文件读写：读取、写入文件内容
- 文件遍历：列出目录中的文件和子目录
- 文件信息：获取文件大小、扩展名、修改时间等信息
- 其他实用工具：MD5计算、文件大小格式化等

### 3. 加解密工具 (common/crypto)

提供常用的加密、解密和编码功能，简化数据安全处理。

主要功能：
- Base64编解码：标准Base64和URL安全的Base64编解码
- URL编解码：字符串URL编码和解码
- 哈希计算：MD5、SHA1、SHA256等哈希值计算
- 密码处理：使用bcrypt对密码进行安全哈希和验证
- AES加解密：使用AES-GCM算法进行对称加密和解密
- RSA加解密：RSA公钥/私钥生成和加解密
- 随机数生成：生成安全的随机字节、字符串和令牌

### 4. HTTP请求工具 (common/http)

提供简洁易用的HTTP客户端，支持各种HTTP请求方式和高级功能。

主要功能：
- 基本请求：GET、POST、PUT、DELETE等请求方法
- JSON处理：自动处理JSON请求和响应
- 表单处理：支持表单数据提交
- 文件上传：简化文件上传过程
- 认证支持：Basic认证、Bearer令牌认证
- 重试机制：自动重试失败的请求
- 超时控制：请求超时设置
- 下载功能：简化文件下载

### 5. 并发线程工具 (common/concurrency)

提供并发编程相关的工具，简化多线程任务处理。

主要功能：
- 线程池：管理和重用goroutine，避免无限制创建
- 并行执行器：并行执行多个任务，支持超时控制
- 线程安全容器：线程安全的计数器、映射等数据结构
- 信号量：控制并发访问资源
- 速率限制器：限制操作频率
- 批量处理：将大量数据分批并行处理

### 6. 字符串处理工具 (common/string)

提供丰富而强大的字符串处理函数，支持各种字符串操作需求。

主要功能：
- 字符串判断：空值判断、类型判断、内容判断
- 字符串转换：驼峰命名、蛇形命名、短横线命名等格式转换
- 大小写转换：首字母大小写、全部大小写转换
- 字符串操作：填充、截断、反转、掩码处理
- 字符串提取：正则提取、标记间提取、子串提取
- 字符串解析：数值解析、布尔值解析
- 高级功能：文本换行、居中对齐、随机子串生成、对比差异等
- 特殊功能：URL友好字符串生成、特殊格式验证

### 7. Excel处理工具 (excel)

提供功能全面的Excel文件处理工具，支持创建、读取、修改Excel文件。

主要功能：
- 文件操作：创建、打开、保存Excel文件
- 工作表管理：创建、删除、切换工作表
- 单元格操作：读写单元格、设置公式、合并单元格
- 样式设置：字体、颜色、边框、对齐方式
- 数据导入导出：从数据结构导入/导出Excel
- 格式转换：Excel与CSV、HTML等格式的互相转换
- 报表模板：支持模板变量替换生成报表
- 批量处理：批量处理多个Excel文件
- 实用工具：日期转换、单元格坐标转换等

#### 示例

```go
package main

import (
	"fmt"
	"github.com/SmartRick/my-go-sdk/common"
	"github.com/SmartRick/my-go-sdk/excel"
)

func main() {
	// 检查文件是否存在
	exists, _ := common.PathExists("./myfile.txt")
	fmt.Println("文件存在:", exists)
	
	// 创建目录
	err := common.CreateDir("./mydir/subdir")
	if err != nil {
		fmt.Println("创建目录失败:", err)
	}
	
	// 写入文件
	err = common.WriteFile("./myfile.txt", "Hello, World!")
	if err != nil {
		fmt.Println("写入文件失败:", err)
	}
	
	// 读取文件
	content, err := common.ReadFile("./myfile.txt")
	if err != nil {
		fmt.Println("读取文件失败:", err)
	}
	fmt.Println("文件内容:", content)
	
	// 获取文件大小
	size, _ := common.GetFileSize("./myfile.txt")
	fmt.Println("文件大小:", common.FormatFileSize(size))
	
	// 计算文件MD5
	md5, _ := common.GetFileMD5("./myfile.txt")
	fmt.Println("文件MD5:", md5)
	
	// 字符串处理
	fmt.Println("\n字符串处理示例:")
	fmt.Println("是否只包含数字:", common.IsNumeric("12345"))
	fmt.Println("转换为蛇形命名:", common.ToSnakeCase("helloWorld"))
	fmt.Println("掩码处理:", common.MaskMiddleChars("13812345678", 3, 4, '*'))
	fmt.Println("提取标签内容:", common.ExtractBetween("<div>内容</div>", "<div>", "</div>"))
	
	// Excel处理
	fmt.Println("\nExcel处理示例:")
	// 创建新的Excel文件
	processor := excel.NewExcelProcessor()
	// 添加数据
	processor.SetCellValue("A1", "ID")
	processor.SetCellValue("B1", "名称")
	processor.SetCellValue("A2", 1)
	processor.SetCellValue("B2", "测试项目")
	// 保存文件
	processor.Save("example.xlsx")
	fmt.Println("Excel文件已创建: example.xlsx")
}

### 使用命令行工具

本SDK还提供了命令行工具，可以快速测试各种功能：

```bash
# 安装
go install github.com/SmartRick/my-go-sdk@latest

# 查看帮助信息
my-go-sdk help

# 运行Excel处理示例
my-go-sdk excel

# 运行字符串处理示例
my-go-sdk string

# 运行所有示例
my-go-sdk all
```

## 后续开发计划

以下是计划添加的功能模块：

1. **日志工具**：提供灵活的日志记录功能，支持日志级别、日志轮转、多输出目标等

2. **配置管理**：支持多种配置格式(JSON, YAML, TOML)的读取和管理

3. **数据库工具**：数据库连接池管理、简化的CRUD操作

4. **时间工具**：时间格式化、时区转换、计时器等

5. **压缩解压**：文件压缩和解压缩功能

6. **PDF处理**：PDF生成、转换、合并等功能

## 贡献

欢迎提交PR或Issue，一起完善这个工具集！

## 许可证

本项目基于MIT许可证开源，详情请查看[LICENSE](./LICENSE)文件。 