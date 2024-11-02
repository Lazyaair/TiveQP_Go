package indexbuilding

import (
	"strings"
)

// Owner 结构体
type Owner struct{}

// OwnerType 返回 type 对应编号的前缀族
func (o *Owner) OwnerType(typeStr string) ([]string, error) {
	index := GetTypeIndex(typeStr)

	prefixList, err := Prefix(11, index)
	if err != nil {
		return nil, err
	}
	return prefixList, nil
}

// OwnerTypeComplement 返回 type 编号对应的补集
func (o *Owner) OwnerTypeComplement(typeStr string) ([]string, error) {
	return TypeEncodingComplement(typeStr)
}

// OwnerLocation 对 Owner 位置编码
func (o *Owner) OwnerLocation(city string, lat, lng float64) ([]string, error) {
	cityName := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(city, ",", "-"), " ", "-"))
	encode, err := LocationEncoding(cityName, lat, lng)
	if err != nil {
		return nil, err
	}
	return AddCityNumber(cityName, encode)
}

// OwnerLocationComplement 求 Owner 位置的补集
func (o *Owner) OwnerLocationComplement(city string, lat, lng float64) ([]string, error) {
	cityName := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(city, ",", "-"), " ", "-"))
	encode, err := LocationEncodingComplement(cityName, lat, lng)
	if err != nil {
		return nil, err
	}
	return AddCityNumber(cityName, encode)
}

// OwnerTime 对输入时间进行编码
func (o *Owner) OwnerTime(hourOpen, minOpen, hourClose, minClose int) ([]string, error) {
	return TimeRangeEncoding(hourOpen, minOpen, hourClose, minClose)
}

// OwnerTimeComplement 补集
func (o *Owner) OwnerTimeComplement(hourOpen, minOpen, hourClose, minClose int) ([]string, error) {
	return TimeRangeEncodingComplement(hourOpen, minOpen, hourClose, minClose)
}
