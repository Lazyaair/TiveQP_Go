package construction

import (
	indexbuilding "TiveQP/Indexbuilding"
	"fmt"
	"strconv"
	"strings"
)

// TwinBitArray 表示一个基于 []uint64 的位数组，用来表示IBF
type TwinBitArray struct {
	data [2][]uint64 // 固定为 2 行，每行是一个 uint64 数组
	Cols int         // 每行的位数
}

// NewBitArray 创建指定大小的位数组
func NewTwinBitArray(cols int) *TwinBitArray {
	// 每行所需的 uint64 数量
	uint64PerRow := (cols + 63) / 64

	// 初始化数据
	tba := &TwinBitArray{
		Cols: cols,
	}

	// 为两行分别分配空间
	for i := 0; i < 2; i++ {
		tba.data[i] = make([]uint64, uint64PerRow)
	}

	return tba
}

// Set 设置位数组中的某一位
func (t *TwinBitArray) Set(row, col int, value bool) {
	if row < 0 || row >= 2 || col < 0 || col >= t.Cols {
		panic("index out of bounds")
	}
	uint64Index := col / 64 // 计算对应的 uint64 索引
	bitOffset := col % 64   // 计算该 uint64 中的位偏移

	if value {
		t.data[row][uint64Index] |= (1 << bitOffset) // 设置为 1
	} else {
		t.data[row][uint64Index] &^= (1 << bitOffset) // 清除该位
	}
}

// Get 获取位数组中的某一位
func (t *TwinBitArray) Get(row, col int) bool {
	if row < 0 || row >= 2 || col < 0 || col >= t.Cols {
		panic("index out of bounds")
	}
	uint64Index := col / 64 // 计算对应的 uint64 索引
	bitOffset := col % 64   // 计算该 uint64 中的位偏移

	return (t.data[row][uint64Index] & (1 << bitOffset)) != 0
}

// Node 表示树的节点
type Node struct {
	// ==== 测试字段 =====
	Typ []string

	// ==== 测试字段 ====

	Owner *indexbuilding.Owner
	// IBF
	IBF *TwinBitArray
	// 补集
	LCS []string
	TCS []string
	YCS []string
	// IBF的插入位置，HMAC==>bigint=(mod)=>index+'|'+0/1
	Bits_LCS [][]string
	Bits_TCS [][]string
	Bits_YCS [][]string
	// Bits的摘要
	HV_LCS [][][]byte
	HV_TCS [][][]byte
	HV_YCS [][][]byte
	// 原始数据的摘要
	HV []byte
	// 密文
	Cipher []byte
	// 子节点
	Left  *Node
	Right *Node
}

//	func (n *Node) Print() {
//		typslice := strings.Join(n.Typ, ",")
//		fmt.Println("Type:", typslice, "Owner:", n.Owner)
//	}
func (n *Node) Print() {
	if n == nil {
		fmt.Println("Node is nil")
		return
	}

	// 输出基本信息
	typslice := strings.Join(n.Typ, ",")
	fmt.Println("=== Node Information ===")
	fmt.Printf("Type: [%s]\n", typslice)
	if n.Owner != nil {
		fmt.Printf("Owner: %+v\n", *n.Owner)
	} else {
		fmt.Println("Owner: nil")
	}

	// 输出补集信息
	fmt.Printf("LCS: %v\n", n.LCS)
	fmt.Printf("TCS: %v\n", n.TCS)
	fmt.Printf("YCS: %v\n", n.YCS)

	// 输出 Bits 信息
	fmt.Printf("Bits_LCS: %v\n", n.Bits_LCS)
	fmt.Printf("Bits_TCS: %v\n", n.Bits_TCS)
	fmt.Printf("Bits_YCS: %v\n", n.Bits_YCS)

	// 输出 HV 信息
	fmt.Printf("HV_LCS: %v\n", n.HV_LCS)
	fmt.Printf("HV_TCS: %v\n", n.HV_TCS)
	fmt.Printf("HV_YCS: %v\n", n.HV_YCS)
	fmt.Printf("HV: %v\n", n.HV)

	// 输出密文
	fmt.Printf("Cipher: %v\n", n.Cipher)
	fmt.Println("====================")
}

// 下层叶节点初始化
func (dln *Node) InitLeafNode(owner *indexbuilding.Owner, ibf_length int, Keylist []string, rb int) error {
	dln.Owner = owner
	dln.IBF = NewTwinBitArray(ibf_length)
	// 关于 Location
	// 取位置编码
	locationcode, err := owner.LocationEncode()
	if err != nil {
		return fmt.Errorf("Location编码失败")
	}
	// Location插入IBF
	for i := 0; i < len(locationcode); i++ {
		Insert(dln.IBF, locationcode[i], Keylist, rb)
	}
	// 取位置补集
	dln.LCS, err = owner.LocationComplementEncode()
	if err != nil {
		return fmt.Errorf("Location补集编码失败")
	}
	dln.Bits_LCS = make([][]string, len(dln.LCS))
	dln.HV_LCS = make([][][]byte, len(dln.LCS))
	// 处理location补集！
	for i := 0; i < len(dln.LCS); i++ {
		dln.Bits_LCS[i] = make([]string, len(Keylist)-1)
		dln.HV_LCS[i] = make([][]byte, len(Keylist)-1)
		InsertCS(dln.IBF, dln.LCS[i], &dln.Bits_LCS[i], Keylist, &dln.HV_LCS[i], rb)
	}

	// 关于 Time
	// 取时间编码
	timecode, err := owner.TimeEncode()
	if err != nil {
		return fmt.Errorf("Time编码失败")
	}
	// Time插入IBF
	for i := 0; i < len(timecode); i++ {
		Insert(dln.IBF, timecode[i], Keylist, rb)
	}
	// 取时间补集
	dln.TCS, err = owner.LocationComplementEncode()
	if err != nil {
		return fmt.Errorf("Time补集编码失败")
	}
	dln.Bits_TCS = make([][]string, len(dln.TCS))
	dln.HV_TCS = make([][][]byte, len(dln.TCS))
	// 处理Time补集！
	for i := 0; i < len(dln.TCS); i++ {
		dln.Bits_TCS[i] = make([]string, len(Keylist)-1)
		dln.HV_TCS[i] = make([][]byte, len(Keylist)-1)
		InsertCS(dln.IBF, dln.TCS[i], &dln.Bits_TCS[i], Keylist, &dln.HV_TCS[i], rb)
	}

	// 计算节点HASH
	dataText := owner.Type + "**" + owner.City +
		"**" + strconv.FormatFloat(owner.Lat, 'f', 6, 64) +
		"**" + strconv.FormatFloat(owner.Lng, 'f', 6, 64) +
		"**" + strconv.Itoa(owner.HourStart) +
		"**" + strconv.Itoa(owner.MinStart) +
		"**" + strconv.Itoa(owner.HourClose) +
		"**" + strconv.Itoa(owner.MinClose)
	ciphertext, err := Encrypt([]byte(dataText), []byte("2bc73dw20ebf4d46"))
	dln.Cipher = ciphertext
	if err != nil {
		return fmt.Errorf("加密失败")
	}
	dln.HV = HashSHA256(ciphertext)
	// 无关字段
	dln.Left = nil
	dln.Right = nil
	dln.YCS = nil
	dln.Bits_YCS = nil
	dln.HV_YCS = nil
	return nil
}

// 下层中间节点初始化
func (dmn *Node) InitMidNode(ibf_length int, Keylist []string, rb int) error {
	// 合并 IBF
	dmn.IBF = OrIBF(dmn.Left.IBF, dmn.Right.IBF)
	// 取补集并集
	dmn.LCS = MergeSet(dmn.Left.LCS, dmn.Right.LCS)
	dmn.TCS = MergeSet(dmn.Left.TCS, dmn.Right.TCS)

	dmn.Bits_LCS = make([][]string, len(dmn.LCS))
	dmn.HV_LCS = make([][][]byte, len(dmn.LCS))
	// 处理location补集！
	for i := 0; i < len(dmn.LCS); i++ {
		dmn.Bits_LCS[i] = make([]string, len(Keylist)-1)
		dmn.HV_LCS[i] = make([][]byte, len(Keylist)-1)
		InsertCS(dmn.IBF, dmn.LCS[i], &dmn.Bits_LCS[i], Keylist, &dmn.HV_LCS[i], rb)
	}
	dmn.Left.LCS = nil
	dmn.Right.LCS = nil

	dmn.Bits_TCS = make([][]string, len(dmn.TCS))
	dmn.HV_TCS = make([][][]byte, len(dmn.TCS))
	// 处理Time补集！
	for i := 0; i < len(dmn.TCS); i++ {
		dmn.Bits_TCS[i] = make([]string, len(Keylist)-1)
		dmn.HV_TCS[i] = make([][]byte, len(Keylist)-1)
		err := InsertCS(dmn.IBF, dmn.TCS[i], &dmn.Bits_TCS[i], Keylist, &dmn.HV_TCS[i], rb)
		if err != nil {
			return fmt.Errorf("IBF:%v,TCS:%v,BITTCS:%v", dmn.IBF, dmn.TCS[i], dmn.Bits_TCS[i])
		}
	}
	dmn.Left.TCS = nil
	dmn.Right.TCS = nil
	// 计算节点HASH
	dmn.HV = HashSHA256(append(dmn.Left.HV, dmn.Right.HV...))
	// 无关字段
	dmn.YCS = nil
	dmn.Bits_YCS = nil
	dmn.HV_YCS = nil
	dmn.Owner = nil
	dmn.Cipher = nil
	return nil
}

// 上层叶节点初始化
func (uln *Node) InitUpLeafNode(typ string, ibf_length int, Keylist []string, rb int) error {
	// 上层叶节点生成 BuildTree 会调用 InitMidNode 初始化
	// if uln.Bits_LCS != nil && uln.Bits_TCS != nil && uln.HV_LCS != nil && uln.HV_TCS != nil && uln.LCS != nil && uln.TCS != nil {
	// 	fmt.Println(true)
	// }
	// 非所需初始化，清空当前内容
	uln.IBF = NewTwinBitArray(ibf_length)
	uln.LCS = nil
	uln.Bits_LCS = nil
	uln.HV_LCS = nil
	uln.TCS = nil
	uln.Bits_TCS = nil
	uln.HV_TCS = nil
	uln.Cipher = nil

	uln.Typ = append(uln.Typ, typ)
	// 取 TypeCode
	typecode, err := indexbuilding.TypeEncoding(typ)
	if err != nil {
		return fmt.Errorf("TypeCoding Error")
	}
	// 插入TypeCode
	for i := 0; i < len(typecode); i++ {
		Insert(uln.IBF, typecode[i], Keylist, rb)
	}

	// 取Typecode补集
	uln.YCS, err = indexbuilding.TypeEncodingComplement(typ)
	if err != nil {
		return fmt.Errorf("TypeCSCoding Error")
	}
	// 处理type补集！
	uln.Bits_YCS = make([][]string, len(uln.YCS))
	uln.HV_YCS = make([][][]byte, len(uln.YCS))
	for i := 0; i < len(uln.YCS); i++ {
		uln.Bits_YCS[i] = make([]string, len(Keylist)-1)
		uln.HV_YCS[i] = make([][]byte, len(Keylist)-1)
		InsertCS(uln.IBF, uln.YCS[i], &uln.Bits_YCS[i], Keylist, &uln.HV_YCS[i], rb)
	}

	// 计算节点HASH
	uln.HV = HashSHA256(append(uln.Left.HV, uln.Right.HV...))
	// 无关字段
	uln.Owner = nil
	return nil
}

// 上层中间和根节点初始化
func (mrn *Node) InitUpMid_RootNode(ibf_length int, Keylist []string, rb int) error {
	// 合并 IBF
	mrn.IBF = OrIBF(mrn.Left.IBF, mrn.Right.IBF)
	// 取补集并集
	mrn.YCS = MergeSet(mrn.Left.YCS, mrn.Right.YCS)
	mrn.Typ = MergeSet(mrn.Left.Typ, mrn.Right.Typ)
	// 处理type补集！
	mrn.Bits_YCS = make([][]string, len(mrn.YCS))
	mrn.HV_YCS = make([][][]byte, len(mrn.YCS))
	for i := 0; i < len(mrn.YCS); i++ {
		mrn.Bits_YCS[i] = make([]string, len(Keylist)-1)
		mrn.HV_YCS[i] = make([][]byte, len(Keylist)-1)
		InsertCS(mrn.IBF, mrn.YCS[i], &mrn.Bits_YCS[i], Keylist, &mrn.HV_YCS[i], rb)
	}

	// 计算节点HASH
	mrn.HV = HashSHA256(append(mrn.Left.HV, mrn.Right.HV...))
	// 无关字段
	mrn.Owner = nil
	mrn.LCS = nil
	mrn.LCS = nil
	mrn.Bits_LCS = nil
	mrn.Bits_TCS = nil
	mrn.HV_LCS = nil
	mrn.HV_TCS = nil
	mrn.Cipher = nil
	return nil
}
