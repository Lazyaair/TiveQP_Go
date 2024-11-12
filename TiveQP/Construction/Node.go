package construction

// DownNode 表示下层树的节点
type DownNode struct {
	Data []string
	LCS  []string
	TCS  []string
	// IBF的插入位置，HMAC==>bigint=(mod)=>index+'|'+0/1
	Bits_LCS [][]string
	Bits_TCS [][]string
	// Bits的摘要
	HV_LCS [][]byte
	HV_TCS [][]byte
	// 原始数据的摘要
	HV []byte
	// 子节点
	Left  *DownNode
	Right *DownNode
}

// UpNode 表示上层树的节点
type UpNode struct {
	Data []string
	LCS  []string
	TCS  []string
	YCS  []string
	// IBF的插入位置，HMAC==>bigint=(mod)=>index+'|'+0/1
	Bits_LCS [][]string
	Bits_TCS [][]string
	Bits_YCS [][]string
	// Bits的摘要
	HV_LCS [][]byte
	HV_TCS [][]byte
	HV_YCS [][]byte
	// 子节点的合并摘要
	HV []byte
	// 子节点
	Left  *UpNode
	Right *UpNode
}
