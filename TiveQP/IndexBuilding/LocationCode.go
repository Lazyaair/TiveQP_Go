package indexbuilding

import "math"

// 分割段数
var splitCount = 50

// 区间投影
func Projection(minVal, maxVal, currentVal float64) float64 {
	result := (currentVal - minVal) / (maxVal - minVal) * float64(splitCount)
	return math.Floor(result)
}

// 网格分割
// 边界判断
// 坐标编码
func LocationEncoding(cityName string, lat, lng float64) {
	cityIndex := GetCityIndex(cityName)
	cityBoundary := GetCityLatLng(cityIndex)

	// 经纬度投影编号
	num_lat := Projection(cityBoundary[0], cityBoundary[1], lat)
	num_lng := Projection(cityBoundary[2], cityBoundary[3], lng)

}
