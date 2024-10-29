package indexbuilding

import (
	"fmt"
	"testing"
)

func TestCity(t *testing.T) {
	fmt.Println(GetCityIndex("ALAFAYA"))                    // 输出 2
	fmt.Println(GetCityIndex("51-RICHARD-BEALL-HWY-17-92")) // 输出 0
	fmt.Println(GetCityLatLng(2))
}
