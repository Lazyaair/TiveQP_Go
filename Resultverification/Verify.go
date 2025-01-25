package resultverification

import (
	config "TiveQP/Config"
	construction "TiveQP/Construction"
	query "TiveQP/Query"
	trapdoor "TiveQP/Trapdoor"
	"bytes"
	"strconv"
)

// 传入root节点hash值，与poof切片，左→右，检测最后两个的height是否相同，相同则合并生成上层节点hash，直至height==0
func CheckHV(rootH []byte, pi []*query.PON) bool {
	if len(pi) == 0 {
		return false
	}
	n_list := make([]*query.PON, 0, len(pi)/2)
	index := 0
	for {
		lenght := len(n_list)
		if !sameHeight(n_list, lenght) {
			if lenght == 1 && n_list[0].Height == 0 {
				if bytes.Equal(rootH, n_list[0].HV) {
					return true
				} else {
					return false
				}
			} else {
				n_list = append(n_list, pi[index])
				index += 1
				continue
			}
		} else {
			midNode := &query.PON{
				HV:     construction.HashSHA256(append(n_list[lenght-2].HV, n_list[lenght-1].HV...)),
				Height: n_list[lenght-1].Height - 1,
				Typ:    "MidNode",
			}
			n_list = n_list[:lenght-2]
			n_list = append(n_list, midNode)
		}
	}
}

func sameHeight(n_list []*query.PON, lenght int) bool {
	if lenght <= 1 {
		return false
	} else {
		if n_list[lenght-1].Height == n_list[lenght-2].Height {
			return true
		} else {
			return false
		}
	}
}

func CheckCompleteness(td *trapdoor.T, pi []*query.PON) bool {
	tc := TCompute(td)
	for i, v := range pi {
		if pi[i].Typ == "UMN" {
			// 检查YCS是否为nil
			if v.Bits_YCS != nil && v.HV_YCS != nil {
				// 如果YCS不为nil，只检查YCS完备性
				if !CompareArrays(v.Bits_YCS, tc.Bits_YCS, v.HV_YCS, tc.HV_YCS) {
					return false
				}
			} else {
				// 如果YCS为nil，检查LCS和TCS完备性
				if !CompareArrays(v.Bits_LCS, tc.Bits_LCS, v.HV_LCS, tc.HV_LCS) &&
					!CompareArrays(v.Bits_TCS, tc.Bits_TCS, v.HV_TCS, tc.HV_TCS) {
					return false
				}
			}
		}
	}
	return true
}

func CompareArrays(A1, A2 [][]string, B1, B2 [][][]byte) bool {
	// 检查维度是否匹配
	if len(A1) != len(B1) || len(A2) != len(B2) {
		return false
	}

	// 使用map存储A1中的字符串
	strMap := make(map[string]struct{})
	for _, arr := range A1 {
		for _, s := range arr {
			strMap[s] = struct{}{}
		}
	}

	// 检查A2中是否有匹配的字符串
	foundStringMatch := false
	for _, arr := range A2 {
		for _, s := range arr {
			if _, exists := strMap[s]; exists {
				foundStringMatch = true
				break
			}
		}
		if foundStringMatch {
			break
		}
	}

	// 如果没有找到字符串匹配，直接返回false
	if !foundStringMatch {
		return false
	}

	strMap = nil
	// 使用map存储B1中的字节数组
	// 使用string作为key因为[]byte不能作为map的key
	byteMap := make(map[string]struct{})
	for _, arr := range B1 {
		for _, b := range arr {
			byteMap[string(b)] = struct{}{}
		}
	}

	// 检查B2中是否有匹配的字节数组
	for _, arr := range B2 {
		for _, b := range arr {
			if _, exists := byteMap[string(b)]; exists {
				return true // 已经找到了字符串匹配和字节数组匹配
			}
		}
	}

	return false
}

type TfComplte struct {
	Bits_LCS [][]string
	Bits_TCS [][]string
	Bits_YCS [][]string

	HV_LCS [][][]byte
	HV_TCS [][][]byte
	HV_YCS [][][]byte
}

func TCompute(td *trapdoor.T) *TfComplte {
	// 提取通用的处理函数
	process := func(data [][]string, keylist []string) ([][]string, [][][]byte) {
		bits := make([][]string, len(data))
		hv := make([][][]byte, len(data))
		for i := 0; i < len(data); i++ {
			bits[i] = make([]string, len(keylist)-1)
			hv[i] = make([][]byte, len(keylist)-1)
			for j := 0; j < len(keylist)-1; j++ {
				row, col := query.ParseIndex(data[i][j], config.IbfLength, config.Rb)
				bitValue := strconv.FormatInt(int64(col), 10) + "|" + strconv.Itoa(row)
				bits[i][j] = bitValue
				hv[i][j] = append(hv[i][j], construction.HMACSHA256([]byte(bitValue), []byte(keylist[j]))...)
			}
		}
		return bits, hv
	}

	// 使用通用的处理函数处理三组数据
	bits_YCS, hv_YCS := process(td.T1, config.Keylist)
	bits_LCS, hv_LCS := process(td.T2, config.Keylist)
	bits_TCS, hv_TCS := process(td.T3, config.Keylist)

	// 返回最终结果
	compute := &TfComplte{
		Bits_YCS: bits_YCS,
		Bits_LCS: bits_LCS,
		Bits_TCS: bits_TCS,
		HV_YCS:   hv_YCS,
		HV_LCS:   hv_LCS,
		HV_TCS:   hv_TCS,
	}
	return compute
}
