package main

import (
	construction "TiveQP/Construction"
	indexbuilding "TiveQP/Indexbuilding"
	query "TiveQP/Query"
	resultverification "TiveQP/Resultverification"
	trapdoor "TiveQP/Trapdoor"
	"fmt"
	"net/http"
	"strings"

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

func main() {
	ibfLength := 200000
	Keylist := []string{"2938879577741549", "8729598049525437", "8418086888563864", "0128636306393258", "2942091695121238", "6518873307787549"}
	rb := 235648

	filename := "E:\\Github\\TiveQP\\TiveQP\\TiveQP\\Data\\20k.txt" // 文件名
	owners, err := construction.LoadOwners(filename)
	if err != nil {
		fmt.Println("加载 Owner 数据出错:", err)
		return
	}
	subroots, err := construction.BuildTreesByChunks(owners, ibfLength, Keylist, rb)
	if err != nil {
		fmt.Println("Error building subroots:", err)
	} else {
		fmt.Println("Subroots built successfully!")
	}
	finalRoot, err := construction.CreateFinalTree(subroots, ibfLength, Keylist, rb)
	if err != nil {
		fmt.Println("Error creating final tree:", err)
	} else {
		fmt.Println("Final tree created successfully!")
	}
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
		// TODO: 使用参数字符串查询数据库
		println("Received params:", params)
		// 去掉前缀
		cleanedStr := strings.TrimPrefix(params, "/api/message?params=")

		u, err := trapdoor.ParseUser(cleanedStr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("User loaded successfully!==Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
		}
		T, err := trapdoor.GenT(u, Keylist, rb)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("TrapDoor created successfully!==Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
		}
		k := 3
		result := make([]*[]byte, 0, k)
		pi := make([]*query.PON, 0, k)
		query.QueryT(finalRoot, T, &k, 0, rb, &result, &pi)
		fmt.Println("======================================================================")
		fmt.Println("check HV==", resultverification.CheckHV(finalRoot.HV, pi))
		fmt.Println("======================================================================")
		fmt.Println("check Completeness==", resultverification.CheckCompleteness(T, pi))
		// // 这里返回示例数据
		shops := []indexbuilding.Owner{}

		for _, v := range result {
			p, _ := construction.Decrypt(*v, []byte("2bc73dw20ebf4d46"))
			fmt.Println(string(p))
			o, err := construction.ParseOwner(string(p))
			if err != nil {
				fmt.Println(err)
			}
			shops = append(shops, *o)
			// fmt.Println("\n========================================")
		}

		// 返回响应
		c.JSON(http.StatusOK, shops)
	})

	r.Run(":8080")
}
