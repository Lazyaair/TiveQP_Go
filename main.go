package main

import (
	construction "TiveQP/Construction"
	indexbuilding "TiveQP/Indexbuilding"
	query "TiveQP/Query"
	resultverification "TiveQP/Resultverification"
	trapdoor "TiveQP/Trapdoor"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 请求参数结构体
type SearchRequest struct {
	Type      string  `form:"type"`      // 店铺类型
	City      string  `form:"city"`      // 城市
	Lat       float64 `form:"lat"`       // 纬度
	Lng       float64 `form:"lng"`       // 经度
	HourStart int     `form:"hourStart"` // 开始时间（小时）
	MinStart  int     `form:"minStart"`  // 开始时间（分钟）
}

// 响应数据结构体
type Shop struct {
	Type      string  `json:"type"`
	City      string  `json:"city"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	HourStart int     `json:"hourStart"`
	MinStart  int     `json:"minStart"`
	HourClose int     `json:"hourClose"`
	MinClose  int     `json:"minClose"`
}

// 完整响应数据结构体
type TiveQPData struct {
	Shops    []indexbuilding.Owner `json:"shops"`
	Trapdoor *trapdoor.T           `json:"trapdoor"`
}

func isShopOpen(shop Shop, hour, minute int) bool {
	currentTime := hour*60 + minute
	shopOpenTime := shop.HourStart*60 + shop.MinStart
	shopCloseTime := shop.HourClose*60 + shop.MinClose
	return currentTime >= shopOpenTime && currentTime <= shopCloseTime
}
func main() {
	ibfLength := 200000
	Keylist := []string{"2938879577741549", "8729598049525437", "8418086888563864", "0128636306393258", "2942091695121238", "6518873307787549"}
	rb := 235648
	lastCount := ""
	filename := "./Data/20k.txt" // 文件名
	owners, err := construction.LoadOwners(filename)
	if err != nil {
		fmt.Println("加载 Owner 数据出错:", err)
		return
	}
	var subroots []*construction.Node
	var finalRoot *construction.Node
	r := gin.Default()

	// 启用CORS中间件
	r.Use(cors.Default())

	r.GET("/api/message", func(c *gin.Context) {
		// 获取参数字符串
		params := c.Query("params")
		if params == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "参数不能为空",
			})
			return
		}

		// 打印接收到的参数字符串（调试用）
		println("Received params:", params)
		// 去掉前缀
		cleanedStr := strings.TrimPrefix(params, "/api/message?params=")
		fmt.Printf("Cleaned params string: %s\n", cleanedStr)

		// 从参数中获取最大店铺数量，并将其从参数字符串中移除
		parts := strings.Split(cleanedStr, "**")
		fmt.Printf("Split parts: %v\n", parts)

		k := 3 // 默认值
		if len(parts) >= 8 {
			if maxShops, err := strconv.Atoi(parts[6]); err == nil {
				k = maxShops
				fmt.Printf("Using maxShops value: %d\n", k)
			} else {
				fmt.Printf("Error parsing maxShops: %v\n", err)
			}
			if parts[7] != lastCount {
				subroots = nil
				finalRoot = nil
				lastCount = parts[7]
				indexbuilding.SplitCount = indexbuilding.LocationMap[parts[7]]
				indexbuilding.Bitsize = indexbuilding.LocationSizeMap[parts[7]]
				subroots, err = construction.BuildTreesByChunks(owners, ibfLength, Keylist, rb)
				if err != nil {
					fmt.Println("Error building subroots:", err)
				} else {
					fmt.Println("Subroots built successfully!")
				}
				finalRoot, err = construction.CreateFinalTree(subroots, ibfLength, Keylist, rb)
				if err != nil {
					fmt.Println("Error creating final tree:", err)
				} else {
					fmt.Println("Final tree created successfully!")
				}
			}
			// 只保留前6个参数传递给ParseUser
			cleanedStr = strings.Join(parts[:6], "**")
			fmt.Printf("Modified params string for ParseUser: %s\n", cleanedStr)
		}

		u, err := trapdoor.ParseUser(cleanedStr)
		if err != nil {
			fmt.Printf("ParseUser error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "解析用户参数失败: " + err.Error(),
			})
			return
		} else {
			fmt.Println("User loaded successfully!")
		}

		T, err := trapdoor.GenT(u, Keylist, rb)
		if err != nil {
			fmt.Printf("GenT error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "生成陷门失败: " + err.Error(),
			})
			return
		} else {
			fmt.Println("TrapDoor created successfully!")
		}

		result := make([]*[]byte, 0, k)
		pi := make([]*query.PON, 0, k)
		query.QueryT(finalRoot, T, &k, 0, rb, &result, &pi)
		fmt.Printf("QueryT completed. Result count: %d\n", len(result))

		fmt.Println("======================================================================")
		fmt.Println("check HV==", resultverification.CheckHV(finalRoot.HV, pi))
		fmt.Println("======================================================================")
		fmt.Println("check Completeness==", resultverification.CheckCompleteness(T, pi))

		shops := []indexbuilding.Owner{}

		for _, v := range result {
			p, err := construction.Decrypt(*v, []byte("2bc73dw20ebf4d46"))
			if err != nil {
				fmt.Printf("Decrypt error: %v\n", err)
				continue
			}
			fmt.Printf("Decrypted data: %s\n", string(p))
			o, err := construction.ParseOwner(string(p))
			if err != nil {
				fmt.Printf("ParseOwner error: %v\n", err)
				continue
			}
			shops = append(shops, *o)
		}

		fmt.Printf("Final shops count: %d\n", len(shops))
		response := TiveQPData{
			Shops:    shops,
			Trapdoor: T,
		}
		// 返回响应
		c.JSON(http.StatusOK, response)
	})

	// 缓存所有店铺数据
	var allShops []Shop
	var shopTypes []string
	var cities []string
	var initialized bool
	var mutex sync.RWMutex

	// 初始化数据的函数
	initializeData := func() error {
		mutex.Lock()
		defer mutex.Unlock()

		if initialized {
			return nil
		}

		file, err := os.Open("./Data/20k.txt")
		if err != nil {
			return err
		}
		defer file.Close()

		typeSet := make(map[string]bool)
		citySet := make(map[string]bool)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			parts := strings.Split(line, "**")
			if len(parts) != 8 {
				continue
			}

			lat, _ := strconv.ParseFloat(parts[2], 64)
			lng, _ := strconv.ParseFloat(parts[3], 64)
			hourStart, _ := strconv.Atoi(parts[4])
			minStart, _ := strconv.Atoi(parts[5])
			hourClose, _ := strconv.Atoi(parts[6])
			minClose, _ := strconv.Atoi(parts[7])

			shop := Shop{
				Type:      parts[0],
				City:      parts[1],
				Lat:       lat,
				Lng:       lng,
				HourStart: hourStart,
				MinStart:  minStart,
				HourClose: hourClose,
				MinClose:  minClose,
			}

			allShops = append(allShops, shop)
			typeSet[parts[0]] = true
			citySet[parts[1]] = true
		}

		// 转换集合为切片
		for t := range typeSet {
			shopTypes = append(shopTypes, t)
		}
		for c := range citySet {
			cities = append(cities, c)
		}

		sort.Strings(shopTypes)
		sort.Strings(cities)

		initialized = true
		return scanner.Err()
	}

	// 获取店铺统计信息
	r.GET("/api/shops/stats", func(c *gin.Context) {
		if err := initializeData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mutex.RLock()
		defer mutex.RUnlock()

		// 计算当前营业的店铺数量
		now := time.Now()
		currentHour := now.Hour()
		currentMin := now.Minute()
		openCount := 0
		for _, shop := range allShops {
			if isShopOpen(shop, currentHour, currentMin) {
				openCount++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"totalShops":  len(allShops),
			"totalTypes":  len(shopTypes),
			"totalCities": len(cities),
			"openShops":   openCount,
			"types":       shopTypes,
			"cities":      cities,
		})
	})

	// 分页获取店铺数据
	r.GET("/api/shops", func(c *gin.Context) {
		if err := initializeData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mutex.RLock()
		defer mutex.RUnlock()

		// 获取查询参数
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		shopType := c.Query("type")
		city := c.Query("city")
		timeStr := c.Query("time")

		// 过滤数据
		var filteredShops []Shop
		for _, shop := range allShops {
			if shopType != "" && shop.Type != shopType {
				continue
			}
			if city != "" && shop.City != city {
				continue
			}
			if timeStr != "" {
				parts := strings.Split(timeStr, ":")
				if len(parts) == 2 {
					hour, _ := strconv.Atoi(parts[0])
					min, _ := strconv.Atoi(parts[1])
					if !isShopOpen(shop, hour, min) {
						continue
					}
				}
			}
			filteredShops = append(filteredShops, shop)
		}

		// 计算分页
		total := len(filteredShops)
		start := (page - 1) * pageSize
		end := start + pageSize
		if start >= total {
			start = total
		}
		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, gin.H{
			"total": total,
			"data":  filteredShops[start:end],
		})
	})

	r.Run(":8080")
}
