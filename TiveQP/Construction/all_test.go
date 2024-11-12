package construction

import (
	"fmt"
	"testing"
)

// // BuildBalancedTreeSequential 使用非并行的递归方式构建平衡二叉树
// func BuildBalancedTreeSequential(values []int) *DownNode {
// 	if len(values) == 0 {
// 		return nil
// 	}

// 	// 找到中间位置，创建根节点
// 	mid := len(values) / 2
// 	root := &DownNode{Value: values[mid]}

// 	// 递归构建左子树和右子树
// 	root.Left = BuildBalancedTreeSequential(values[:mid])
// 	root.Right = BuildBalancedTreeSequential(values[mid+1:])

// 	return root
// }

func TestTypeTree(t *testing.T) {
	// // 构建一个有序数组用于测试
	// values := make([]int, 1000000)
	// for i := 0; i < 1000000; i++ {
	// 	values[i] = i + 1
	// }

	// // 测试非并行构建平衡二叉树的时间
	// start := time.Now()
	// rootSequential := BuildBalancedTreeSequential(values)
	// elapsedSequential := time.Since(start)
	// fmt.Printf("非并行构建平衡二叉树耗时: %v\n", elapsedSequential)

	// // 测试并行构建平衡二叉树的时间
	// start = time.Now()
	// rootParallel := BuildTypeTree(values)
	// elapsedParallel := time.Since(start)
	// fmt.Printf("并行构建平衡二叉树耗时: %v\n", elapsedParallel)
	// _ = rootParallel
	// _ = rootSequential
}

func TestLoadData(t *testing.T) {
	filename := "E:\\Github\\TiveQP\\TiveQP\\TiveQP\\Data\\20k.txt" // 文件名
	owners, err := LoadOwners(filename)
	if err != nil {
		fmt.Println("加载 Owner 数据出错:", err)
		return
	}

	// 输出所有加载的 Owner 对象
	for _, owner := range owners {
		fmt.Printf("%+v\n", *owner)
	}
}
