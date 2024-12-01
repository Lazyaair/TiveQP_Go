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
	for index, trapdoor := range T {
		fmt.Println("T", index+1, ":")
		for _, v1 := range trapdoor {
			for _, v2 := range v1 {
				fmt.Print(v2, "||")
			}
			fmt.Println()
		}
		fmt.Println("=============")
	}
}
