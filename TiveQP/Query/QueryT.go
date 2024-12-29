package query

// E:\Github\TiveQP\TiveQP\TiveQP\Cachemidnode 0.txt
// k = 10
// eve = Restaurants**ATLANTA**33.846335**-84.3635778**7**0**21**0
// Query start
// Query end
// Verification :
// Result number：10
// Correctness :
// Restaurants**ATLANTA**33.846335**-84.3635778**7**0**21**0
// Restaurants**ATLANTA**33.8428749**-84.3785385**10**0**17**0
// Restaurants**ATLANTA**33.8275823**-84.328604**6**0**22**0
// Restaurants**ATLANTA**33.84269152**-84.37045581**4**0**23**0
// Restaurants**ATLANTA**33.847331**-84.372713**11**0**21**0
// Restaurants**ATLANTA**33.84703603**-84.36566228**11**30**21**0
// Restaurants**ATLANTA**33.820152**-84.387432**11**30**21**0
// Restaurants**ATLANTA**33.8475491**-84.3737958**1**0**20**0
// Restaurants**ATLANTA**33.8167245**-84.3356715**10**0**15**30
// Restaurants**ATLANTA**33.89221993**-84.32607651**10**0**22**0
// true
// true
// Completeness :
// true
// time_query = 577017
// time_proof = 5
// Verification time: 359623
// proof size ：
// 13948
import (
	construction "TiveQP/Construction"
	trapdoor "TiveQP/TrapDoor"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

func QueryT(root *construction.Node, td *trapdoor.T, k *int, height int, rb int, result *[]*construction.Node, pi *[]*PON) {
	// root.Print()
	// 空节点
	if root == nil {
		return
	}
	// k 个查询结束 UNN
	if *k == 0 {
		poof := &PON{
			HV:     root.HV,
			Height: height,
			Typ:    "UNN",
		}
		*pi = append(*pi, poof)
		fmt.Println("UNN height=", height)
		return
	}
	// UMN
	if !check(root, td, rb) {
		poof := &PON{
			HV:     root.HV,
			Height: height,
			Typ:    "UMN",
		}
		*pi = append(*pi, poof)
		fmt.Println("UMN height=", height)
		return
	}
	// leafNode == MLN
	if root.Left == nil && root.Right == nil {
		fmt.Println("k=", *k, "height=", height, "<=================================================")
		//root.Print()
		*k -= 1
		*result = append(*result, root)
		poof := &PON{
			HV:     root.HV,
			Height: height,
			Typ:    "MLN",
		}
		*pi = append(*pi, poof)
		return
	}

	if root.Left != nil {
		QueryT(root.Left, td, k, height+1, rb, result, pi)
	}
	if root.Right != nil {
		QueryT(root.Right, td, k, height+1, rb, result, pi)
	}

}

func check(root *construction.Node, t *trapdoor.T, rb int) bool {
	if root.YCS != nil {
		// 只关 T1
		for _, rowVal := range t.T1 {
			for _, val := range rowVal {
				row, col := ParseIndex(val, root.IBF.Cols, rb)
				if !root.IBF.Get(row, col) {
					return false
				}
			}
			// if count == len(rowVal) {
			// 	return true
			// }
		}
		return true
	} else {
		// 只关 T2+T3
		flag := false
		count := 0
		for _, rowVal := range t.T2 {
			for _, val := range rowVal {
				row, col := ParseIndex(val, root.IBF.Cols, rb)
				if root.IBF.Get(row, col) {
					count++
				}
			}
			if count == len(rowVal) {
				flag = true
				count = 0
				break
			}
		}
		if flag {
			for _, rowVal := range t.T3 {
				for _, val := range rowVal {
					row, col := ParseIndex(val, root.IBF.Cols, rb)
					if root.IBF.Get(row, col) {
						count++
					}
				}
				if count == len(rowVal) {
					return true
				}
			}
		}
	}
	return false
}

func ParseIndex(str string, ibf_len, rb int) (_, _ int) {
	h := strings.Split(str, ",")
	// h[0]:outbytes,h[1]:hkp1
	outbytes, err := hex.DecodeString(h[0])
	if err != nil {
		fmt.Println("err:", err)
		return -1, -1
	}
	hkp1, err := hex.DecodeString(h[1])
	if err != nil {
		fmt.Println("err:", err)
		return -1, -1
	}

	bi := new(big.Int).SetBytes(outbytes)
	col := bi.Mod(bi, big.NewInt(int64(ibf_len))).Int64() // twin_id

	hkp1bi := new(big.Int).SetBytes(hkp1)
	sha1bytes := sha256.Sum256(hkp1bi.Xor(hkp1bi, big.NewInt(int64(rb))).Bytes())
	row := new(big.Int).SetBytes(sha1bytes[:]).Mod(new(big.Int).SetBytes(sha1bytes[:]), big.NewInt(2)).Int64()
	return int(row), int(col)
}
