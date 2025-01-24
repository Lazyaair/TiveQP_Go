package trapdoor

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	// "Fast Food**AUSTIN**30.2795878**-97.806248**12**12"
	u, err := ParseUser("Fast Food**AUSTIN**30.2795878**-97.806248**12**12")
	if err != nil {
		fmt.Println(err)
	}
	Keylist := []string{"2938879577741549", "8729598049525437", "8418086888563864", "0128636306393258", "2942091695121238", "6518873307787549"}
	rb := 235648
	T, err := GenT(u, Keylist, rb)
	if err != nil {
		fmt.Println(err)
	}
	// Ts := []interface{}{T.T1, T.T2, T.T3}
	Ts := []interface{}{T.T1}
	for index, trapdoor := range Ts {
		fmt.Println("T", index+1, ":")
		// 对 trapdoor 进行类型断言，确保它是 [][]string 类型
		if td, ok := trapdoor.([][]string); ok {
			for _, v1 := range td {
				for _, v2 := range v1 {
					fmt.Print(v2, "||")
				}
				fmt.Println()
			}
		} else {
			fmt.Println("Invalid type")
		}
		fmt.Println("=============")
	}
}
