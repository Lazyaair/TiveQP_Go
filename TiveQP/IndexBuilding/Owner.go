package indexbuilding

import (
	"strings"
)

// Owner 结构体
type Owner struct{}

// OwnerType 返回 type 对应编号的前缀族
func (o *Owner) OwnerType(typeStr string) ([]string, error) {
	index, err := TypeEncoding(typeStr)
	if err != nil {
		return nil, err
	}
	return Prefix(11, index), nil
}

// OwnerTypeComplement 返回 type 编号对应的补集
func (o *Owner) OwnerTypeComplement(typeStr string) ([]string, error) {
	return TypeComplement(typeStr)
}

// OwnerLocation 对 Owner 位置编码
func (o *Owner) OwnerLocation(city string, lat, lng float64) []string {
	cityName := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(city, ",", "-"), " ", "-"))
	encode := CityEncoding(cityName, lat, lng)
	return AddCityNumber(cityName, encode)
}

// OwnerLocationComplement 求 Owner 位置的补集
func (o *Owner) OwnerLocationComplement(city string, lat, lng float64) ([]string, error) {
	cityName := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(city, ",", "-"), " ", "-"))
	encode := CityEncodingComplement(cityName, lat, lng)
	return AddCityNumber(cityName, encode), nil
}

// OwnerTime 对输入时间进行编码
func (o *Owner) OwnerTime(hourOpen, minOpen, hourClose, minClose int) []string {
	return TimeEncodingOwner(hourOpen, minOpen, hourClose, minClose)
}

// OwnerTimeComplement 补集
func (o *Owner) OwnerTimeComplement(hourOpen, minOpen, hourClose, minClose int) []string {
	return TimeEncodingOwnerComplement(hourOpen, minOpen, hourClose, minClose)
}
