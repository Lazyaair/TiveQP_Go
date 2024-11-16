package construction

import (
	indexbuilding "TiveQP/IndexBuilding"
	"fmt"
)

// BuildTreesByChunks 函数用于将owners分为多个小块，每1000个元素建一棵树
func BuildTreesByChunks(owners []*indexbuilding.Owner, ibfLength int, Keylist []string, rb int) ([]*Node, error) {
	var subroots []*Node
	chunkSize := 1000
	numChunks := len(owners) / chunkSize

	// 对每一块数据（每1000个元素），构建一棵树
	for i := 0; i < numChunks; i++ {
		startIdx := i * chunkSize
		endIdx := (i + 1) * chunkSize
		chunk := owners[startIdx:endIdx]

		// 构建每个小块的树
		treeRoot, err := BuildTree(chunk, ibfLength, Keylist, rb)
		if err != nil {
			return nil, fmt.Errorf("failed to build tree for chunk %d: %v", i, err)
		}

		// 获取子树的根节点（即每棵树的最上层节点）
		err = treeRoot.InitUpLeafNode(chunk[0].Type, ibfLength, Keylist, rb)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize subroot for chunk %d: %v", i, err)
		}

		subroots = append(subroots, treeRoot)
	}

	return subroots, nil
}

// BuildTree 函数构建单棵树
func BuildTree(owners []*indexbuilding.Owner, ibfLength int, Keylist []string, rb int) (*Node, error) {
	// 如果 owners 数组为空，返回 nil
	if len(owners) == 0 {
		return nil, fmt.Errorf("owners list is empty")
	}

	// 当只有一个 Owner 时，创建叶节点并返回
	if len(owners) == 1 {
		leafNode := &Node{}
		err := leafNode.InitLeafNode(owners[0], ibfLength, Keylist, rb)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize leaf node: %v", err)
		}
		return leafNode, nil
	}

	// 将 owners 切割成两个部分
	mid := len(owners) / 2
	leftOwners := owners[:mid]
	rightOwners := owners[mid:]

	// 递归构建左子树和右子树
	leftNode, err := BuildTree(leftOwners, ibfLength, Keylist, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to build left subtree: %v", err)
	}

	rightNode, err := BuildTree(rightOwners, ibfLength, Keylist, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to build right subtree: %v", err)
	}

	// 创建一个中间节点，并合并左右子树
	midNode := &Node{
		Left:  leftNode,
		Right: rightNode,
	}

	// 初始化中间节点
	err = midNode.InitMidNode(ibfLength, Keylist, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mid node: %v", err)
	}

	return midNode, nil
}

// CreateFinalTree 函数用于将20个subroot节点再次合并成一棵树
func CreateFinalTree(subroots []*Node, ibfLength int, Keylist []string, rb int) (*Node, error) {
	// 检查 subroots 的数量
	if len(subroots) == 0 {
		return nil, fmt.Errorf("subroots list is empty")
	}

	// 递归地构建最终的树
	for len(subroots) > 1 {
		var nextLevel []*Node
		for i := 0; i < len(subroots); i += 2 {
			// 如果剩余节点是奇数个，则最后一个单独成对
			if i+1 < len(subroots) {
				left := subroots[i]
				right := subroots[i+1]

				// 创建中间节点，并将左右节点合并
				midNode := &Node{
					Left:  left,
					Right: right,
				}

				// 初始化中间节点
				err := midNode.InitMidNode(ibfLength, Keylist, rb)
				if err != nil {
					return nil, fmt.Errorf("failed to initialize mid node: %v", err)
				}
				nextLevel = append(nextLevel, midNode)
			} else {
				nextLevel = append(nextLevel, subroots[i])
			}
		}
		// 进入下一轮合并
		subroots = nextLevel
	}

	// 最后生成最终根节点，调用 InitUpMid_RootNode 初始化根节点
	finalRoot := subroots[0]
	err := finalRoot.InitUpMid_RootNode(ibfLength, Keylist, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize final root node: %v", err)
	}

	return finalRoot, nil
}

// PreOrderTraversal 实现前序遍历
func (n *Node) PreOrderTraversal() {
	if n == nil {
		return
	}

	// 访问当前节点（根节点）
	fmt.Printf("%v", n) // 可以根据需要打印其他信息，比如节点的值

	// 递归遍历左子树
	if n.Left != nil {
		n.Left.PreOrderTraversal()
	}

	// 递归遍历右子树
	if n.Right != nil {
		n.Right.PreOrderTraversal()
	}
}
