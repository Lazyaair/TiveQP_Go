package query

// poof of node
type PON struct {
	Typ string

	// ======正确性=======
	// 节点hash
	HV []byte
	// 层高
	Height int

	// ======完备性=======

	Bits_LCS [][]string
	Bits_TCS [][]string
	Bits_YCS [][]string

	HV_LCS [][][]byte
	HV_TCS [][][]byte
	HV_YCS [][][]byte
}
