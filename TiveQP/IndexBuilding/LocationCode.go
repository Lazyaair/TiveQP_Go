package indexbuilding

import (
	"fmt"
	"math"
)

// 分割段数
// 位长
var splitCount = 50
var bitsize = 12

// 区间投影
func Projection(minVal, maxVal, currentVal float64) int {
	result := (currentVal - minVal) / (maxVal - minVal) * float64(splitCount)
	return int(math.Floor(result))
}

// 网格分割
// 边界判断
// 坐标编码
func LocationEncoding(cityName string, lat, lng float64) ([]string, error) {
	cityIndex := GetCityIndex(cityName)
	cityBoundary := GetCityLatLng(cityIndex)

	// 经纬度投影编号
	num_lat := Projection(cityBoundary[0], cityBoundary[1], lat)
	num_lng := Projection(cityBoundary[2], cityBoundary[3], lng)

	// 网格编号
	place_index := num_lng*splitCount + num_lat

	result := make([]string, 0, 6)
	switch {
	case num_lat == 0 && num_lng == 0:
		// 左下角
		n1, _ := Range(bitsize, place_index, place_index+1)
		n2, _ := Range(bitsize, place_index+splitCount, place_index+splitCount+1)
		result = append(result, n1...)
		result = append(result, n2...)
		return result, nil
	case num_lat == 0 && num_lng == splitCount-1:
		// 左上角
		n1, _ := Range(bitsize, place_index, place_index+1)
		n2, _ := Range(bitsize, place_index-splitCount, place_index-splitCount+1)
		result = append(result, n2...)
		result = append(result, n1...)
		return result, nil
	case num_lat == splitCount-1 && num_lng == 0:
		// 右下角
		n1, _ := Range(bitsize, place_index-1, place_index)
		n2, _ := Range(bitsize, place_index+splitCount-1, place_index+splitCount)
		result = append(result, n1...)
		result = append(result, n2...)
		return result, nil
	case num_lat == splitCount-1 && num_lng == splitCount-1:
		// 右上角
		n1, _ := Range(bitsize, place_index-1, place_index)
		n2, _ := Range(bitsize, place_index-splitCount-1, place_index-splitCount)
		result = append(result, n2...)
		result = append(result, n1...)
		return result, nil
	case num_lat == 0 && num_lng > 0 && num_lng < splitCount-1:
		// 左 边
		n1, _ := Range(bitsize, place_index-splitCount, place_index-splitCount+1)
		n2, _ := Range(bitsize, place_index, place_index+1)
		n3, _ := Range(bitsize, place_index+splitCount, place_index+splitCount+1)
		result = append(result, n1...)
		result = append(result, n2...)
		result = append(result, n3...)
		return result, nil
	case num_lat == splitCount-1 && 0 < num_lng && num_lng < splitCount-1:
		// 右 边
		n1, _ := Range(bitsize, place_index-splitCount-1, place_index-splitCount)
		n2, _ := Range(bitsize, place_index-1, place_index)
		n3, _ := Range(bitsize, place_index+splitCount-1, place_index+splitCount)
		result = append(result, n1...)
		result = append(result, n2...)
		result = append(result, n3...)
		return result, nil
	case 0 < num_lat && num_lat < splitCount-1 && num_lng == 0:
		// 下 边
		n1, _ := Range(bitsize, place_index-1, place_index+1)
		n2, _ := Range(bitsize, place_index+splitCount-1, place_index+splitCount+1)
		result = append(result, n1...)
		result = append(result, n2...)
		return result, nil
	case 0 < num_lat && num_lat < splitCount-1 && num_lng == splitCount-1:
		// 上 边
		n1, _ := Range(bitsize, place_index-1, place_index+1)
		n2, _ := Range(bitsize, place_index-splitCount-1, place_index-splitCount+1)
		result = append(result, n2...)
		result = append(result, n1...)
		return result, nil
	case 0 < num_lat && num_lat < splitCount-1 && 0 < num_lng && num_lng < splitCount-1:
		// 中 间
		n1, _ := Range(bitsize, place_index-splitCount-1, place_index-splitCount+1)
		n2, _ := Range(bitsize, place_index-1, place_index+1)
		n3, _ := Range(bitsize, place_index+splitCount-1, place_index+splitCount+1)
		result = append(result, n1...)
		result = append(result, n2...)
		result = append(result, n3...)
		return result, nil
	default:
		return nil, fmt.Errorf("经纬度编号有误！")
	}
}
