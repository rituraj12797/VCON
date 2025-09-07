package engine

import "fmt"

// currentl we will pass the version x ====> Array of sentence identifier
//								   x+1 ====> Array of sentence identifier

// 							LCS = dmallest identifier subsequence common in both

/* CONVERT THIS TO ITERATIVE DP FOR AVOIDING STACK MEMORY EXPLOSION IN LARGER DOCUMENTS */
func recursion(ind1 int, ind2 int, arr *[]string, brr *[]string, dp *[][]int) int {

	// prune base case
	if ind1 >= len(*arr) || ind2 >= len(*brr) {
		return 0
	}

	// cache miss
	val := (*dp)[ind1][ind2]

	if val != -1 {
		return (*dp)[ind1][ind2]
	}

	var ans int = 0

	if (*arr)[ind1] == (*brr)[ind2] {

		ans = 1 + recursion(ind1+1, ind2+1, arr, brr, dp)
	} else {
		ans = max(recursion(ind1+1, ind2, arr, brr, dp), recursion(ind1, ind2+1, arr, brr, dp))
	}

	(*dp)[ind1][ind2] = ans
	return ans
}

func lcsGenerator(ind1 int, ind2 int, arr *[]string, brr *[]string, dp *[][]int, result *[]string) {

	if ind1 >= len(*arr) || ind2 >= len(*brr) {
		return
	}

	val := (*dp)[ind1][ind2]

	if (*arr)[ind1] == (*brr)[ind2] {
		// ans = max(1+recursion(ind1+1,ind2+1,arr,brr,dp), max(recursion(ind1+1,ind2,arr,brr,dp),recursion(ind1,ind2+1,arr,brr,dp)));      coule be one of the 3 values
		(*result) = append((*result), (*arr)[ind1])
		lcsGenerator(ind1+1, ind2+1, arr, brr, dp, result)
	} else {
		if val == recursion(ind1+1, ind2, arr, brr, dp) {
			lcsGenerator(ind1+1, ind2, arr, brr, dp, result)
		} else {
			lcsGenerator(ind1, ind2+1, arr, brr, dp, result)
		}
	}

}

func LCS(version_x1 *[]string, version_x2 *[]string) []string {

	var dp [][]int

	for i := 0; i < len(*version_x1)+1; i++ {
		var temp []int
		for j := 0; j < len(*version_x2)+1; j++ {
			temp = append(temp, -1)
		}
		dp = append(dp, temp)
	}

	recursion(0, 0, version_x1, version_x2, &dp)
	var res []string
	lcsGenerator(0, 0, version_x1, version_x2, &dp, &res)

	fmt.Println(" The LCS : ", res)

	return res, nil
}
