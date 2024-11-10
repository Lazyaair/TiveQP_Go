package construction

import (
	"fmt"
	"sync"
)

// BuildTypeTree 并行构建平衡二叉树
func BuildTypeTree(values []int) *LeafNode {
	const minParallelSize = 50000
	return buildTypeTree(values, minParallelSize)
}

// buildTypeTree 是递归构建平衡二叉树的内部函数，使用 `WaitGroup` 同步并行任务
func buildTypeTree(values []int, minParallelSize int) *LeafNode {
	if len(values) == 0 {
		return nil
	}

	// 找到中间位置，创建根节点
	mid := len(values) / 2
	root := &LeafNode{Value: values[mid]}

	// 根据数组大小决定是否并行构建左右子树
	if len(values) > minParallelSize {
		// 创建一个新的 WaitGroup
		var wg sync.WaitGroup
		wg.Add(2)

		// 使用 goroutine 构建左子树
		go func() {
			defer wg.Done()
			root.Left = buildTypeTree(values[:mid], minParallelSize)
		}()

		// 使用 goroutine 构建右子树
		go func() {
			defer wg.Done()
			root.Right = buildTypeTree(values[mid+1:], minParallelSize)
		}()

		// 等待左右子树构建完成
		wg.Wait()
	} else {
		// 当数组较小时，直接递归构建左右子树
		root.Left = buildTypeTree(values[:mid], minParallelSize)
		root.Right = buildTypeTree(values[mid+1:], minParallelSize)
	}

	return root
}

// PrintTree 中序遍历打印树节点
func PrintTree(node *LeafNode) {
	if node == nil {
		return
	}
	PrintTree(node.Left)
	fmt.Printf("%d ", node.Value)
	PrintTree(node.Right)
}
