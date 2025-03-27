package indexbuilding

import (
	"fmt"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

var SplitCount = 50 // 分割段数
var Bitsize = 12    // 位长
var LocationMap = map[string]int{"1": 50, "2": 25, "3": 17, "4": 13, "5": 10}
var LocationSizeMap = map[string]int{"1": 12, "2": 11, "3": 9, "4": 8, "5": 7}

// 区间投影
func Projection(minVal, maxVal, currentVal float64) int {
	// // 用 *big.Float 来表示高精度的经纬度
	// minF := new(big.Float).SetFloat64(minVal)
	// maxF := new(big.Float).SetFloat64(maxVal)
	// currentF := new(big.Float).SetFloat64(currentVal)
	// // 计算比例: (currentVal - minVal) / (maxVal - minVal) * splitCount
	// delta := new(big.Float).Sub(currentF, minF)                                 // currentVal - minVal
	// rangeVal := new(big.Float).Sub(maxF, minF)                                  // maxVal - minVal
	// proportion := new(big.Float).Quo(delta, rangeVal)                           // (currentVal - minVal) / (maxVal - minVal)
	// scaled := new(big.Float).Mul(proportion, big.NewFloat(float64(splitCount))) // 比例 * splitCount
	// // 向下取整: 使用 Int 来丢弃小数部分
	// result, _ := scaled.Int(nil) // 转换为整数
	// return (int(result.Int64()))
	if currentVal <= minVal {
		return 0
	}
	if currentVal >= maxVal {
		return SplitCount - 1
	}
	result := int((currentVal - minVal) / (maxVal - minVal) * float64(SplitCount))
	if result == SplitCount {
		return SplitCount - 1
	} else {
		return result
	}
}

// 为编码添加城市索引
func AddCityIndex(cityName string, code []string) error { //
	index, err := GetCityIndex(cityName)
	if err != nil {
		return fmt.Errorf("city not exists")
	}
	prefix := strconv.Itoa(index) + ":"
	for i, v := range code {
		var builder strings.Builder
		builder.Grow(len(prefix) + len(v)) // 预分配内存，避免多次扩容
		builder.WriteString(prefix)
		builder.WriteString(v)
		code[i] = builder.String()
	}
	return nil
}

// 对用户
func LocationEncodingUser(cityName string, x int, lat, lng float64) ([]string, error) {
	cityIndex, err := GetCityIndex(cityName)
	if err != nil {
		return nil, fmt.Errorf("city not exists")
	}
	cityBoundary := GetCityLatLng(cityIndex)

	// 经纬度投影编号
	num_lat := Projection(cityBoundary[0], cityBoundary[1], lat)
	num_lng := Projection(cityBoundary[2], cityBoundary[3], lng)

	// if x == 1 {
	// 	return Prefix(Bitsize, num_lat*SplitCount+num_lng)
	// }
	k_x := x - 1
	sub_row_start := max(0, num_lat-k_x)
	sub_row_end := min(49, num_lat+k_x)
	sub_col_start := max(0, num_lng-k_x)
	sub_col_end := min(49, num_lng+k_x)

	result := make([]string, 0, 2*x-1)
	set := mapset.NewSet()

	for i := sub_row_start; i <= sub_row_end; i++ {
		// lSet, _ := prefixRangeUniqueWithStars(Bitsize, i*50+sub_col_start, i*50+sub_col_end)
		// for _, num := range lSet {
		// 	set.Add(num)
		// }
		for j := i*50 + sub_col_start; j <= i*50+sub_col_end; j++ {
			lSet, _ := Prefix(Bitsize, j)
			for _, num := range lSet {
				set.Add(num)
			}
		}
		// fmt.Println("[", i*50+sub_col_start, ",", i*50+sub_col_end, "]")
	}
	for elem := range set.Iter() {
		result = append(result, elem.(string))
	}
	return result, nil
}

// 网格分割
// 边界判断
// 坐标编码
// 对拥有者
func LocationEncoding(cityName string, lat, lng float64) ([]string, error) {
	cityIndex, err := GetCityIndex(cityName)
	if err != nil {
		return nil, fmt.Errorf("city not exists")
	}
	cityBoundary := GetCityLatLng(cityIndex)
	// fmt.Println(cityBoundary)

	// 经纬度投影编号
	num_lat := Projection(cityBoundary[0], cityBoundary[1], lat)
	num_lng := Projection(cityBoundary[2], cityBoundary[3], lng)

	result := make([]string, 0, 9)
	row_start := max(0, num_lat-4)
	row_end := min(49, num_lat+4)
	col_start := max(0, num_lng-4)
	col_end := min(49, num_lng+4)
	for i := row_start; i <= row_end; i++ {
		lSet, _ := Range(Bitsize, i*SplitCount+col_start, i*SplitCount+col_end)
		// fmt.Println("[", i*SplitCount+col_start, ",", i*SplitCount+col_end, "]")
		result = append(result, lSet...)
	}
	return result, nil
}

// 对拥有者
func LocationEncodingComplement(cityName string, lat, lng float64) ([]string, error) {
	cityIndex, err := GetCityIndex(cityName)
	if err != nil {
		return nil, fmt.Errorf("city not exists")
	}
	cityBoundary := GetCityLatLng(cityIndex)

	// 经纬度投影编号
	num_lat := Projection(cityBoundary[0], cityBoundary[1], lat)
	num_lng := Projection(cityBoundary[2], cityBoundary[3], lng)

	result := make([]string, 0, 18)

	row_start := max(0, num_lat-4)
	row_end := min(49, num_lat+4)
	col_start := max(0, num_lng-4)
	col_end := min(49, num_lng+4)

	// (1) 下方行和中间行左侧的合并
	if col_start > 0 {
		// 合并下方行和中间行左侧第一行
		if row_start > 0 {
			lSet, _ := Range(Bitsize, 0, row_start*50+col_start-1)
			// fmt.Println("[", 0, ",", row_start*50+col_start-1, "]")
			result = append(result, lSet...)
		}

		// 中间行左侧剩余部分逐行处理
		for i := row_start + 1; i <= row_end; i++ {
			lSet, _ := Range(Bitsize, i*50, i*50+col_start-1)
			// fmt.Println("[", i*50, ",", i*50+col_start-1, "]")
			result = append(result, lSet...)
		}
	}

	// (2) 中间行右侧和上方行的合并
	if col_end < 49 {
		// 中间行右侧的前几行逐行处理
		for i := row_start; i < row_end; i++ {
			lSet, _ := Range(Bitsize, i*50+col_end+1, i*50+49)
			// fmt.Println("[", i*50+col_end+1, ",", i*50+49, "]")
			result = append(result, lSet...)
		}
		// 合并中间行右侧最后一行和上方行
		if row_end < 49 {
			lSet, _ := Range(Bitsize, row_end*50+col_end+1, 49*50+49)
			// fmt.Println("[", row_end*50+col_end+1, ",", 49*50+49, "]")
			result = append(result, lSet...)
		}
	}
	return result, nil
}
