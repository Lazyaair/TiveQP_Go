package indexbuilding

import (
	"fmt"
	"testing"
)

func TestPrefix(t *testing.T) {
	for i := 24; i < 47; i++ {
		fmt.Println(Prefix(6, i))
	}
}

// func TestSamesize(t *testing.T) {
// 	fmt.Println(Samesize("00101110011", "0010111001*"))
// }

func TestRange(t *testing.T) {
	fmt.Println(Range(6, 16, 16))
}
