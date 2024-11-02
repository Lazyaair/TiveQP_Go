package indexbuilding

import (
	"fmt"
	"testing"
)

func TestOwnerMethods(t *testing.T) {
	owner := &Owner{}

	// 测试 OwnerType 方法
	t.Run("OwnerType", func(t *testing.T) {
		typeStr := "someType"
		prefixes, err := owner.OwnerType(typeStr)
		if err != nil {
			t.Fatalf("OwnerType failed: %v", err)
		}
		fmt.Println("OwnerType Prefixes:", prefixes)
	})

	// 测试 OwnerTypeComplement 方法
	t.Run("OwnerTypeComplement", func(t *testing.T) {
		typeStr := "someType"
		complement, err := owner.OwnerTypeComplement(typeStr)
		if err != nil {
			t.Fatalf("OwnerTypeComplement failed: %v", err)
		}
		fmt.Println("OwnerTypeComplement:", complement)
	})

	// 测试 OwnerLocation 方法
	t.Run("OwnerLocation", func(t *testing.T) {
		city := "Atlanta, GA"
		lat := 33.7660237
		lng := -84.5301237
		locationCode, err := owner.OwnerLocation(city, lat, lng)
		if err != nil {
			t.Fatalf("OwnerLocation failed: %v", err)
		}
		fmt.Println("OwnerLocation Code:", locationCode)
	})

	// 测试 OwnerLocationComplement 方法
	t.Run("OwnerLocationComplement", func(t *testing.T) {
		city := "Atlanta, GA"
		lat := 33.7660237
		lng := -84.5301237
		locationComplementCode, err := owner.OwnerLocationComplement(city, lat, lng)
		if err != nil {
			t.Fatalf("OwnerLocationComplement failed: %v", err)
		}
		fmt.Println("OwnerLocationComplement Code:", locationComplementCode)
	})

	// 测试 OwnerTime 方法
	t.Run("OwnerTime", func(t *testing.T) {
		hourOpen := 9
		minOpen := 0
		hourClose := 17
		minClose := 30
		timeCodes, err := owner.OwnerTime(hourOpen, minOpen, hourClose, minClose)
		if err != nil {
			t.Fatalf("OwnerTime failed: %v", err)
		}
		fmt.Println("OwnerTime Codes:", timeCodes)
	})

	// 测试 OwnerTimeComplement 方法
	t.Run("OwnerTimeComplement", func(t *testing.T) {
		hourOpen := 9
		minOpen := 0
		hourClose := 17
		minClose := 30
		timeComplementCodes, err := owner.OwnerTimeComplement(hourOpen, minOpen, hourClose, minClose)
		if err != nil {
			t.Fatalf("OwnerTimeComplement failed: %v", err)
		}
		fmt.Println("OwnerTimeComplement Codes:", timeComplementCodes)
	})
}
