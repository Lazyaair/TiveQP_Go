package main

import (
	construction "TiveQP/Construction"
	query "TiveQP/Query"
	resultverification "TiveQP/Resultverification"
	trapdoor "TiveQP/Trapdoor"
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	ibfLength := 200000
	Keylist := []string{"2938879577741549", "8729598049525437", "8418086888563864", "0128636306393258", "2942091695121238", "6518873307787549"}
	rb := 235648

	filename := "E:\\Github\\TiveQP\\TiveQP\\TiveQP\\Data\\20k.txt" // 文件名
	owners, err := construction.LoadOwners(filename)
	if err != nil {
		fmt.Println("加载 Owner 数据出错:", err)
		return
	}
	subroots, err := construction.BuildTreesByChunks(owners, ibfLength, Keylist, rb)
	if err != nil {
		fmt.Println("Error building subroots:", err)
	} else {
		fmt.Println("Subroots built successfully!")
	}
	finalRoot, err := construction.CreateFinalTree(subroots, ibfLength, Keylist, rb)
	if err != nil {
		fmt.Println("Error creating final tree:", err)
	} else {
		fmt.Println("Final tree created successfully!")
	}
	// "Fast Food**AUSTIN**30.2795878**-97.806248**12**12"
	// Restaurants**ATLANTA**33.846335**-84.3635778**12**12
	// Shopping**AUSTIN**30.3575044**-97.7321061**10**0**19**0
	u, err := trapdoor.ParseUser("Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("User loaded successfully!==Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
	}
	T, err := trapdoor.GenT(u, Keylist, rb, 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("TrapDoor created successfully!==Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
	}
	// u1, err := trapdoor.ParseUser("Shopping**AUSTIN**30.3575044**-97.7321061**11**11")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("User loaded successfully!==Shopping**AUSTIN**30.3575044**-97.7321061**11**11")
	// }
	// T1, err := trapdoor.GenT(u1, Keylist, rb)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("TrapDoor created successfully!==Shopping**AUSTIN**30.3575044**-97.7321061**11**11")
	// }

	fmt.Println("Query begin!===Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
	k := 3
	result := make([]*[]byte, 0, k)
	pi := make([]*query.PON, 0, k)
	query.QueryT(finalRoot, T, &k, 0, rb, &result, &pi)
	fmt.Println("Query ended!===Restaurants**ATLANTA**33.846335**-84.3635778**12**12")
	// fmt.Println(pi)
	for _, v := range result {
		p, _ := construction.Decrypt(*v, []byte("2bc73dw20ebf4d46"))
		fmt.Println(string(p))
		// fmt.Println("\n========================================")
	}
	// for _, v := range pi {
	// 	fmt.Println(v)
	// 	// fmt.Println("\n========================================")
	// }
	fmt.Println("======================================================================")
	fmt.Println("check HV==", resultverification.CheckHV(finalRoot.HV, pi))
	fmt.Println("======================================================================")
	fmt.Println("check Completeness==", resultverification.CheckCompleteness(T, pi))

	// fmt.Println("Query begin!===Shopping**AUSTIN**30.3575044**-97.7321061**11**11")
	// k1 := 2
	// result1 := make([]*construction.Node, 0, k1)
	// pi1 := make([]*query.PON, 0, k)
	// query.QueryT(finalRoot, T1, &k1, 0, rb, &result1, &pi1)
	// fmt.Println("Query ended!===Shopping**AUSTIN**30.3575044**-97.7321061**11**11")
	// fmt.Println(pi1)

	// for _, v := range result {
	// 	print(v)
	// 	// fmt.Println("\n========================================")
	// }
	// number := 0
	// level := 0
	// finalRoot.PreOrderTraversal(&number, level)
}
