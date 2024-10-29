package main

import (
	indexbuilding "TiveQP/IndexBuilding"
	"fmt"
)

func main() {
	fmt.Println(indexbuilding.GetCityIndex("ALAFAYA")) // 输出 2
	fmt.Println(indexbuilding.GetCityLatLng(2))        // 输出 -1
}
