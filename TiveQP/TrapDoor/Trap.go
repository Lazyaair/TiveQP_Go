package trapdoor

import (
	indexbuilding "TiveQP/IndexBuilding"
	"encoding/hex"
)

type T struct {
	T1 [][]string
	T2 [][]string
	T3 [][]string
}

func GenT(u *indexbuilding.User, keylist []string, rd int) (*T, error) {

	// T1 --- Type
	type_prefix, err := u.TypeEncode()
	if err != nil {
		return nil, err
	}
	t1 := make([][]string, len(type_prefix))
	for i := 0; i < len(t1); i++ {
		t1[i] = make([]string, len(keylist)-1)
		for j := 0; j < len(keylist)-1; j++ {
			outbytes := HMACSHA256([]byte(type_prefix[i]), []byte(keylist[j]))
			hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
			t1[i][j] = hex.EncodeToString(outbytes) + "," + hex.EncodeToString(hkp1)
		}
	}

	// T2 --- Location
	location_prefix, err := u.LocationEncode()
	if err != nil {
		return nil, err
	}
	t2 := make([][]string, len(location_prefix))
	for i := 0; i < len(t2); i++ {
		t2[i] = make([]string, len(keylist)-1)
		for j := 0; j < len(keylist)-1; j++ {
			outbytes := HMACSHA256([]byte(location_prefix[i]), []byte(keylist[j]))
			hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
			t2[i][j] = hex.EncodeToString(outbytes) + "," + hex.EncodeToString(hkp1)
		}
	}

	// T3 --- Time
	time_prefix, err := u.TimeEncode()
	if err != nil {
		return nil, err
	}
	t3 := make([][]string, len(time_prefix))
	for i := 0; i < len(t3); i++ {
		t3[i] = make([]string, len(keylist)-1)
		for j := 0; j < len(keylist)-1; j++ {
			outbytes := HMACSHA256([]byte(time_prefix[i]), []byte(keylist[j]))
			hkp1 := HashSHA256(append(outbytes, []byte(keylist[len(keylist)-1])...))
			t3[i][j] = hex.EncodeToString(outbytes) + "," + hex.EncodeToString(hkp1)
		}
	}
	td := &T{
		T1: t1,
		T2: t2,
		T3: t3,
	}
	return td, nil
}
