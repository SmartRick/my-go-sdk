package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HTTPClient 封装HTTP客户端
type HTTPClient struct {
	client    *http.Client
	baseURL   string
	headers   map[string]string
	timeout   time.Duration
	maxRetry  int
	retryWait time.Duration
}

// NewHTTPClient 创建一个新的HTTP客户端
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client:    &http.Client{},
		headers:   make(map[string]string),
		timeout:   30 * time.Second,
		maxRetry:  3,
		retryWait: 1 * time.Second,
	}
}

// SetBaseURL 设置基础URL
func (c *HTTPClient) SetBaseURL(baseURL string) *HTTPClient {
	c.baseURL = baseURL
	return c
}

// SetTimeout 设置请求超时时间
func (c *HTTPClient) SetTimeout(timeout time.Duration) *HTTPClient {
	c.timeout = timeout
	c.client.Timeout = timeout
	return c
}

// SetHeader 设置请求头
func (c *HTTPClient) SetHeader(key, value string) *HTTPClient {
	c.headers[key] = value
	return c
}

// SetHeaders 批量设置请求头
func (c *HTTPClient) SetHeaders(headers map[string]string) *HTTPClient {
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

// SetBasicAuth 设置基本认证
func (c *HTTPClient) SetBasicAuth(username, password string) *HTTPClient {
	auth := Base64Encode([]byte(username + ":" + password))
	return c.SetHeader("Authorization", "Basic "+auth)
}

// SetBearerAuth 设置Bearer令牌认证
func (c *HTTPClient) SetBearerAuth(token string) *HTTPClient {
	return c.SetHeader("Authorization", "Bearer "+token)
}

// SetRetry 设置重试参数
func (c *HTTPClient) SetRetry(maxRetry int, retryWait time.Duration) *HTTPClient {
	c.maxRetry = maxRetry
	c.retryWait = retryWait
	return c
}

// buildURL 构建完整URL
func (c *HTTPClient) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}

	baseURL := strings.TrimSuffix(c.baseURL, "/")
	path = strings.TrimPrefix(path, "/")

	return baseURL + "/" + path
}

// applyHeaders 应用请求头
func (c *HTTPClient) applyHeaders(req *http.Request) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

// doRequest 执行HTTP请求
func (c *HTTPClient) doRequest(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(req.Context(), c.timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// 应用请求头
	c.applyHeaders(req)

	// 重试逻辑
	for i := 0; i <= c.maxRetry; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			return resp, nil
		}

		if i < c.maxRetry {
			time.Sleep(c.retryWait)
			// 克隆请求（因为req.Body可能已被消费）
			var newReq *http.Request
			if req.GetBody != nil {
				body, _ := req.GetBody()
				newReq, _ = http.NewRequestWithContext(ctx, req.Method, req.URL.String(), body)
			} else {
				newReq, _ = http.NewRequestWithContext(ctx, req.Method, req.URL.String(), nil)
			}

			// 复制请求头
			newReq.Header = req.Header
			req = newReq
		}
	}

	return nil, fmt.Errorf("请求失败，已重试%d次: %w", c.maxRetry, err)
}

// Get 发送GET请求
func (c *HTTPClient) Get(path string, params map[string]string) (*http.Response, error) {
	fullURL := c.buildURL(path)

	// 构建URL参数
	if len(params) > 0 {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, v)
		}
		fullURL = fullURL + "?" + query.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置默认的请求头
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

// Post 发送POST请求
func (c *HTTPClient) Post(path string, body interface{}) (*http.Response, error) {
	fullURL := c.buildURL(path)

	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, reqBody)
	if err != nil {
		return nil, err
	}

	// 设置默认的请求头
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

// Put 发送PUT请求
func (c *HTTPClient) Put(path string, body interface{}) (*http.Response, error) {
	fullURL := c.buildURL(path)

	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(http.MethodPut, fullURL, reqBody)
	if err != nil {
		return nil, err
	}

	// 设置默认的请求头
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

// Delete 发送DELETE请求
func (c *HTTPClient) Delete(path string) (*http.Response, error) {
	fullURL := c.buildURL(path)

	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置默认的请求头
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

// PostForm 发送表单POST请求
func (c *HTTPClient) PostForm(path string, formData map[string]string) (*http.Response, error) {
	fullURL := c.buildURL(path)

	// 构建表单数据
	data := url.Values{}
	for k, v := range formData {
		data.Add(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// 设置表单请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

// UploadFile 上传文件
func (c *HTTPClient) UploadFile(path string, fieldName, filePath string, params map[string]string) (*http.Response, error) {
	fullURL := c.buildURL(path)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件表单字段
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// 写入文件内容
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// 添加其他表单字段
	for k, v := range params {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return nil, err
	}

	// 设置multipart表单请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRequest(req)
}

// GetJSON 发送GET请求并解析JSON响应
func (c *HTTPClient) GetJSON(path string, params map[string]string, v interface{}) error {
	resp, err := c.Get(path, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleJSONResponse(resp, v)
}

// PostJSON 发送POST请求并解析JSON响应
func (c *HTTPClient) PostJSON(path string, body, v interface{}) error {
	resp, err := c.Post(path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleJSONResponse(resp, v)
}

// PutJSON 发送PUT请求并解析JSON响应
func (c *HTTPClient) PutJSON(path string, body, v interface{}) error {
	resp, err := c.Put(path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleJSONResponse(resp, v)
}

// DeleteJSON 发送DELETE请求并解析JSON响应
func (c *HTTPClient) DeleteJSON(path string, v interface{}) error {
	resp, err := c.Delete(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleJSONResponse(resp, v)
}

// handleJSONResponse 处理JSON响应
func (c *HTTPClient) handleJSONResponse(resp *http.Response, v interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，状态码：%d，响应：%s", resp.StatusCode, string(body))
	}

	if v == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// Download 下载文件
func (c *HTTPClient) Download(url, savePath string) error {
	// 创建目录
	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建保存文件
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发送GET请求
	resp, err := c.Get(url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码：%d", resp.StatusCode)
	}

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// SimpleGet 简单的GET请求
func SimpleGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// SimplePost 简单的POST请求
func SimplePost(url string, contentType string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// SimplePostJSON 简单的POST JSON请求
func SimplePostJSON(url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return SimplePost(url, "application/json", bytes.NewBuffer(jsonData))
}

// IsSuccess 检查HTTP状态码是否表示成功
func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// IsRedirect 检查HTTP状态码是否表示重定向
func IsRedirect(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// IsClientError 检查HTTP状态码是否表示客户端错误
func IsClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

// IsServerError 检查HTTP状态码是否表示服务器错误
func IsServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}
