package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"

	"github.com/SmartRick/my-go-sdk/common" // 修改为正确的导入路径
	"github.com/SmartRick/my-go-sdk/excel"  // 修改为正确的导入路径
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "excel":
			runExcelExamples()
			return
		case "string":
			runStringExamples()
			return
		case "all":
			runAllExamples()
			return
		case "help":
			printHelp()
			return
		}
	}

	// 默认操作：打印帮助信息
	printHelp()
}

func printHelp() {
	fmt.Println("使用说明：")
	fmt.Println("  my-go-sdk excel    - 运行Excel处理示例")
	fmt.Println("  my-go-sdk string   - 运行字符串处理示例")
	fmt.Println("  my-go-sdk all      - 运行所有示例")
	fmt.Println("  my-go-sdk help     - 显示此帮助信息")
}

// 运行所有示例
func runAllExamples() {
	fmt.Println("====== 运行所有功能示例 ======")

	fmt.Println("\n=== Excel处理示例 ===")
	runExcelExamples()

	fmt.Println("\n=== 字符串处理示例 ===")
	runStringExamples()

	fmt.Println("\n====== 所有示例运行完毕 ======")
}

// 运行Excel示例代码
func runExcelExamples() {
	fmt.Println("运行Excel处理包示例...")

	// 确保输出目录存在
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// 切换到输出目录
	currentDir, _ := os.Getwd()
	os.Chdir(outputDir)

	// 运行示例
	excel.RunAllExamples()

	// 额外创建一个简单的Excel文件
	createSimpleExcel()

	// 切回原目录
	os.Chdir(currentDir)
}

// 运行字符串处理示例
func runStringExamples() {
	fmt.Println("运行字符串处理包示例...")
	common.StringExamples()
}

// 创建一个简单的Excel文件
func createSimpleExcel() {
	fmt.Println("\n创建简单Excel文件...")

	// 创建Excel处理器
	processor := excel.NewExcelProcessor()

	// 添加标题行
	processor.SetCellValue("A1", "ID")
	processor.SetCellValue("B1", "名称")
	processor.SetCellValue("C1", "描述")
	processor.SetCellValue("D1", "创建日期")

	// 添加数据
	processor.SetCellValue("A2", 1)
	processor.SetCellValue("B2", "项目A")
	processor.SetCellValue("C2", "这是项目A的描述信息")
	processor.SetCellValue("D2", "2023-06-15")

	processor.SetCellValue("A3", 2)
	processor.SetCellValue("B3", "项目B")
	processor.SetCellValue("C3", "这是项目B的描述信息")
	processor.SetCellValue("D3", "2023-07-22")

	processor.SetCellValue("A4", 3)
	processor.SetCellValue("B4", "项目C")
	processor.SetCellValue("C4", "这是项目C的描述信息")
	processor.SetCellValue("D4", "2023-08-30")

	// 设置列宽
	processor.SetColumnWidth("A", "A", 8)
	processor.SetColumnWidth("B", "B", 15)
	processor.SetColumnWidth("C", "C", 30)
	processor.SetColumnWidth("D", "D", 15)

	// 创建样式
	style, err := processor.CreateStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	if err == nil {
		// 应用样式到标题行
		processor.SetCellStyle("A1", "D1", style)
	}

	// 保存文件
	outputPath := "simple_excel.xlsx"
	err = processor.Save(outputPath)
	if err != nil {
		fmt.Printf("保存Excel文件失败: %v\n", err)
		return
	}

	fmt.Printf("Excel文件已成功创建: %s\n", filepath.Join("output", outputPath))
}
