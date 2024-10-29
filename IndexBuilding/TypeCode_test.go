package indexbuilding

import (
	"fmt"
	"testing"
)

func TestGetTypeIndex(t *testing.T) {
	fmt.Println(GetTypeIndex("Financial Services")) // 1
	fmt.Println(GetTypeIndex("Medical Supplies"))   // 2
}

func TestTypeEncoding(t *testing.T) {
	fmt.Println(GetTypeIndex("Dentists"))
	fmt.Println(TypeEncoding("Dentists"))
}

func TestTypeEncodingComplement(t *testing.T) {
	fmt.Println(TypeEncodingComplement("Dentists"))
}
