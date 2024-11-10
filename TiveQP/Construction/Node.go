package construction

// LeafNode 表示下层树的节点
type LeafNode struct {
	Value int
	Left  *LeafNode
	Right *LeafNode
}

// TypeNode 表示上层树的节点
type TypeNode struct {
}
