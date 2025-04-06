package indexbuilding

import (
	"fmt"
	"testing"
)

// haveCommonElements 判断两个字符串切片是否含有相同元素
func haveCommonElements(slice1, slice2 []string) bool {
	// 如果任一切片为空，则无共同元素
	if len(slice1) == 0 || len(slice2) == 0 {
		return false
	}

	// 创建 map 用于存储第一个切片的元素
	set := make(map[string]bool)
	for _, elem := range slice1 {
		set[elem] = true
	}

	// 检查第二个切片的元素是否在 map 中
	for _, elem := range slice2 {
		if set[elem] {
			fmt.Println("===>", elem, "<===")
			return true // 发现共同元素，返回 true
		}
	}

	// 未找到共同元素，返回 false
	return false
}

// city test
func TestCity(t *testing.T) {
	fmt.Println(GetCityIndex("ALAFAYA"))                    // 输出 2
	fmt.Println(GetCityIndex("51-RICHARD-BEALL-HWY-17-92")) // 输出 0
	fmt.Println(GetCityLatLng(2))
}

// element test
func TestPrefix(t *testing.T) {
	// for i := 3; i < 19; i++ {
	// 	fmt.Println(Prefix(6, i))
	// }
	rn, _ := Range(12, 2450, 2456)
	n, _ := prefixRangeUniqueWithStars(12, 2452, 2452)
	fmt.Println("User code: ", rn)
	fmt.Println("CS   code: ", rn)
	fmt.Println(haveCommonElements(rn, n))
	// n1, _ := Prefix(12, 2451)
	// n2, _ := Prefix(12, 2452)
	// n3, _ := Prefix(12, 2453)
	// fmt.Println(haveCommonElements(rn, n1))
	// fmt.Println(haveCommonElements(rn, n2))
	// fmt.Println(haveCommonElements(rn, n3))
	// fmt.Println(Prefix(6, 35))
}

// func TestSamesize(t *testing.T) {
// 	fmt.Println(Samesize("00101110011", "0010111001*"))
// }

func TestRange(t *testing.T) {
	fmt.Println(Range(4, 1, 3))
	fmt.Println(Range(4, 0, 6))
}

// location test
func TestLocationEncoding(t *testing.T) {
	cityName := "ALTANTA"
	lat := 42.3328027
	lng := -71.1389101
	locationCode, _ := LocationEncoding(cityName, lat, lng)
	fmt.Println("Code          :", locationCode)
	UserCode, _ := LocationEncodingUser(cityName, 2, lat, lng)
	fmt.Println("UserCode      :", UserCode)
	fmt.Println(haveCommonElements(locationCode, UserCode))
	locationCodeComplement, _ := LocationEncodingComplement(cityName, lat, lng)
	fmt.Println("ComplementCode:", locationCodeComplement)
	fmt.Println(haveCommonElements(locationCodeComplement, UserCode))
	fmt.Println(haveCommonElements(locationCodeComplement, locationCode))
}

func TestAddCityIndex(t *testing.T) {
	cityName := "ATLANTA"
	lat := 33.7660237
	lng := -84.5301237
	locationcode, _ := LocationEncodingUser(cityName, 3, lat, lng)
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
	fmt.Println(GetTypeIndex("Restaurants"))
	fmt.Println(TypeEncoding("Restaurants"))
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
