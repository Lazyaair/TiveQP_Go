package construction

import (
	indexbuilding "TiveQP/Indexbuilding"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseOwner 解析每一行并创建 Owner 对象
func ParseOwner(line string) (*indexbuilding.Owner, error) {
	// 使用 "**" 分割字段
	fields := strings.Split(line, "**")
	if len(fields) != 8 {
		return nil, fmt.Errorf("数据格式错误，期望 8 个字段，但得到 %d 个字段", len(fields))
	}

	// 解析字段并转换数据类型
	lat, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("解析 Lat 出错: %v", err)
	}

	lng, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("解析 Lng 出错: %v", err)
	}

	hourStart, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("解析 HourStart 出错: %v", err)
	}

	minStart, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, fmt.Errorf("解析 MinStart 出错: %v", err)
	}

	hourClose, err := strconv.Atoi(fields[6])
	if err != nil {
		return nil, fmt.Errorf("解析 HourClose 出错: %v", err)
	}

	minClose, err := strconv.Atoi(fields[7])
	if err != nil {
		return nil, fmt.Errorf("解析 MinClose 出错: %v", err)
	}

	// 创建 Owner 对象并返回
	owner := &indexbuilding.Owner{
		Type:      fields[0],
		City:      fields[1],
		Lat:       lat,
		Lng:       lng,
		HourStart: hourStart,
		MinStart:  minStart,
		HourClose: hourClose,
		MinClose:  minClose,
	}

	return owner, nil
}

// LoadOwners 从文件中加载 Owner 对象
func LoadOwners(filename string) ([]*indexbuilding.Owner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件出错: %v", err)
	}
	defer file.Close()

	var owners []*indexbuilding.Owner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		owner, err := ParseOwner(line)
		if err != nil {
			return nil, fmt.Errorf("解析行出错: %v", err)
		}
		owners = append(owners, owner)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时出错: %v", err)
	}

	return owners, nil
}
