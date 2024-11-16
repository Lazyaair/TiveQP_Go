package indexbuilding

import (
	"fmt"
	"testing"
)

// city test
func TestCity(t *testing.T) {
	fmt.Println(GetCityIndex("ALAFAYA"))                    // 输出 2
	fmt.Println(GetCityIndex("51-RICHARD-BEALL-HWY-17-92")) // 输出 0
	fmt.Println(GetCityLatLng(2))
}

// element test
func TestPrefix(t *testing.T) {
	for i := 24; i < 47; i++ {
		fmt.Println(Prefix(6, i))
	}
}

// func TestSamesize(t *testing.T) {
// 	fmt.Println(Samesize("00101110011", "0010111001*"))
// }

func TestRange(t *testing.T) {
	fmt.Println(Range(12, 1136, 1182))
}

// location test
func TestLocationEncoding(t *testing.T) {
	cityName := "QUINCY"
	lat := 42.3328027
	lng := -71.1389101
	locationCode, _ := LocationEncoding(cityName, lat, lng)
	fmt.Println("Code:", locationCode)
	locationCodeComplement, _ := LocationEncodingComplement(cityName, lat, lng)
	fmt.Println("ComplementCode:", locationCodeComplement)
}

func TestAddCityIndex(t *testing.T) {
	cityName := "ATLANTA"
	lat := 33.7660237
	lng := -84.5301237
	locationcode, _ := LocationEncodingUser(cityName, lat, lng)
	fmt.Println("Code:", locationcode)
	AddCityIndex(cityName, locationcode)
	fmt.Println("Code with index:", locationcode)
}

// time test
func TestTimePointEncoding(t *testing.T) {
	fmt.Println(TimePointEncoding(12, 11))
}

func TestTimeRangeEncoding(t *testing.T) {
	fmt.Println(TimeRangeEncoding(8, 0, 9, 0))
}

func TestTimeRangeEncodingComplement(t *testing.T) {
	fmt.Println(TimeRangeEncodingComplement(8, 48, 8, 50))
}

// type test
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

// user test
func TestUser(t *testing.T) {

}

// owner test
func TestOwner(t *testing.T) {

}
