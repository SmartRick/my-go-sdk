package excel

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 示例：展示Excel工具包的基本使用方法
func ExampleBasicUsage() {
	// 创建一个新的Excel处理器
	processor := NewExcelProcessor()

	// 设置单元格值
	processor.SetCellValue("A1", "姓名")
	processor.SetCellValue("B1", "年龄")
	processor.SetCellValue("C1", "生日")
	processor.SetCellValue("D1", "工资")
	processor.SetCellValue("E1", "部门")

	// 添加数据
	processor.SetCellValue("A2", "张三")
	processor.SetCellValue("B2", 28)
	processor.SetCellValue("C2", time.Date(1995, 3, 15, 0, 0, 0, 0, time.Local))
	processor.SetCellValue("D2", 8500.50)
	processor.SetCellValue("E2", "研发部")

	processor.SetCellValue("A3", "李四")
	processor.SetCellValue("B3", 32)
	processor.SetCellValue("C3", time.Date(1991, 7, 8, 0, 0, 0, 0, time.Local))
	processor.SetCellValue("D3", 12500.80)
	processor.SetCellValue("E3", "市场部")

	// 设置列宽
	processor.SetColumnWidth("A", "A", 12)
	processor.SetColumnWidth("B", "B", 8)
	processor.SetColumnWidth("C", "C", 12)
	processor.SetColumnWidth("D", "D", 10)
	processor.SetColumnWidth("E", "E", 12)

	// 添加新的工作表
	processor.CreateSheet("部门列表")
	processor.SetCellValue("A1", "部门ID")
	processor.SetCellValue("B1", "部门名称")
	processor.SetCellValue("C1", "部门主管")

	processor.SetCellValue("A2", 1)
	processor.SetCellValue("B2", "研发部")
	processor.SetCellValue("C2", "王经理")

	processor.SetCellValue("A3", 2)
	processor.SetCellValue("B3", "市场部")
	processor.SetCellValue("C3", "赵经理")

	// 切换回第一个工作表
	processor.SetActiveSheet("Sheet1")

	// 保存文件
	err := processor.Save("员工信息.xlsx")
	if err != nil {
		log.Fatalf("保存Excel文件失败: %v", err)
	}

	fmt.Println("Excel文件已成功创建")
}

// 示例：使用模板填充数据
func ExampleUsingTemplate() {
	// 假设有一个模板文件
	templatePath := "报表模板.xlsx"

	// 准备数据
	data := map[string]interface{}{
		"公司名称":  "ABC科技有限公司",
		"报表日期":  time.Now().Format("2006-01-02"),
		"制表人":   "系统管理员",
		"销售总额":  125680.50,
		"利润":    45990.75,
		"销售增长率": "15.8%",
	}

	// 创建模板对象
	tmpl := &ReportTemplate{
		TemplatePath: templatePath,
		Values:       data,
	}

	// 填充模板并保存为新文件
	outputPath := "销售报表_" + time.Now().Format("20060102") + ".xlsx"
	// 注意：实际使用时需要确保模板文件存在
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Printf("模板文件不存在，跳过模板示例\n")
		return
	}

	err := tmpl.FillTemplate(outputPath)
	if err != nil {
		log.Fatalf("填充模板失败: %v", err)
	}

	fmt.Println("报表已成功生成")
}

// 示例：批量处理多个Excel文件
func ExampleBatchProcessing() {
	// 假设有一个目录，包含多个Excel文件需要处理
	dirPath := "excel_files"

	// 确保目录存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0755)
		fmt.Printf("创建目录: %s\n", dirPath)
	}

	// 遍历目录下的所有Excel文件
	files, _ := filepath.Glob(filepath.Join(dirPath, "*.xlsx"))

	if len(files) == 0 {
		fmt.Println("没有Excel文件可处理")
		return
	}

	for _, file := range files {
		// 打开Excel文件
		processor, err := OpenExcelFile(file)
		if err != nil {
			log.Printf("打开文件 %s 失败: %v", file, err)
			continue
		}

		// 处理逻辑 - 这里只是示例，实际应用中根据需求来处理
		fmt.Printf("处理文件: %s\n", file)
		sheets := processor.GetSheetList()
		for _, sheet := range sheets {
			fmt.Printf("  - 工作表: %s\n", sheet)
		}

		// 关闭文件
		processor.Close()
	}

	fmt.Println("批量处理完成")
}

// 示例：创建图表
func ExampleCreateChart() {
	// 创建一个新的Excel处理器
	processor := NewExcelProcessor()

	// 设置数据
	processor.SetCellValue("A1", "季度")
	processor.SetCellValue("B1", "销售额(万元)")

	processor.SetCellValue("A2", "Q1")
	processor.SetCellValue("B2", 124.5)

	processor.SetCellValue("A3", "Q2")
	processor.SetCellValue("B3", 168.3)

	processor.SetCellValue("A4", "Q3")
	processor.SetCellValue("B4", 152.7)

	processor.SetCellValue("A5", "Q4")
	processor.SetCellValue("B5", 189.2)

	// 保存文件
	err := processor.Save("季度销售图表.xlsx")
	if err != nil {
		log.Fatalf("保存Excel文件失败: %v", err)
	}

	fmt.Println("Excel图表文件已创建 (注：excelize库支持添加图表，但此示例中未实现，请参考官方文档)")
}

// 示例：导出为其他格式
func ExampleExportToOtherFormats() {
	// 创建一个新的Excel处理器
	processor := NewExcelProcessor()

	// 设置数据
	processor.SetCellValue("A1", "产品ID")
	processor.SetCellValue("B1", "产品名称")
	processor.SetCellValue("C1", "价格")
	processor.SetCellValue("D1", "库存")

	processor.SetCellValue("A2", "P001")
	processor.SetCellValue("B2", "笔记本电脑")
	processor.SetCellValue("C2", 5699)
	processor.SetCellValue("D2", 125)

	processor.SetCellValue("A3", "P002")
	processor.SetCellValue("B3", "智能手机")
	processor.SetCellValue("C3", 3299)
	processor.SetCellValue("D3", 230)

	processor.SetCellValue("A4", "P003")
	processor.SetCellValue("B4", "无线耳机")
	processor.SetCellValue("C4", 899)
	processor.SetCellValue("D4", 310)

	// 保存Excel文件
	processor.Save("产品列表.xlsx")

	// 导出为CSV
	// 注意：当前的writeRowsToCSV实现不完整，实际使用时需要修改
	// processor.ExportAsCSV("产品列表.csv")

	// 导出为HTML
	// 注意：当前的writeStringToFile实现不完整，实际使用时需要修改
	// processor.ExportAsHTML("产品列表.html")

	fmt.Println("Excel文件已导出")
}

// 主函数：运行所有示例
func RunAllExamples() {
	fmt.Println("====== Excel工具包使用示例 ======")

	fmt.Println("\n=== 基本使用示例 ===")
	ExampleBasicUsage()

	fmt.Println("\n=== 模板使用示例 ===")
	ExampleUsingTemplate()

	fmt.Println("\n=== 批量处理示例 ===")
	ExampleBatchProcessing()

	fmt.Println("\n=== 图表创建示例 ===")
	ExampleCreateChart()

	fmt.Println("\n=== 格式导出示例 ===")
	ExampleExportToOtherFormats()

	fmt.Println("\n====== 所有示例运行完毕 ======")
}
