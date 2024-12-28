package query

// poof of node
type PON struct {
	Typ string

	// ======正确性=======
	// 节点hash
	HV []byte
	// 层高
	Height int

	// ======完整性=======
}
