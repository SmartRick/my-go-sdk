package common

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PathExists 检查路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir 检查路径是否是目录
func IsDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

// IsFile 检查路径是否是文件
func IsFile(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return !stat.IsDir(), nil
}

// CreateDir 创建目录，支持递归创建
func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// RemoveFile 删除文件
func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

// RemoveDir 删除目录及其内容
func RemoveDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 同步文件内容到磁盘
	return destFile.Sync()
}

// MoveFile 移动/重命名文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// ReadFile 读取文件内容为字符串
func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile 将字符串写入文件
func WriteFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// AppendToFile 追加内容到文件
func AppendToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// ListFiles 列出目录中的所有文件
func ListFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ListDirs 列出目录中的所有子目录
func ListDirs(dirPath string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != dirPath {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

// GetFileSize 获取文件大小（字节）
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileExt 获取文件扩展名
func GetFileExt(filePath string) string {
	return filepath.Ext(filePath)
}

// GetBaseName 获取文件基本名称（不包含扩展名）
func GetBaseName(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	return base[:len(base)-len(ext)]
}

// GetFileMD5 计算文件的MD5哈希值
func GetFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ReadLines 读取文件的所有行
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// FindFilesByExt 在指定目录中查找特定扩展名的文件
func FindFilesByExt(dirPath string, ext string) ([]string, error) {
	var result []string

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == strings.ToLower(ext) {
			result = append(result, path)
		}

		return nil
	})

	return result, err
}

// FileModTime 获取文件的最后修改时间
func FileModTime(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.ModTime().Unix(), nil
}

// IsFileOlderThan 检查文件是否比指定时间戳更早
func IsFileOlderThan(filePath string, timestamp int64) (bool, error) {
	modTime, err := FileModTime(filePath)
	if err != nil {
		return false, err
	}
	return modTime < timestamp, nil
}

// DirSize 递归计算目录大小
func DirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// TouchFile 创建空文件或更新现有文件的访问和修改时间
func TouchFile(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		return file.Close()
	}

	currentTime := time.Now()
	return os.Chtimes(filePath, currentTime, currentTime)
}

// FormatFileSize 格式化文件大小（转换为KB/MB/GB等）
func FormatFileSize(sizeInBytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case sizeInBytes < KB:
		return fmt.Sprintf("%d B", sizeInBytes)
	case sizeInBytes < MB:
		return fmt.Sprintf("%.2f KB", float64(sizeInBytes)/KB)
	case sizeInBytes < GB:
		return fmt.Sprintf("%.2f MB", float64(sizeInBytes)/MB)
	case sizeInBytes < TB:
		return fmt.Sprintf("%.2f GB", float64(sizeInBytes)/GB)
	default:
		return fmt.Sprintf("%.2f TB", float64(sizeInBytes)/TB)
	}
}
