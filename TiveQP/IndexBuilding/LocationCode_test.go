package indexbuilding

import (
	"fmt"
	"testing"
)

func TestLocationEncoding(t *testing.T) {
	cityName := "ATLANTA"
	lat := 33.7660237
	lng := -84.5301237
	locationCode, _ := LocationEncoding(cityName, lat, lng)
	fmt.Println("Code:", locationCode)
	locationCodeComplement, _ := LocationEncodingComplement(cityName, lat, lng)
	fmt.Println("ComplementCode:", locationCodeComplement)
}
