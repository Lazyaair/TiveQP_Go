package trapdoor

import (
	indexbuilding "TiveQP/IndexBuilding"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func HMACSHA256(message, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

func HashSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func Encrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	padding := blockSize - len(plainText)%blockSize
	paddedText := append(plainText, bytes.Repeat([]byte{byte(padding)}, padding)...)
	ciphertext := make([]byte, len(paddedText))
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)
	return append(iv, ciphertext...), nil
}

func Decrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := cipherText[:blockSize]
	ciphertext := cipherText[blockSize:]
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	padding := int(plaintext[len(plaintext)-1])
	if padding > len(plaintext) || padding > blockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padding]
	return plaintext, nil
}

// ParseUser 解析每一行并创建 Owner 对象
func ParseUser(line string) (*indexbuilding.User, error) {
	// 使用 "**" 分割字段
	fields := strings.Split(line, "**")
	if len(fields) != 6 {
		return nil, fmt.Errorf("数据格式错误，期望 6 个字段，但得到 %d 个字段", len(fields))
	}

	// 解析字段并转换数据类型
	lat, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("解析 Lat 出错: %v", err)
	}

	lng, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("解析 Lng 出错: %v", err)
	}

	hourStart, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("解析 HourStart 出错: %v", err)
	}

	minStart, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, fmt.Errorf("解析 MinStart 出错: %v", err)
	}

	// 创建 User 对象并返回
	user := &indexbuilding.User{
		Type:      fields[0],
		City:      fields[1],
		Lat:       lat,
		Lng:       lng,
		HourStart: hourStart,
		MinStart:  minStart,
	}

	return user, nil
}
