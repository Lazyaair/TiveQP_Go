package query

// E:\Github\TiveQP\TiveQP\TiveQP\Cachemidnode 0.txt
// k = 10
// eve = Restaurants**ATLANTA**33.846335**-84.3635778**7**0**21**0
// Query start
// Query end
// Verification :
// Result number：10
// Correctness :
// Restaurants**ATLANTA**33.846335**-84.3635778**7**0**21**0
// Restaurants**ATLANTA**33.8428749**-84.3785385**10**0**17**0
// Restaurants**ATLANTA**33.8275823**-84.328604**6**0**22**0
// Restaurants**ATLANTA**33.84269152**-84.37045581**4**0**23**0
// Restaurants**ATLANTA**33.847331**-84.372713**11**0**21**0
// Restaurants**ATLANTA**33.84703603**-84.36566228**11**30**21**0
// Restaurants**ATLANTA**33.820152**-84.387432**11**30**21**0
// Restaurants**ATLANTA**33.8475491**-84.3737958**1**0**20**0
// Restaurants**ATLANTA**33.8167245**-84.3356715**10**0**15**30
// Restaurants**ATLANTA**33.89221993**-84.32607651**10**0**22**0
// true
// true
// Completeness :
// true
// time_query = 577017
// time_proof = 5
// Verification time: 359623
// proof size ：
// 13948
import (
	construction "TiveQP/Construction"
	trapdoor "TiveQP/Trapdoor"
)

func QueryT(root *construction.Node, td *trapdoor.T, k *int) {
	if root == nil {
		return
	}
	if *k == 0 {
		return
	}

	// 判断根节点的左子节点是否符合条件P
	if root.Left != nil && check(root.Left, td) {
		// 如果左子节点符合条件P，继续判断左子节点的左子节点
		QueryT(root.Left, td, k)
	} else {
		// 如果左子节点不符合条件P，判断根的右子节点
		if root.Right != nil {
			QueryT(root.Right, td, k)
		}
	}
}

func check(root *construction.Node, t *trapdoor.T) bool {
	if root.YCS != nil {
		// 只关 T1
		return true
	} else {
		// 只关 T2+T3
		return true
	}
}
