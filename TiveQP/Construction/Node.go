package construction

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// BitArray 表示一个基于 []uint64 的位数组，用来表示IBF
type BitArray struct {
	data []uint64
	size int
}

// NewBitArray 创建指定大小的位数组
func NewBitArray(size int) *BitArray {
	// size/64 向上取整以创建足够的 uint64 数组
	return &BitArray{
		data: make([]uint64, (size+63)/64), // +63 是为了向上取整
		size: size,
	}
}

// Set 设置位数组中的某一位
func (b *BitArray) Set(pos int, value bool) {
	if pos >= b.size || pos < 0 {
		panic("index out of bounds")
	}
	uint64Index := pos / 64 // 找到对应的 uint64 索引
	bitOffset := pos % 64   // 找到该 uint64 中的具体位位置

	if value {
		b.data[uint64Index] |= (1 << bitOffset) // 使用 OR 操作设置为 1
	} else {
		b.data[uint64Index] &^= (1 << bitOffset) // 使用 AND NOT 操作清零
	}
}

// Get 获取位数组中的某一位
func (b *BitArray) Get(pos int) bool {
	if pos >= b.size || pos < 0 {
		panic("index out of bounds")
	}
	uint64Index := pos / 64 // 找到对应的 uint64 索引
	bitOffset := pos % 64   // 找到该 uint64 中的具体位位置

	return (b.data[uint64Index] & (1 << bitOffset)) != 0
}

// DownNode 表示下层树的节点
type DownNode struct {
	Data []string
	// IBF
	IBF *BitArray
	// 补集
	LCS []string
	TCS []string
	// IBF的插入位置，HMAC==>bigint=(mod)=>index+'|'+0/1
	Bits_LCS [][]string
	Bits_TCS [][]string
	// Bits的摘要
	HV_LCS [][]byte
	HV_TCS [][]byte
	// 原始数据的摘要
	HV []byte
	// 子节点
	Left  *DownNode
	Right *DownNode
}

// UpNode 表示上层树的节点
type UpNode struct {
	Data []string
	// IBF
	IBF []byte
	// 补集
	LCS []string
	TCS []string
	YCS []string
	// IBF的插入位置，HMAC==>bigint=(mod)=>index+'|'+0/1
	Bits_LCS [][]string
	Bits_TCS [][]string
	Bits_YCS [][]string
	// Bits的摘要
	HV_LCS [][]byte
	HV_TCS [][]byte
	HV_YCS [][]byte
	// 子节点的合并摘要
	HV []byte
	// 子节点
	Left  *UpNode
	Right *UpNode
}

// 初始化所需的各个方法

// bits-HAMC计算
func HMACSHA256(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// HV计算
func HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
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
func OrIBF(a, b *BitArray) *BitArray {
	if a.size != b.size {
		panic("bit arrays must be the same size for OR operation")
	}

	// 创建结果位数组
	result := NewBitArray(a.size)
	for i := 0; i < len(a.data); i++ {
		result.data[i] = a.data[i] | b.data[i] // 对每个 uint64 执行 OR 操作
	}

	return result
}
