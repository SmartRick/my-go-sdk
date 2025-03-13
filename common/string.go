package common

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// --------------------------------
// 字符串判断
// --------------------------------

// IsEmpty 判断字符串是否为空
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank 判断字符串是否为空白（空字符串或只包含空白字符）
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNumeric 判断字符串是否只包含数字
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsAlpha 判断字符串是否只包含字母
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsAlphaNumeric 判断字符串是否只包含字母和数字
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// ContainsAny 判断字符串是否包含任意给定的子字符串
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAll 判断字符串是否包含所有给定的子字符串
func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// StartsWith 判断字符串是否以指定前缀开始（忽略大小写）
func StartsWith(s, prefix string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
	}
	return strings.HasPrefix(s, prefix)
}

// EndsWith 判断字符串是否以指定后缀结束（忽略大小写）
func EndsWith(s, suffix string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix))
	}
	return strings.HasSuffix(s, suffix)
}

// MatchPattern 判断字符串是否匹配正则表达式
func MatchPattern(s, pattern string) (bool, error) {
	return regexp.MatchString(pattern, s)
}

// --------------------------------
// 字符串转换
// --------------------------------

// ToLowerCamel 转换为小驼峰命名（例如 helloWorld）
func ToLowerCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := splitToWords(s)
	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if words[i] != "" {
			result += strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
		}
	}

	return result
}

// ToUpperCamel 转换为大驼峰命名（例如 HelloWorld）
func ToUpperCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := splitToWords(s)
	var result string
	for _, word := range words {
		if word != "" {
			result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return result
}

// ToSnakeCase 转换为蛇形命名（例如 hello_world）
func ToSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := splitToWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "_")
}

// ToKebabCase 转换为短横线命名（例如 hello-world）
func ToKebabCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := splitToWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "-")
}

// splitToWords 将字符串拆分为单词数组
func splitToWords(s string) []string {
	// 处理已有的分隔符（下划线、短横线等）
	s = strings.NewReplacer("_", " ", "-", " ", ".", " ").Replace(s)

	// 在大写字母前添加空格
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteByte(' ')
		}
		result.WriteRune(r)
	}

	// 分割单词
	return strings.Fields(result.String())
}

// Capitalize 首字母大写
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// Uncapitalize 首字母小写
func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}

// SwapCase 大小写反转
func SwapCase(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, r := range s {
		if unicode.IsUpper(r) {
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsLower(r) {
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ReverseString 字符串反转
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// PadLeft 左填充
func PadLeft(s string, padChar rune, totalLength int) string {
	padCount := totalLength - utf8.RuneCountInString(s)
	if padCount <= 0 {
		return s
	}
	return strings.Repeat(string(padChar), padCount) + s
}

// PadRight 右填充
func PadRight(s string, padChar rune, totalLength int) string {
	padCount := totalLength - utf8.RuneCountInString(s)
	if padCount <= 0 {
		return s
	}
	return s + strings.Repeat(string(padChar), padCount)
}

// TruncateString 截断字符串并添加省略号
func TruncateString(s string, maxLength int) string {
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}

	runes := []rune(s)
	if maxLength < 3 {
		return string(runes[:maxLength])
	}

	return string(runes[:maxLength-3]) + "..."
}

// --------------------------------
// 字符串解析与格式化
// --------------------------------

// ParseInt 解析整数（包含异常处理）
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseFloat 解析浮点数（包含异常处理）
func ParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseBool 解析布尔值（包含异常处理）
func ParseBool(s string, defaultValue bool) bool {
	if s == "" {
		return defaultValue
	}

	s = strings.ToLower(s)
	switch s {
	case "true", "yes", "y", "1", "on":
		return true
	case "false", "no", "n", "0", "off":
		return false
	default:
		return defaultValue
	}
}

// FormatInt 格式化整数为千分位形式
func FormatInt(n int) string {
	inStr := strconv.Itoa(n)
	numOfDigits := len(inStr)
	if n < 0 {
		numOfDigits-- // 减去负号
	}

	if numOfDigits <= 3 {
		return inStr
	}

	var result strings.Builder
	var count int

	if n < 0 {
		result.WriteByte('-')
		inStr = inStr[1:]
	}

	for i := len(inStr) - 1; i >= 0; i-- {
		result.WriteByte(inStr[i])
		count++
		if count%3 == 0 && i > 0 {
			result.WriteByte(',')
		}
	}

	return ReverseString(result.String())
}

// --------------------------------
// 特殊目的字符串函数
// --------------------------------

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int, charset string) string {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// ToJSON 将对象转换为JSON字符串
func ToJSON(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串解析对象
func FromJSON(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

// PrettyJSON 将对象转换为格式化的JSON字符串
func PrettyJSON(obj interface{}) (string, error) {
	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// EscapeHTML 转义HTML特殊字符
func EscapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// UnescapeHTML 反转义HTML特殊字符
func UnescapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return s
}

// --------------------------------
// 字符串验证
// --------------------------------

// IsChinaPhoneNumber 验证是否为中国手机号
func IsChinaPhoneNumber(s string) bool {
	pattern := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

// IsEmail 验证是否为邮箱地址
func IsEmail(s string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

// IsIDCard 验证是否为中国身份证号（18位）
func IsIDCard(s string) bool {
	pattern := `^\d{17}[\dX]$`
	match, _ := regexp.MatchString(pattern, s)
	if !match {
		return false
	}

	// 验证校验位
	if len(s) != 18 {
		return false
	}

	// 加权因子
	factor := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验位对应值
	parity := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	sum := 0
	for i := 0; i < 17; i++ {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			return false
		}
		sum += n * factor[i]
	}

	mod := sum % 11
	if s[17] != parity[mod] && (s[17] != 'x' || parity[mod] != 'X') {
		return false
	}

	return true
}

// IsURL 验证是否为URL
func IsURL(s string) bool {
	pattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

// IsIPv4 验证是否为IPv4地址
func IsIPv4(s string) bool {
	pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

// WordCount 统计单词数量
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// CountSubstring 计算子字符串出现次数
func CountSubstring(s, sub string) int {
	return strings.Count(s, sub)
}

// JoinStrings 连接字符串数组，带分隔符
func JoinStrings(elements []string, separator string) string {
	return strings.Join(elements, separator)
}

// SplitAndTrim 分割字符串并去除空白
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// FormatTemplate 简单的字符串模板替换
func FormatTemplate(template string, data map[string]interface{}) string {
	for key, value := range data {
		placeholder := fmt.Sprintf("${%s}", key)
		valueStr := fmt.Sprintf("%v", value)
		template = strings.ReplaceAll(template, placeholder, valueStr)
	}
	return template
}

// --------------------------------
// 高级字符串处理
// --------------------------------

// ExtractBetween 提取两个标记之间的字符串
func ExtractBetween(s, start, end string) string {
	startIdx := strings.Index(s, start)
	if startIdx == -1 {
		return ""
	}

	startIdx += len(start)
	endIdx := strings.Index(s[startIdx:], end)
	if endIdx == -1 {
		return ""
	}

	return s[startIdx : startIdx+endIdx]
}

// ExtractAllBetween 提取两个标记之间的所有字符串
func ExtractAllBetween(s, start, end string) []string {
	var results []string

	for {
		startIdx := strings.Index(s, start)
		if startIdx == -1 {
			break
		}

		startIdx += len(start)
		endIdx := strings.Index(s[startIdx:], end)
		if endIdx == -1 {
			break
		}

		results = append(results, s[startIdx:startIdx+endIdx])
		s = s[startIdx+endIdx+len(end):]
	}

	return results
}

// MaskString 对字符串进行掩码处理，例如银行卡号、手机号等
func MaskString(s string, start, end int, maskChar rune) string {
	if s == "" || start >= end || start < 0 || end > len(s) {
		return s
	}

	runes := []rune(s)
	if end > len(runes) {
		end = len(runes)
	}

	for i := start; i < end; i++ {
		runes[i] = maskChar
	}

	return string(runes)
}

// MaskMiddleChars 对字符串中间部分进行掩码处理，保留前n个和后m个字符
func MaskMiddleChars(s string, keepStart, keepEnd int, maskChar rune) string {
	runes := []rune(s)
	length := len(runes)

	if length <= keepStart+keepEnd {
		return s
	}

	start := keepStart
	end := length - keepEnd

	return string(runes[:start]) + strings.Repeat(string(maskChar), end-start) + string(runes[end:])
}

// SplitByLength 按照指定长度分割字符串
func SplitByLength(s string, length int) []string {
	if length <= 0 {
		return []string{s}
	}

	var result []string
	runes := []rune(s)

	for i := 0; i < len(runes); i += length {
		end := i + length
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}

	return result
}

// RemoveAccents 移除字符串中的变音符号
func RemoveAccents(s string) string {
	// 基于Unicode标准化表单进行转换
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// SlugifyString 将字符串转换为URL友好的格式
func SlugifyString(s string) string {
	// 转换为小写
	s = strings.ToLower(s)

	// 移除变音符号
	s = RemoveAccents(s)

	// 替换非字母数字字符为连字符
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")

	// 移除开头和结尾的连字符
	s = strings.Trim(s, "-")

	return s
}

// CenterAlign 使字符串居中对齐
func CenterAlign(s string, width int, padChar rune) string {
	length := utf8.RuneCountInString(s)
	if length >= width {
		return s
	}

	leftPad := (width - length) / 2
	rightPad := width - length - leftPad

	return strings.Repeat(string(padChar), leftPad) + s + strings.Repeat(string(padChar), rightPad)
}

// WrapText 根据指定宽度进行文本换行
func WrapText(s string, width int) string {
	if width <= 0 {
		return s
	}

	var result strings.Builder
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}

		runes := []rune(line)

		for len(runes) > width {
			// 寻找合适的换行点
			wrapAt := width
			for ; wrapAt > 0 && !unicode.IsSpace(runes[wrapAt]); wrapAt-- {
			}

			if wrapAt == 0 {
				// 无法找到空格，强制在width处换行
				wrapAt = width
			} else {
				// 跳过空格
				for wrapAt > 0 && unicode.IsSpace(runes[wrapAt-1]) {
					wrapAt--
				}
			}

			result.WriteString(string(runes[:wrapAt]) + "\n")

			// 跳过开头的空格
			for wrapAt < len(runes) && unicode.IsSpace(runes[wrapAt]) {
				wrapAt++
			}

			runes = runes[wrapAt:]
		}

		result.WriteString(string(runes))
	}

	return result.String()
}

// CountWords 统计文本中的单词数量
func CountWords(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// 使用正则表达式匹配单词（考虑中英文混合情况）
	re := regexp.MustCompile(`[\p{Han}]|[a-zA-Z]+[']?[a-zA-Z]*`)
	matches := re.FindAllString(s, -1)

	return len(matches)
}

// FindLongestWord 查找文本中最长的单词
func FindLongestWord(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	longest := words[0]
	for _, word := range words[1:] {
		if len(word) > len(longest) {
			longest = word
		}
	}

	return longest
}

// FirstN 返回字符串的前N个字符
func FirstN(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return s
	}
	if n < 0 {
		n = 0
	}
	return string(runes[:n])
}

// LastN 返回字符串的后N个字符
func LastN(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return s
	}
	if n < 0 {
		n = 0
	}
	return string(runes[len(runes)-n:])
}

// RandomSubstring 随机返回字符串的一个子串
func RandomSubstring(s string, length int) string {
	runes := []rune(s)
	if len(runes) <= length {
		return s
	}

	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(len(runes) - length + 1)

	return string(runes[start : start+length])
}

// InsertAt 在指定位置插入字符串
func InsertAt(s string, index int, insert string) string {
	runes := []rune(s)
	if index < 0 {
		index = 0
	}
	if index > len(runes) {
		index = len(runes)
	}

	return string(runes[:index]) + insert + string(runes[index:])
}

// RemoveAt 删除指定位置的n个字符
func RemoveAt(s string, index, count int) string {
	runes := []rune(s)
	if index < 0 || index >= len(runes) || count <= 0 {
		return s
	}

	end := index + count
	if end > len(runes) {
		end = len(runes)
	}

	return string(runes[:index]) + string(runes[end:])
}

// CountLines 统计文本的行数
func CountLines(s string) int {
	if s == "" {
		return 0
	}

	return strings.Count(s, "\n") + 1
}

// DiffWords 比较两个字符串的差异（返回两个字符串中不同的单词）
func DiffWords(s1, s2 string) ([]string, []string) {
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	map1 := make(map[string]int)
	map2 := make(map[string]int)

	for _, word := range words1 {
		map1[word]++
	}

	for _, word := range words2 {
		map2[word]++
	}

	var onlyInS1, onlyInS2 []string

	for word, count := range map1 {
		if count2, exists := map2[word]; !exists || count > count2 {
			diff := count
			if exists {
				diff -= count2
			}
			for i := 0; i < diff; i++ {
				onlyInS1 = append(onlyInS1, word)
			}
		}
	}

	for word, count := range map2 {
		if count1, exists := map1[word]; !exists || count > count1 {
			diff := count
			if exists {
				diff -= count1
			}
			for i := 0; i < diff; i++ {
				onlyInS2 = append(onlyInS2, word)
			}
		}
	}

	return onlyInS1, onlyInS2
}
