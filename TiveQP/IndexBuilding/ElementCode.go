package indexbuilding

import (
	"fmt"
	"math"
	"strings"
)

// Prefix 返回一个数值x的bitsize位前缀族
func Prefix(bitsize int, x int) ([]string, error) {
	// 使用位移运算计算
	if (1<<bitsize)-1 < x {
		return nil, fmt.Errorf("x超出了bitsize位能表示的范围")
	}

	// 将 x 转换为二进制字符串，长度为 bitsize
	/*
		%b：这是表示整数的二进制表示。它会将 x 以二进制形式输出。
		%0b：前面加上 0 表示当输出的位数不够时，在前面用 0 来补齐。
		%0*b：前面加上一个星号 *，意味着这个宽度（位数）由函数参数提供。这个宽度就是我们传入的 bitsize。
	*/
	result := fmt.Sprintf("%0*b", bitsize, x)

	// 必然是bitsize+1个
	list := make([]string, bitsize+1)

	// 生成前缀族
	for i := 0; i < bitsize; i++ {
		list[i] = result[:bitsize-i] + strings.Repeat("*", i)
	}
	// 最后一个全是 '*'
	list[bitsize] = strings.Repeat("*", bitsize)
	return list, nil
}

// Samesize 返回相同前缀的长度
// 使用64位块的方式来进行加速
// 废弃没什么用
/*
func Samesize(s1, s2 string) (int, error) {
	i := 0
	length := len(s1)
	// 将字符串转换为字节切片
	b1 := []byte(s1)
	b2 := []byte(s2)

	// 比较64位块（8字节）
	for i+8 <= length {
		// 将字符串的当前字节块转换为64位整数
		if *(*uint64)(unsafe.Pointer(&b1[i])) != *(*uint64)(unsafe.Pointer(&b2[i])) {
			break
		}
		i += 8
	}

	// 比较剩下的字符（逐字符比较）
	for i < length && s1[i] == s2[i] {
		i++
	}

	return i, nil
}
*/

// Range 返回数值区间[down, up]的bitsize范围
/*
树结构
第一步：从down到up依次生成bitsize位的2进制字符串（不够bitsize位的左边补0），放入strs中；
第二步：i=[0,log(up-down+1)/log2],检查strs[0]和strs[1]的从右往左数i+1位是否相同。
	是则合并strs[j]和strs[j+1],合并规则为从右往左数第i位置为“*“；
	否则则从strs[1]和strs[2]开始合并（合并规则不变），并将str[0]放入result中。
	不论是否相同，若strs中合并剩下1个，则将其加入result中。
*/
func Range(bitsize, down, up int) ([]string, error) {
	// 结果result
	result := make([]string, 0, 2*(int(math.Log2(float64(up-down+1)))+1))
	// 第一步：生成从down到up的bitsize位二进制字符串
	// 临时strs
	strs := make([]string, 0, up-down+1)
	for i := down; i <= up; i++ {
		// 使用%0*b格式将数字转换为bitsize位的二进制字符串
		str := fmt.Sprintf("%0*b", bitsize, i)
		strs = append(strs, str)
	}

	// 第二步：开始合并操作
	iMax := int(math.Log2(float64(up-down+1)) + 1) // log2(up-down+1)
	for i := 0; i <= iMax; i++ {
		// 若只剩一个字符串，则加入result并返回
		if len(strs) == 1 {
			result = append(result, strs[0])
			return result, nil
		}
		// 第一项为在result中，添加、删除
		if strs[0][bitsize-i-2] != strs[1][bitsize-i-2] {
			result = append(result, strs[0])
			strs = strs[1:]
		}
		// 最后一项是否在result中，添加、删除
		if len(strs)%2 != 0 {
			result = append(result, strs[len(strs)-1])
			strs = strs[:len(strs)-1]
		}
		temp := make([]string, 0, len(strs)/2)
		for j := 0; j < len(strs)-1; j = j + 2 {
			// 将字符串转换为rune切片，处理多字节字符
			runes := []rune(strs[j])
			//只用修改第bitsize-i-1位
			runes[bitsize-i-1] = '*'
			temp = append(temp, string(runes))
		}
		strs = temp
	}
	return result, nil
}
