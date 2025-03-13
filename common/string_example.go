package common

import (
	"fmt"
)

// StringExamples 展示字符串处理工具的使用示例
func StringExamples() {
	fmt.Println("====== 字符串处理工具使用示例 ======")

	// 基本字符串判断
	fmt.Println("\n--- 基本字符串判断 ---")
	fmt.Printf("IsEmpty(\"\") = %v\n", IsEmpty(""))
	fmt.Printf("IsEmpty(\"hello\") = %v\n", IsEmpty("hello"))
	fmt.Printf("IsBlank(\"   \") = %v\n", IsBlank("   "))
	fmt.Printf("IsNumeric(\"123\") = %v\n", IsNumeric("123"))
	fmt.Printf("IsNumeric(\"123a\") = %v\n", IsNumeric("123a"))
	fmt.Printf("IsAlpha(\"abc\") = %v\n", IsAlpha("abc"))
	fmt.Printf("IsAlpha(\"abc123\") = %v\n", IsAlpha("abc123"))
	fmt.Printf("IsAlphaNumeric(\"abc123\") = %v\n", IsAlphaNumeric("abc123"))
	fmt.Printf("ContainsAny(\"hello world\", \"world\", \"universe\") = %v\n", ContainsAny("hello world", "world", "universe"))
	fmt.Printf("ContainsAll(\"hello world\", \"hello\", \"world\") = %v\n", ContainsAll("hello world", "hello", "world"))

	// 字符串转换
	fmt.Println("\n--- 字符串转换 ---")
	fmt.Printf("ToLowerCamel(\"hello world\") = %v\n", ToLowerCamel("hello world"))
	fmt.Printf("ToUpperCamel(\"hello world\") = %v\n", ToUpperCamel("hello world"))
	fmt.Printf("ToSnakeCase(\"helloWorld\") = %v\n", ToSnakeCase("helloWorld"))
	fmt.Printf("ToKebabCase(\"HelloWorld\") = %v\n", ToKebabCase("HelloWorld"))
	fmt.Printf("Capitalize(\"hello\") = %v\n", Capitalize("hello"))
	fmt.Printf("Uncapitalize(\"Hello\") = %v\n", Uncapitalize("Hello"))
	fmt.Printf("SwapCase(\"Hello World\") = %v\n", SwapCase("Hello World"))
	fmt.Printf("ReverseString(\"Hello\") = %v\n", ReverseString("Hello"))

	// 字符串填充与截断
	fmt.Println("\n--- 字符串填充与截断 ---")
	fmt.Printf("PadLeft(\"123\", '0', 5) = %v\n", PadLeft("123", '0', 5))
	fmt.Printf("PadRight(\"123\", ' ', 5) = %v\n", PadRight("123", ' ', 5))
	fmt.Printf("TruncateString(\"Hello World\", 5) = %v\n", TruncateString("Hello World", 5))
	fmt.Printf("TruncateString(\"Hello World\", 8) = %v\n", TruncateString("Hello World", 8))

	// 解析与格式化
	fmt.Println("\n--- 解析与格式化 ---")
	fmt.Printf("ParseInt(\"123\", 0) = %v\n", ParseInt("123", 0))
	fmt.Printf("ParseInt(\"abc\", 0) = %v\n", ParseInt("abc", 0))
	fmt.Printf("ParseFloat(\"123.45\", 0) = %v\n", ParseFloat("123.45", 0))
	fmt.Printf("ParseBool(\"true\", false) = %v\n", ParseBool("true", false))
	fmt.Printf("FormatInt(1234567) = %v\n", FormatInt(1234567))

	// 高级字符串处理
	fmt.Println("\n--- 高级字符串处理 ---")
	fmt.Printf("ExtractBetween(\"<div>内容</div>\", \"<div>\", \"</div>\") = %v\n", ExtractBetween("<div>内容</div>", "<div>", "</div>"))

	html := "<div>标题1</div><p>段落1</p><div>标题2</div>"
	divContents := ExtractAllBetween(html, "<div>", "</div>")
	fmt.Printf("ExtractAllBetween - 结果数量: %d, 内容: %v\n", len(divContents), divContents)

	fmt.Printf("MaskString(\"13812345678\", 3, 7, '*') = %v\n", MaskString("13812345678", 3, 7, '*'))
	fmt.Printf("MaskMiddleChars(\"13812345678\", 3, 4, '*') = %v\n", MaskMiddleChars("13812345678", 3, 4, '*'))

	parts := SplitByLength("HelloWorld", 3)
	fmt.Printf("SplitByLength(\"HelloWorld\", 3) = %v\n", parts)

	fmt.Printf("SlugifyString(\"Hello World! 你好世界\") = %v\n", SlugifyString("Hello World! 你好世界"))
	fmt.Printf("CenterAlign(\"标题\", 10, '=') = %v\n", CenterAlign("标题", 10, '='))

	longText := "这是一段很长的文本，需要按照指定宽度进行换行处理。这样可以保证文本显示的美观性和可读性。"
	wrapped := WrapText(longText, 10)
	fmt.Println("WrapText 结果:")
	fmt.Println(wrapped)

	fmt.Printf("CountWords(\"Hello 世界，这是一个测试\") = %v\n", CountWords("Hello 世界，这是一个测试"))
	fmt.Printf("FindLongestWord(\"Hello world programming\") = %v\n", FindLongestWord("Hello world programming"))

	fmt.Printf("FirstN(\"Hello World\", 5) = %v\n", FirstN("Hello World", 5))
	fmt.Printf("LastN(\"Hello World\", 5) = %v\n", LastN("Hello World", 5))

	fmt.Printf("RandomSubstring(\"HelloWorldExample\", 5) = %v\n", RandomSubstring("HelloWorldExample", 5))
	fmt.Printf("InsertAt(\"HelloWorld\", 5, \"-\") = %v\n", InsertAt("HelloWorld", 5, "-"))
	fmt.Printf("RemoveAt(\"HelloWorld\", 5, 2) = %v\n", RemoveAt("HelloWorld", 5, 2))

	fmt.Printf("CountLines(\"Hello\\nWorld\\nTest\") = %v\n", CountLines("Hello\nWorld\nTest"))

	text1 := "The quick brown fox jumps"
	text2 := "The brown fox jumps over lazy dog"
	diff1, diff2 := DiffWords(text1, text2)
	fmt.Printf("DiffWords - 只在text1中: %v, 只在text2中: %v\n", diff1, diff2)

	fmt.Println("\n====== 字符串处理工具示例结束 ======")
}
