package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

// --------------------------------
// Base64 编解码
// --------------------------------

// Base64Encode 将字节数组编码为base64字符串
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 将base64字符串解码为字节数组
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// Base64UrlEncode 将字节数组编码为URL安全的base64字符串
func Base64UrlEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64UrlDecode 将URL安全的base64字符串解码为字节数组
func Base64UrlDecode(encoded string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(encoded)
}

// --------------------------------
// URL 编解码
// --------------------------------

// UrlEncode 对字符串进行URL编码
func UrlEncode(data string) string {
	return url.QueryEscape(data)
}

// UrlDecode 对URL编码的字符串进行解码
func UrlDecode(encoded string) (string, error) {
	return url.QueryUnescape(encoded)
}

// --------------------------------
// 哈希函数
// --------------------------------

// MD5Hash 计算字符串的MD5哈希值
func MD5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SHA1Hash 计算字符串的SHA1哈希值
func SHA1Hash(data string) string {
	hash := sha1.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SHA256Hash 计算字符串的SHA256哈希值
func SHA256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// --------------------------------
// 密码哈希和验证
// --------------------------------

// HashPassword 使用bcrypt对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 验证密码与哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --------------------------------
// AES 加解密
// --------------------------------

// AESEncrypt 使用AES-GCM算法加密数据
func AESEncrypt(plainText []byte, key []byte) ([]byte, error) {
	// 确保密钥长度为16, 24或32字节(128, 192或256位)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("密钥长度必须为16, 24或32字节")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// GCM模式不需要填充
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机的nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据，nonce附加在密文前面
	cipherText := aesgcm.Seal(nonce, nonce, plainText, nil)
	return cipherText, nil
}

// AESDecrypt 使用AES-GCM算法解密数据
func AESDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	// 确保密钥长度为16, 24或32字节(128, 192或256位)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("密钥长度必须为16, 24或32字节")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 从密文中提取nonce
	nonceSize := aesgcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("密文长度不足")
	}

	nonce, cipherTextData := cipherText[:nonceSize], cipherText[nonceSize:]

	// 解密数据
	plainText, err := aesgcm.Open(nil, nonce, cipherTextData, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// GenerateAESKey 生成指定位数的AES密钥
func GenerateAESKey(bits int) ([]byte, error) {
	// AES密钥长度必须是128, 192或256位
	if bits != 128 && bits != 192 && bits != 256 {
		return nil, errors.New("AES密钥位数必须为128, 192或256")
	}

	bytes := bits / 8
	key := make([]byte, bytes)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// --------------------------------
// RSA 加解密
// --------------------------------

// GenerateRSAKeyPair 生成RSA公钥和私钥对
func GenerateRSAKeyPair(bits int) (string, string, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 将私钥编码为PEM格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	// 从私钥中提取公钥并编码为PEM格式
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	return string(publicKeyPEM), string(privateKeyPEM), nil
}

// RSAEncrypt 使用RSA公钥加密数据
func RSAEncrypt(plainText []byte, publicKeyPEM string) ([]byte, error) {
	// 解析PEM格式的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("无效的公钥")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("不是有效的RSA公钥")
	}

	// 加密数据
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, plainText)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

// RSADecrypt 使用RSA私钥解密数据
func RSADecrypt(cipherText []byte, privateKeyPEM string) ([]byte, error) {
	// 解析PEM格式的私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("无效的私钥")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 解密数据
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// --------------------------------
// 辅助函数
// --------------------------------

// GenerateRandomBytes 生成指定数量的随机字节
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateSecureToken 生成安全的随机令牌
func GenerateSecureToken(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return Base64UrlEncode(bytes), nil
}

// EncryptString 使用AES加密字符串并返回Base64编码的结果
func EncryptString(plainText string, key []byte) (string, error) {
	cipherText, err := AESEncrypt([]byte(plainText), key)
	if err != nil {
		return "", err
	}
	return Base64Encode(cipherText), nil
}

// DecryptString 解密Base64编码的AES加密字符串
func DecryptString(encryptedText string, key []byte) (string, error) {
	cipherText, err := Base64Decode(encryptedText)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	plainText, err := AESDecrypt(cipherText, key)
	if err != nil {
		return "", fmt.Errorf("AES解密失败: %w", err)
	}

	return string(plainText), nil
}
