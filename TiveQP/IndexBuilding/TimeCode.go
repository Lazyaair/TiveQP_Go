package indexbuilding

//             open 0 close 1  1 bit
//             星期一 到 星期日 0，1，2，3，4，5，6   3 bit
//             时间 0000 - 2350  6 bit 半小时为最小单位
//             日期 转化为 星期

// hour:min 时间点 用户使用
func TimePointEncoding(hour, min int) ([]string, error) {
	time := hour * 2
	if min >= 30 {
		time++
	}
	return Prefix(6, time)
}

// 时间范围 [start,end) 左闭右开 拥有者使用
func TimeRangeEncoding(hour_start, min_start, hour_end, min_end int) ([]string, error) {
	time_start := hour_start * 2
	if min_start >= 30 {
		time_start++
	}
	// [8:00-12:00]===>[16,17,...,22,23]==[2*i,2*j-1]
	time_end := hour_end*2 - 1
	if min_end >= 30 {
		time_end++
	}
	return Range(6, time_start, time_end)
}

// 时间范围补集
func TimeRangeEncodingComplement(hour_start, min_start, hour_end, min_end int) ([]string, error) {
	time_start := hour_start * 2
	if min_start >= 30 {
		time_start++
	}
	time_end := hour_end*2 - 1
	if min_end >= 30 {
		time_end++
	}

	result := make([]string, 0, 10)

	if time_start != 0 {
		f1, _ := Range(6, 0, time_start-1)
		result = append(result, f1...)
	}
	if time_end != 47 {
		f2, _ := Range(6, time_end+1, 47)
		result = append(result, f2...)
	}
	return result, nil
}
