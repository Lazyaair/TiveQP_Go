package resultverification

import (
	construction "TiveQP/Construction"
	query "TiveQP/Query"
	"bytes"
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
