package construction

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"math/big"
	"strconv"
)

// 初始化所需的各个方法
// bits-HAMC计算
func HMACSHA256(message, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

// HV计算
func HashSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// 集合的合并去重
func MergeSet(slice1, slice2 []string) []string {
	// 让 slice1 是较小的切片
	if len(slice1) > len(slice2) {
		slice1, slice2 = slice2, slice1
	}

	// 使用较大的切片初始化 map
	// struct{}不占用额外的内存空间
	uniqueMap := make(map[string]struct{}, len(slice2))
	for _, item := range slice2 {
		uniqueMap[item] = struct{}{}
	}

	// 只添加 slice1 中 map 中不存在的元素
	for _, item := range slice1 {
		uniqueMap[item] = struct{}{}
	}

	// 提取结果
	result := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		result = append(result, key)
	}

	return result
}

// Or 执行两个位数组的按位或运算，返回结果位数组
func OrIBF(a, b *TwinBitArray) *TwinBitArray {
	// 检查两者的列数是否相同
	if a.cols != b.cols {
		panic("列数不一致！")
	}

	// 创建一个新的 TwinBitArray 作为结果
	result := NewTwinBitArray(a.cols)

	// 对两行数据执行 OR 操作
	for i := 0; i < 2; i++ { // 两行
		for j := 0; j < len(a.data[i]); j++ { // 每行的 uint64 数组
			result.data[i][j] = a.data[i][j] | b.data[i][j] // 对应的 uint64 执行 OR 操作
		}
	}

	return result
}

// 对 MS 的插入
func Insert(twinlist *TwinBitArray, data string, keylist []string, rb int) error {
	// 循环计算每个 key 对应的位置
	for i := 0; i < len(keylist)-1; i++ {
		// 计算 HMAC(w, k_i)
		outbytes := HMACSHA256([]byte(data), []byte(keylist[i]))
		bi := new(big.Int).SetBytes(outbytes)
		twinIndex := bi.Mod(bi, big.NewInt(int64(twinlist.cols))).Int64() // twin_id

		// 计算 (h_k+1)
		hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
		hkp1bi := new(big.Int).SetBytes(hkp1)

		// 和随机数 XOR
		sha1bytes := sha256.Sum256(hkp1bi.Xor(hkp1bi, big.NewInt(int64(rb))).Bytes())
		location := new(big.Int).SetBytes(sha1bytes[:]).Mod(new(big.Int).SetBytes(sha1bytes[:]), big.NewInt(2)).Int64()

		// 置位基于位置 (0 or 1)
		if location == 0 {
			twinlist.Set(0, int(twinIndex), true)  // Set bit to 1 for twinlist[0][twinIndex]
			twinlist.Set(1, int(twinIndex), false) // Set bit to 0 for twinlist[1][twinIndex]
		} else {
			twinlist.Set(1, int(twinIndex), true)  // Set bit to 1 for twinlist[1][twinIndex]
			twinlist.Set(0, int(twinIndex), false) // Set bit to 0 for twinlist[0][twinIndex]
		}
	}
	return nil
}

// 对补集的处理
func InsertCS(twinlist *TwinBitArray, data string, bit_CS_i, keylist []string, hv_cs []byte, rb int) error {
	// 循环计算每个 key 对应的位置
	for i := 0; i < len(keylist)-1; i++ {
		// 计算 HMAC(w, k_i)
		outbytes := HMACSHA256([]byte(data), []byte(keylist[i]))
		bi := new(big.Int).SetBytes(outbytes)
		twinIndex := bi.Mod(bi, big.NewInt(int64(twinlist.cols))).Int64() // twin_id

		// 计算 (h_k+1)
		hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
		hkp1bi := new(big.Int).SetBytes(hkp1)

		// 和随机数 XOR
		sha1bytes := sha256.Sum256(hkp1bi.Xor(hkp1bi, big.NewInt(int64(rb))).Bytes())
		location := new(big.Int).SetBytes(sha1bytes[:]).Mod(new(big.Int).SetBytes(sha1bytes[:]), big.NewInt(2)).Int64()

		// 置位基于位置 (0 or 1)
		// 就是这里（！！！！高学长未置位！！！！）
		if location == 0 {
			// twinlist.Set(0, int(twinIndex), true)  // Set bit to 1 for twinlist[0][twinIndex]
			// twinlist.Set(1, int(twinIndex), false) // Set bit to 0 for twinlist[1][twinIndex]
			bit_CS_i = append(bit_CS_i, strconv.FormatInt(twinIndex, 10)+"|"+strconv.Itoa(1))
			hv_cs = append(hv_cs, HMACSHA256([]byte(bit_CS_i[i]), []byte(keylist[i]))...)
		} else {
			// twinlist.Set(1, int(twinIndex), true)  // Set bit to 1 for twinlist[1][twinIndex]
			// twinlist.Set(0, int(twinIndex), false) // Set bit to 0 for twinlist[0][twinIndex]
			bit_CS_i = append(bit_CS_i, strconv.FormatInt(twinIndex, 10)+"|"+strconv.Itoa(0))
			hv_cs = append(hv_cs, HMACSHA256([]byte(bit_CS_i[i]), []byte(keylist[i]))...)
		}
	}
	return nil
}

// AES 加密
// AES加密函数（CBC模式）
func Encrypt(plaintext, key []byte) ([]byte, error) {
	// 确保密钥长度为16、24或32字节（分别对应AES-128, AES-192, AES-256）
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个随机的初始化向量（IV），长度为AES的块大小（16字节）
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 创建一个AES加密模式（CBC模式）
	mode := cipher.NewCBCEncrypter(block, iv)

	// 对明文进行填充，使其长度为16的倍数
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	// 填充后的明文长度为 aes.BlockSize 的倍数
	paddedPlaintext := append(plaintext, byte(padding))

	// 使用CBC模式加密明文
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedPlaintext)

	// 返回加密后的密文
	return ciphertext, nil
}
