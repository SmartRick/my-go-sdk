# Common工具包

本目录包含各种通用工具函数，包括文件处理、加解密、HTTP请求和并发控制等实用功能。

## 1. 文件工具 (file.go)

`file.go` 提供了丰富的文件和目录操作函数，简化常见的文件处理任务。

### 文件检测函数

```go
// 检查路径是否存在
exists, err := PathExists("/path/to/file")

// 检查路径是否是目录
isDir, err := IsDir("/path/to/dir")

// 检查路径是否是文件
isFile, err := IsFile("/path/to/file")
```

### 文件创建和删除

```go
// 创建目录（支持递归创建）
err := CreateDir("/path/to/new/dir")

// 删除文件
err := RemoveFile("/path/to/file")

// 删除目录及其内容
err := RemoveDir("/path/to/dir")

// 创建空文件或更新已有文件的访问和修改时间
err := TouchFile("/path/to/file")
```

### 文件复制和移动

```go
// 复制文件
err := CopyFile("/path/source", "/path/destination")

// 移动/重命名文件
err := MoveFile("/path/source", "/path/destination")
```

### 文件内容读写

```go
// 读取文件内容为字符串
content, err := ReadFile("/path/to/file")

// 将字符串写入文件（覆盖）
err := WriteFile("/path/to/file", "文件内容")

// 追加内容到文件
err := AppendToFile("/path/to/file", "追加的内容")

// 读取文件的所有行
lines, err := ReadLines("/path/to/file")
```

### 文件遍历

```go
// 列出目录中的所有文件（递归）
files, err := ListFiles("/path/to/dir")

// 列出目录中的所有子目录（递归）
dirs, err := ListDirs("/path/to/dir")

// 在指定目录中查找特定扩展名的文件
jpgFiles, err := FindFilesByExt("/path/to/dir", ".jpg")
```

### 文件信息获取

```go
// 获取文件大小（字节）
size, err := GetFileSize("/path/to/file")

// 格式化文件大小（转换为KB/MB/GB等）
formattedSize := FormatFileSize(sizeInBytes)

// 获取文件扩展名
ext := GetFileExt("/path/to/file.txt") // 返回 ".txt"

// 获取文件基本名称（不包含扩展名）
name := GetBaseName("/path/to/file.txt") // 返回 "file"

// 计算文件的MD5哈希值
md5hash, err := GetFileMD5("/path/to/file")

// 获取文件的最后修改时间（Unix时间戳）
modTime, err := FileModTime("/path/to/file")

// 检查文件是否比指定时间戳更早
isOlder, err := IsFileOlderThan("/path/to/file", timestamp)

// 递归计算目录大小
dirSize, err := DirSize("/path/to/dir")
```

## 2. 加解密工具 (crypto.go)

`crypto.go` 提供了常用的加密、解密和编码功能，简化数据安全处理。

### Base64 编解码

```go
// 标准Base64编码
encodedStr := Base64Encode([]byte("Hello World"))

// 标准Base64解码
decodedBytes, err := Base64Decode(encodedStr)

// URL安全的Base64编码
urlSafeStr := Base64UrlEncode([]byte("Hello World"))

// URL安全的Base64解码
decodedBytes, err := Base64UrlDecode(urlSafeStr)
```

### URL 编解码

```go
// URL编码
encoded := UrlEncode("Hello World & More")

// URL解码
decoded, err := UrlDecode(encoded)
```

### 哈希计算

```go
// 计算MD5哈希
md5hash := MD5Hash("Hello World")

// 计算SHA1哈希
sha1hash := SHA1Hash("Hello World")

// 计算SHA256哈希
sha256hash := SHA256Hash("Hello World")
```

### 密码处理

```go
// 对密码进行哈希处理
hashedPassword, err := HashPassword("my-secure-password")

// 验证密码与哈希值是否匹配
isMatch := CheckPasswordHash("my-secure-password", hashedPassword)
```

### AES 加解密

```go
// 生成AES密钥 (128, 192 或 256 位)
key, err := GenerateAESKey(256)

// AES加密
cipherText, err := AESEncrypt([]byte("明文数据"), key)

// AES解密
plainText, err := AESDecrypt(cipherText, key)

// 字符串加密（AES加密+Base64编码）
encryptedStr, err := EncryptString("明文字符串", key)

// 字符串解密
decryptedStr, err := DecryptString(encryptedStr, key)
```

### RSA 加解密

```go
// 生成RSA密钥对
publicKey, privateKey, err := GenerateRSAKeyPair(2048)

// RSA公钥加密
cipherText, err := RSAEncrypt([]byte("明文数据"), publicKey)

// RSA私钥解密
plainText, err := RSADecrypt(cipherText, privateKey)
```

### 随机数生成

```go
// 生成随机字节
randomBytes, err := GenerateRandomBytes(16)

// 生成随机字符串
randomString, err := GenerateRandomString(12)

// 生成安全的随机令牌
token, err := GenerateSecureToken(32)
```

## 3. HTTP请求工具 (http.go)

`http.go` 提供了简洁易用的HTTP客户端，支持各种HTTP请求方式和高级功能。

### 基本用法

```go
// 创建HTTP客户端
client := NewHTTPClient()

// 设置基础URL
client.SetBaseURL("https://api.example.com")

// 设置超时
client.SetTimeout(30 * time.Second)

// 设置请求头
client.SetHeader("User-Agent", "MyApp/1.0")

// 设置认证
client.SetBasicAuth("username", "password")
// 或
client.SetBearerAuth("your-token")

// 设置重试
client.SetRetry(3, 1 * time.Second)
```

### 发送请求

```go
// GET请求
resp, err := client.Get("/users", map[string]string{"page": "1"})

// POST请求 (JSON)
resp, err := client.Post("/users", map[string]interface{}{
    "name": "张三",
    "age": 30,
})

// PUT请求
resp, err := client.Put("/users/1", userData)

// DELETE请求
resp, err := client.Delete("/users/1")

// 表单POST请求
resp, err := client.PostForm("/login", map[string]string{
    "username": "user",
    "password": "pass",
})

// 上传文件
resp, err := client.UploadFile(
    "/upload", 
    "file", 
    "./document.pdf",
    map[string]string{"description": "My Document"}
)
```

### JSON处理

```go
// GET请求并解析JSON响应
var users []User
err := client.GetJSON("/users", nil, &users)

// POST请求并解析JSON响应
var user User
err := client.PostJSON("/users", userData, &user)
```

### 文件下载

```go
// 下载文件
err := client.Download("https://example.com/file.zip", "./downloads/file.zip")
```

### 简单请求函数

```go
// 简单GET请求
data, err := SimpleGet("https://example.com/api/data")

// 简单POST请求 (JSON)
data, err := SimplePostJSON("https://example.com/api/users", userData)
```

## 4. 并发线程工具 (concurrency.go)

`concurrency.go` 提供了并发编程相关的工具，简化多线程任务处理。

### 线程安全的数据结构

```go
// 线程安全的计数器
counter := SafeCounter{}
counter.Increment()
value := counter.Get()
counter.Decrement()

// 线程安全的映射
safeMap := NewSafeMap()
safeMap.Set("key", "value")
value := safeMap.Get("key")
safeMap.Delete("key")
```

### 线程池

```go
// 创建线程池 (工作线程数=5, 队列大小=100)
pool := NewThreadPool(5, 100)

// 提交任务
pool.Submit(func() (interface{}, error) {
    // 任务逻辑
    return nil, nil
})

// 等待所有任务完成
pool.Wait()

// 停止线程池
pool.Stop()
```

### 并行执行器

```go
// 创建并行执行器
parallelizer := NewParallelizer(10, 30 * time.Second)

// 准备任务
tasks := []Task{
    func() (interface{}, error) { return "任务1结果", nil },
    func() (interface{}, error) { return "任务2结果", nil },
}

// 并行执行任务
results := parallelizer.Run(tasks)

// 处理结果
for _, result := range results {
    if result.Error != nil {
        fmt.Println("错误:", result.Error)
    } else {
        fmt.Println("结果:", result.Value)
    }
}
```

### 简便的并发执行函数

```go
// 使用超时并行执行任务
results := RunTasksWithTimeout(10 * time.Second, task1, task2, task3)

// 限制并发数量执行任务
results := RunTasksConcurrently(5, task1, task2, task3, task4, task5)
```

### 信号量

```go
// 创建一个大小为3的信号量
sem := NewSemaphore(3)

// 获取信号量
sem.Acquire()

// 带超时获取信号量
success := sem.AcquireWithTimeout(5 * time.Second)

// 尝试获取信号量（不阻塞）
if sem.TryAcquire() {
    // 成功获取信号量
}

// 释放信号量
sem.Release()
```

### 速率限制器

```go
// 创建每秒10个请求的速率限制器
limiter := NewRateLimiter(10)

// 等待获取令牌
limiter.Wait()

// 带速率限制执行函数
RunWithRateLimit(10, func() {
    // 限制频率的操作
})
```

### 批量处理

```go
// 批量处理数据
items := []string{"item1", "item2", "item3", "item4", "item5"}
err := BatchProcess(items, 2, 3, func(batch []string) error {
    // 处理一批数据
    return nil
})
```

## 文件工具使用示例

```go
package main

import (
	"fmt"
	"github.com/SmartRick/my-go-sdk/common"
	"time"
)

func main() {
	// 创建目录
	err := common.CreateDir("./test_dir")
	if err != nil {
		fmt.Println("创建目录失败:", err)
		return
	}
	
	// 写入文件
	err = common.WriteFile("./test_dir/test.txt", "Hello, World!")
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}
	
	// 检查文件是否存在
	exists, _ := common.PathExists("./test_dir/test.txt")
	fmt.Println("文件存在:", exists)
	
	// 获取文件信息
	size, _ := common.GetFileSize("./test_dir/test.txt")
	fmt.Println("文件大小:", common.FormatFileSize(size))
	
	// 计算MD5
	md5, _ := common.GetFileMD5("./test_dir/test.txt")
	fmt.Println("文件MD5:", md5)
	
	// 复制文件
	err = common.CopyFile("./test_dir/test.txt", "./test_dir/test_copy.txt")
	if err != nil {
		fmt.Println("复制文件失败:", err)
	}
	
	// 列出目录文件
	files, _ := common.ListFiles("./test_dir")
	fmt.Println("目录中的文件:")
	for _, file := range files {
		fmt.Println("- ", file)
	}
	
	// 追加内容
	common.AppendToFile("./test_dir/test.txt", "\n这是追加的内容")
	
	// 读取文件内容
	content, _ := common.ReadFile("./test_dir/test.txt")
	fmt.Println("更新后的文件内容:", content)
	
	// 清理测试文件
	time.Sleep(2 * time.Second) // 等待一下，以便看到测试结果
	common.RemoveDir("./test_dir")
	fmt.Println("测试完成，清理测试文件")
}
```

## 未来计划

未来将添加更多通用工具功能，例如：
- 日志工具：提供灵活的日志记录功能
- 配置管理：支持多种配置格式的读取和管理
- 时间工具：时间格式化、时区转换、计时器等
- 字符串处理：字符串验证、转换、格式化等
- 数据库工具：数据库连接池管理、简化的CRUD操作 