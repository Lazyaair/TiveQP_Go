package construction

import (
	indexbuilding "TiveQP/Indexbuilding"
	"fmt"
	"sync"
)

// 将owners分为多个小块，每1000个元素（每个Type）建一棵树
func BuildTreesByChunks(owners []*indexbuilding.Owner, ibfLength int, Keylist []string, rb int) ([]*Node, error) {
	var subroots []*Node
	chunkSize := 1000
	numChunks := (len(owners) + chunkSize - 1) / chunkSize // 计算总块数

	// 定义通道用于接收结果和错误
	results := make(chan *Node, numChunks)
	errors := make(chan error, numChunks)

	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		startIdx := i * chunkSize
		endIdx := startIdx + chunkSize
		if endIdx > len(owners) {
			endIdx = len(owners)
		}
		chunk := owners[startIdx:endIdx]

		// 启动 Goroutine 构建每个小块的树
		wg.Add(1)
		go func(chunk []*indexbuilding.Owner, idx int) {
			defer wg.Done()
			treeRoot, err := BuildTree(chunk, ibfLength, Keylist, rb)
			if err != nil {
				errors <- fmt.Errorf("failed to build tree for chunk %d: %v", idx, err)
				return
			}

			// 初始化子树的根节点（上层叶节点）
			err = treeRoot.InitUpLeafNode(chunk[0].Type, ibfLength, Keylist, rb)
			// fmt.Println(chunk[0].Type)
			if err != nil {
				errors <- fmt.Errorf("failed to initialize subroot for chunk %d: %v", idx, err)
				return
			}

			results <- treeRoot
		}(chunk, i)
	}

	// 等待所有 Goroutine 完成
	wg.Wait()
	close(results)
	close(errors)

	// 检查是否有错误
	if len(errors) > 0 {
		return nil, <-errors
	}

	// 收集结果
	for root := range results {
		subroots = append(subroots, root)
	}

	return subroots, nil
}

// 递归根据Type构建下层单棵树
func BuildTree(owners []*indexbuilding.Owner, ibfLength int, Keylist []string, rb int) (*Node, error) {
	if len(owners) == 0 {
		return nil, fmt.Errorf("owners list is empty")
	}

	if len(owners) == 1 {
		leafNode := &Node{}
		err := leafNode.InitLeafNode(owners[0], ibfLength, Keylist, rb)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize leaf node: %v", err)
		}
		return leafNode, nil
	}

	mid := len(owners) / 2
	leftOwners := owners[:mid]
	rightOwners := owners[mid:]

	var leftNode, rightNode *Node
	var leftErr, rightErr error
	var wg sync.WaitGroup

	// 并行构建左右子树
	wg.Add(1)
	go func() {
		defer wg.Done()
		leftNode, leftErr = BuildTree(leftOwners, ibfLength, Keylist, rb)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rightNode, rightErr = BuildTree(rightOwners, ibfLength, Keylist, rb)
	}()

	// 等待左右子树构建完成
	wg.Wait()

	if leftErr != nil {
		return nil, leftErr
	}
	if rightErr != nil {
		return nil, rightErr
	}

	// 创建并初始化中间节点
	midNode := &Node{
		Left:  leftNode,
		Right: rightNode,
	}
	err := midNode.InitMidNode(ibfLength, Keylist, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mid node: %v", err)
	}

	return midNode, nil
}

// 将20个subroot节点再次合并成一棵树
// 上层数取决于Type个数没必要开协程
func CreateFinalTree(subroots []*Node, ibfLength int, Keylist []string, rb int) (*Node, error) {
	// kn := 0
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
				upMidNode := &Node{
					Left:  left,
					Right: right,
				}
				// 初始化中间节点
				err := upMidNode.InitUpMid_RootNode(ibfLength, Keylist, rb)
				// typslice := strings.Join(upMidNode.Typ, ",")
				// fmt.Println("num:", kn, "type:", typslice)
				// kn += 1
				if err != nil {
					return nil, fmt.Errorf("failed to initialize mid node: %v", err)
				}
				nextLevel = append(nextLevel, upMidNode)
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

// 实现前序遍历
func (n *Node) PreOrderTraversal(num *int, level int) {
	if n == nil {
		return
	}

	// 访问当前节点（根节点）
	fmt.Printf("%d节点%d层", *num, level)
	*num = *num + 1
	// 递归遍历左子树
	if n.Left != nil {
		n.Left.PreOrderTraversal(num, level+1)
	}

	// 递归遍历右子树
	if n.Right != nil {
		n.Right.PreOrderTraversal(num, level+1)
	}
}
