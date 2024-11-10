package construction

import (
	"fmt"
	"testing"
	"time"
)

func TestTypeTree(t *testing.T) {
	// 构建一个有序数组用于测试
	values := make([]int, 20000)
	for i := 0; i < 20000; i++ {
		values[i] = i + 1
	}

	// 测试并行构建平衡二叉树的时间
	start := time.Now()
	rootParallel := BuildTypeTree(values)
	elapsedParallel := time.Since(start)
	fmt.Printf("并行构建平衡二叉树耗时: %v\n", elapsedParallel)
	_ = rootParallel

}
