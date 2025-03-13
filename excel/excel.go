package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelProcessor Excel处理器
type ExcelProcessor struct {
	file       *excelize.File
	sheetName  string
	activeCell string // 当前活动单元格，例如"A1"
}

// NewExcelProcessor 创建新的Excel处理器
func NewExcelProcessor() *ExcelProcessor {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(0) // 获取默认工作表名
	return &ExcelProcessor{
		file:       file,
		sheetName:  sheetName,
		activeCell: "A1",
	}
}

// OpenExcelFile 打开Excel文件
func OpenExcelFile(filePath string) (*ExcelProcessor, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	sheetName := file.GetSheetName(0) // 获取第一个工作表名
	return &ExcelProcessor{
		file:       file,
		sheetName:  sheetName,
		activeCell: "A1",
	}, nil
}

// Save 保存Excel文件
func (p *ExcelProcessor) Save(filePath string) error {
	if filePath == "" {
		return p.file.Save()
	}
	return p.file.SaveAs(filePath)
}

// Close 关闭Excel文件
func (p *ExcelProcessor) Close() error {
	return p.file.Close()
}

// CreateSheet 创建新工作表
func (p *ExcelProcessor) CreateSheet(sheetName string) int {
	index, err := p.file.NewSheet(sheetName)
	if err != nil {
		return -1
	}
	p.sheetName = sheetName
	p.activeCell = "A1"
	return index
}

// SetActiveSheet 设置活动工作表
func (p *ExcelProcessor) SetActiveSheet(sheetName string) error {
	if !p.SheetExists(sheetName) {
		return fmt.Errorf("工作表 %s 不存在", sheetName)
	}
	p.sheetName = sheetName
	p.activeCell = "A1"
	return nil
}

// SheetExists 检查工作表是否存在
func (p *ExcelProcessor) SheetExists(sheetName string) bool {
	sheets := p.file.GetSheetList()
	for _, sheet := range sheets {
		if sheet == sheetName {
			return true
		}
	}
	return false
}

// RemoveSheet 删除工作表
func (p *ExcelProcessor) RemoveSheet(sheetName string) error {
	if !p.SheetExists(sheetName) {
		return fmt.Errorf("工作表 %s 不存在", sheetName)
	}
	if p.sheetName == sheetName {
		// 如果删除的是当前活动工作表，需要设置活动工作表为别的工作表
		sheets := p.file.GetSheetList()
		for _, sheet := range sheets {
			if sheet != sheetName {
				p.sheetName = sheet
				p.activeCell = "A1"
				break
			}
		}
	}
	return p.file.DeleteSheet(sheetName)
}

// GetSheetList 获取所有工作表列表
func (p *ExcelProcessor) GetSheetList() []string {
	return p.file.GetSheetList()
}

// SetCellValue 设置单元格值
func (p *ExcelProcessor) SetCellValue(cell string, value interface{}) error {
	p.activeCell = cell
	return p.file.SetCellValue(p.sheetName, cell, value)
}

// GetCellValue 获取单元格值
func (p *ExcelProcessor) GetCellValue(cell string) (string, error) {
	p.activeCell = cell
	return p.file.GetCellValue(p.sheetName, cell)
}

// SetCellFormula 设置单元格公式
func (p *ExcelProcessor) SetCellFormula(cell, formula string) error {
	p.activeCell = cell
	return p.file.SetCellFormula(p.sheetName, cell, formula)
}

// GetCellFormula 获取单元格公式
func (p *ExcelProcessor) GetCellFormula(cell string) (string, error) {
	p.activeCell = cell
	return p.file.GetCellFormula(p.sheetName, cell)
}

// SetColumnWidth 设置列宽度
func (p *ExcelProcessor) SetColumnWidth(startCol, endCol string, width float64) error {
	return p.file.SetColWidth(p.sheetName, startCol, endCol, width)
}

// SetRowHeight 设置行高度
func (p *ExcelProcessor) SetRowHeight(row int, height float64) error {
	return p.file.SetRowHeight(p.sheetName, row, height)
}

// MergeCell 合并单元格
func (p *ExcelProcessor) MergeCell(startCell, endCell string) error {
	return p.file.MergeCell(p.sheetName, startCell, endCell)
}

// UnmergeCell 取消合并单元格
func (p *ExcelProcessor) UnmergeCell(startCell, endCell string) error {
	return p.file.UnmergeCell(p.sheetName, startCell, endCell)
}

// SetCellStyle 设置单元格样式
func (p *ExcelProcessor) SetCellStyle(startCell, endCell string, styleID int) error {
	return p.file.SetCellStyle(p.sheetName, startCell, endCell, styleID)
}

// CreateStyle 创建样式
func (p *ExcelProcessor) CreateStyle(style *excelize.Style) (int, error) {
	return p.file.NewStyle(style)
}

// InsertRow 插入行
func (p *ExcelProcessor) InsertRow(row int) error {
	return p.file.InsertRows(p.sheetName, row, 1)
}

// RemoveRow 删除行
func (p *ExcelProcessor) RemoveRow(row int) error {
	return p.file.RemoveRow(p.sheetName, row)
}

// InsertCol 插入列
func (p *ExcelProcessor) InsertCol(col string) error {
	return p.file.InsertCols(p.sheetName, col, 1)
}

// RemoveCol 删除列
func (p *ExcelProcessor) RemoveCol(col string) error {
	return p.file.RemoveCol(p.sheetName, col)
}

// SetSheetBackground 设置工作表背景图片
func (p *ExcelProcessor) SetSheetBackground(picturePath string) error {
	return p.file.SetSheetBackground(p.sheetName, picturePath)
}

// AddPicture 插入图片
func (p *ExcelProcessor) AddPicture(cell, picturePath string, widthScale, heightScale float64) error {
	p.activeCell = cell
	return p.file.AddPicture(p.sheetName, cell, picturePath, &excelize.GraphicOptions{
		ScaleX: widthScale,
		ScaleY: heightScale,
	})
}

// SetCellHyperlink 设置单元格超链接
func (p *ExcelProcessor) SetCellHyperlink(cell, linkType, location, tooltip string) error {
	p.activeCell = cell
	// 创建HyperlinkOpts对象，传递字符串指针
	display := tooltip
	tooltipPtr := tooltip
	opts := excelize.HyperlinkOpts{
		Display: &display,
		Tooltip: &tooltipPtr,
	}
	return p.file.SetCellHyperLink(p.sheetName, cell, location, linkType, opts)
}

// DataImporter 数据导入接口
type DataImporter interface {
	// Import 导入数据
	Import(processor *ExcelProcessor) error
}

// DataExporter 数据导出接口
type DataExporter interface {
	// Export 导出数据
	Export(processor *ExcelProcessor) error
}

// --------------------------------
// 便捷函数
// --------------------------------

// ReadExcel 读取Excel文件内容
func ReadExcel(filePath string) (map[string][][]string, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string][][]string)
	sheets := file.GetSheetList()

	for _, sheet := range sheets {
		rows, err := file.GetRows(sheet)
		if err != nil {
			continue
		}
		result[sheet] = rows
	}

	return result, nil
}

// CreateExcel 创建并保存Excel文件
func CreateExcel(filePath string, data map[string][][]interface{}) error {
	file := excelize.NewFile()
	defaultSheet := file.GetSheetName(0)

	sheetIndex := 0
	for sheet, rows := range data {
		// 创建工作表
		if sheetIndex == 0 {
			// 第一个工作表，使用默认工作表
			file.SetSheetName(defaultSheet, sheet)
		} else {
			file.NewSheet(sheet)
		}

		// 写入数据
		for rowIndex, row := range rows {
			for colIndex, cell := range row {
				cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				if err != nil {
					return err
				}
				file.SetCellValue(sheet, cellName, cell)
			}
		}

		sheetIndex++
	}

	// 保存文件
	return file.SaveAs(filePath)
}

// ExcelToCSV 将Excel文件转换为CSV文件
func ExcelToCSV(excelPath, csvPath string, sheetName string) error {
	// 打开Excel文件
	file, err := excelize.OpenFile(excelPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 如果未指定工作表，使用第一个工作表
	if sheetName == "" {
		sheetName = file.GetSheetName(0)
	}

	// 获取工作表数据
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return err
	}

	// 创建CSV文件
	return writeRowsToCSV(csvPath, rows)
}

// writeRowsToCSV 将行数据写入CSV文件
func writeRowsToCSV(filePath string, rows [][]string) error {
	file := excelize.NewFile()

	sheet := "Sheet1"
	for rowIndex, row := range rows {
		for colIndex, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return err
			}
			file.SetCellValue(sheet, cellName, cell)
		}
	}

	// 保存为CSV
	return file.SaveAs(filePath)
}

// CellRangeToSlice 将单元格范围转换为二维数组
func CellRangeToSlice(file *excelize.File, sheet, startCell, endCell string) ([][]string, error) {
	// 解析单元格范围
	startCol, startRow, err := excelize.CellNameToCoordinates(startCell)
	if err != nil {
		return nil, err
	}
	endCol, endRow, err := excelize.CellNameToCoordinates(endCell)
	if err != nil {
		return nil, err
	}

	// 确保起始单元格在左上角
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}

	// 读取数据
	var result [][]string
	for row := startRow; row <= endRow; row++ {
		var rowData []string
		for col := startCol; col <= endCol; col++ {
			cellName, err := excelize.CoordinatesToCellName(col, row)
			if err != nil {
				return nil, err
			}
			val, err := file.GetCellValue(sheet, cellName)
			if err != nil {
				return nil, err
			}
			rowData = append(rowData, val)
		}
		result = append(result, rowData)
	}

	return result, nil
}

// --------------------------------
// 报表相关函数
// --------------------------------

// ReportTemplate 报表模板
type ReportTemplate struct {
	TemplatePath string
	Values       map[string]interface{}
}

// FillTemplate 填充模板
func (t *ReportTemplate) FillTemplate(outputPath string) error {
	// 打开模板文件
	file, err := excelize.OpenFile(t.TemplatePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取所有工作表
	sheets := file.GetSheetList()

	// 在每个工作表中替换变量
	for _, sheet := range sheets {
		// 获取工作表中的行
		rows, err := file.GetRows(sheet)
		if err != nil {
			continue
		}

		// 遍历每个单元格
		for rowIndex, row := range rows {
			for colIndex, cell := range row {
				if cell != "" && strings.Contains(cell, "${") {
					// 查找并替换变量
					newValue := t.replaceTemplateVars(cell)
					cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
					if err != nil {
						continue
					}
					file.SetCellValue(sheet, cellName, newValue)
				}
			}
		}
	}

	// 保存新文件
	return file.SaveAs(outputPath)
}

// replaceTemplateVars 替换模板变量
func (t *ReportTemplate) replaceTemplateVars(text string) string {
	result := text
	for key, value := range t.Values {
		placeholder := fmt.Sprintf("${%s}", key)
		if strings.Contains(result, placeholder) {
			valueStr := fmt.Sprintf("%v", value)
			result = strings.ReplaceAll(result, placeholder, valueStr)
		}
	}
	return result
}

// --------------------------------
// 实用工具函数
// --------------------------------

// DetectExcelFormat 检测Excel文件格式（xls或xlsx）
func DetectExcelFormat(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".xlsx" || ext == ".xlsm" || ext == ".xltx" || ext == ".xltm" {
		return "xlsx"
	} else if ext == ".xls" || ext == ".xlt" {
		return "xls"
	}
	return "unknown"
}

// ConvertDateToCellValue 将日期转换为Excel单元格值
func ConvertDateToCellValue(date time.Time) float64 {
	// Excel基准日期是1900年1月1日，但有一个误差 (1900年不是闰年但Excel认为是)
	// 所以对于1900年3月1日之后的日期，我们需要-1来修正这个误差
	baseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	days := date.Sub(baseDate).Hours() / 24

	if days > 59 { // 1900年2月29日是第60天，它实际上不存在
		days--
	}

	return days
}

// ConvertCellValueToDate 将Excel单元格值转换为日期
func ConvertCellValueToDate(excelDate float64) time.Time {
	// 处理Excel日期值的特殊情况
	if excelDate > 60 {
		excelDate-- // 修正Excel关于1900年2月29日的错误
	}

	baseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	days := int(excelDate)
	seconds := int((excelDate-float64(days))*86400 + 0.5)

	return baseDate.AddDate(0, 0, days).Add(time.Second * time.Duration(seconds))
}

// ColumnLetterToNumber 将列字母转换为数字（如：A->1, Z->26, AA->27）
func ColumnLetterToNumber(column string) (int, error) {
	col, _, err := excelize.CellNameToCoordinates(column + "1")
	return col, err
}

// NumberToColumnLetter 将数字转换为列字母（如：1->A, 26->Z, 27->AA）
func NumberToColumnLetter(column int) (string, error) {
	cellName, err := excelize.CoordinatesToCellName(column, 1)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(cellName, "1"), nil
}

// SimplifyCSS 简化CSS样式定义
func SimplifyCSS(styleID int, processor *ExcelProcessor) {
	// 这个函数可以在将来实现，用于简化Excel样式的应用
	// 目前只是一个占位符
}

// --------------------------------
// 高级功能
// --------------------------------

// AutoFilter 设置自动筛选
func (p *ExcelProcessor) AutoFilter(startCell, endCell string) error {
	return p.file.AutoFilter(p.sheetName, startCell+":"+endCell, nil)
}

// AddDataValidation 添加数据验证
func (p *ExcelProcessor) AddDataValidation(startCell, endCell string, validationType string, criteria interface{}) error {
	dv := excelize.NewDataValidation(true)
	dv.SetSqref(startCell + ":" + endCell)

	// 设置验证类型
	dv.Type = validationType

	// 设置验证条件
	switch c := criteria.(type) {
	case []string:
		err := dv.SetDropList(c)
		if err != nil {
			return err
		}
	case []float64:
		if len(c) >= 2 {
			var operator excelize.DataValidationOperator
			operator = excelize.DataValidationOperatorBetween
			var dvType excelize.DataValidationType

			switch validationType {
			case "decimal":
				dvType = excelize.DataValidationTypeDecimal
			case "whole":
				dvType = excelize.DataValidationTypeWhole
			case "date":
				dvType = excelize.DataValidationTypeDate
			case "time":
				dvType = excelize.DataValidationTypeTime
			case "textLength":
				dvType = excelize.DataValidationTypeTextLength
			default:
				dvType = excelize.DataValidationTypeDecimal
			}

			err := dv.SetRange(c[0], c[1], dvType, operator)
			if err != nil {
				return err
			}
		}
	case string:
		dv.Formula1 = c
	}

	return p.file.AddDataValidation(p.sheetName, dv)
}

// BatchSetValues 批量设置单元格值
func (p *ExcelProcessor) BatchSetValues(data map[string]interface{}) error {
	for cell, value := range data {
		err := p.file.SetCellValue(p.sheetName, cell, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExportAsCSV 将当前工作表导出为CSV
func (p *ExcelProcessor) ExportAsCSV(csvPath string) error {
	// 获取当前工作表数据
	rows, err := p.file.GetRows(p.sheetName)
	if err != nil {
		return err
	}

	// 创建CSV文件
	return writeRowsToCSV(csvPath, rows)
}

// ExportAsHTML 将当前工作表导出为HTML表格
func (p *ExcelProcessor) ExportAsHTML(htmlPath string) error {
	// 获取当前工作表数据
	rows, err := p.file.GetRows(p.sheetName)
	if err != nil {
		return err
	}

	// 构建HTML内容
	var html strings.Builder
	html.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	html.WriteString("<meta charset=\"UTF-8\">\n")
	html.WriteString("<title>Excel Export</title>\n")
	html.WriteString("<style>\n")
	html.WriteString("table { border-collapse: collapse; width: 100%; }\n")
	html.WriteString("th, td { border: 1px solid #ddd; padding: 8px; }\n")
	html.WriteString("tr:nth-child(even) { background-color: #f2f2f2; }\n")
	html.WriteString("th { padding-top: 12px; padding-bottom: 12px; text-align: left; background-color: #4CAF50; color: white; }\n")
	html.WriteString("</style>\n")
	html.WriteString("</head>\n<body>\n")
	html.WriteString(fmt.Sprintf("<h2>%s</h2>\n", p.sheetName))
	html.WriteString("<table>\n")

	// 添加表头和内容
	if len(rows) > 0 {
		html.WriteString("<tr>\n")
		for _, cell := range rows[0] {
			html.WriteString(fmt.Sprintf("<th>%s</th>\n", cell))
		}
		html.WriteString("</tr>\n")

		for i := 1; i < len(rows); i++ {
			html.WriteString("<tr>\n")
			for _, cell := range rows[i] {
				html.WriteString(fmt.Sprintf("<td>%s</td>\n", cell))
			}
			html.WriteString("</tr>\n")
		}
	}

	html.WriteString("</table>\n")
	html.WriteString("</body>\n</html>")

	// 写入文件
	return writeStringToFile(htmlPath, html.String())
}

// writeStringToFile 将字符串写入文件
func writeStringToFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
