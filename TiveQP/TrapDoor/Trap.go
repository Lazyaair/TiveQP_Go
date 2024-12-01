package trapdoor

import (
	indexbuilding "TiveQP/IndexBuilding"
	"encoding/hex"
)

type TwinBitArray struct {
	data [2][]uint64
	cols int
}

func NewTwinBitArray(cols int) *TwinBitArray {
	uint64PerRow := (cols + 63) / 64

	tba := &TwinBitArray{
		cols: cols,
	}

	for i := 0; i < 2; i++ {
		tba.data[i] = make([]uint64, uint64PerRow)
	}

	return tba
}

func (t *TwinBitArray) Set(row, col int, value bool) {
	if row < 0 || row >= 2 || col < 0 || col >= t.cols {
		panic("index out of bounds")
	}
	uint64Index := col / 64
	bitOffset := col % 64

	if value {
		t.data[row][uint64Index] |= (1 << bitOffset)
	} else {
		t.data[row][uint64Index] &^= (1 << bitOffset)
	}
}

func (t *TwinBitArray) Get(row, col int) bool {
	if row < 0 || row >= 2 || col < 0 || col >= t.cols {
		panic("index out of bounds")
	}
	uint64Index := col / 64
	bitOffset := col % 64

	return (t.data[row][uint64Index] & (1 << bitOffset)) != 0
}

func GenT(u *indexbuilding.User, keylist []string, rd int) ([][][]string, error) {
	t := make([][][]string, 3)

	// T1=t[0] ---Type
	type_prefix, err := u.TypeEncode()
	if err != nil {
		return nil, err
	}
	t[0] = make([][]string, len(type_prefix)-1)
	for i := 0; i < len(t[0]); i++ {
		t[0][i] = make([]string, len(keylist))
		for j := 0; j < len(keylist)-1; j++ {
			outbytes := HMACSHA256([]byte(type_prefix[i]), []byte(keylist[j]))
			hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
			t[0][i][j] = hex.EncodeToString(outbytes) + "," + hex.EncodeToString(hkp1)
		}
	}

	// T2=t[1] ---Location
	location_prefix, err := u.LocationEncode()
	if err != nil {
		return nil, err
	}
	t[1] = make([][]string, len(location_prefix))
	for i := 0; i < len(t[1]); i++ {
		t[1][i] = make([]string, len(keylist)-1)
		for j := 0; j < len(keylist)-1; j++ {
			outbytes := HMACSHA256([]byte(location_prefix[i]), []byte(keylist[j]))
			hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
			t[1][i][j] = hex.EncodeToString(outbytes) + "," + hex.EncodeToString(hkp1)
		}
	}

	// T3=t[2] ---Time
	time_prefix, err := u.TimeEncode()
	if err != nil {
		return nil, err
	}
	t[2] = make([][]string, len(time_prefix))
	for i := 0; i < len(t[2]); i++ {
		t[2][i] = make([]string, len(keylist)-1)
		for j := 0; j < len(keylist)-1; j++ {
			outbytes := HMACSHA256([]byte(time_prefix[i]), []byte(keylist[j]))
			hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
			t[2][i][j] = hex.EncodeToString(outbytes) + "," + hex.EncodeToString(hkp1)
		}
	}

	return t, nil
}
