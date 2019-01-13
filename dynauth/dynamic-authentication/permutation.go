/*
	Title:	Permutations package
	Author:	Connor Peters
	Date:	2/11/2018
	Desc:	To generate a limited number of permutations (i.e. 1234, 5678, 78910 from 10P4)
			Call LimPerms with a slice of strings to permutate
*/

package dynauthcore

import (
	"strings"
)

// LimPerms - to create limited permutations using subsets and heap's algorithm.
// Needs an slice of strings to create the limited permutations from and an int to limit the number of permutations.
// For example, (["1", "2", "3", "4"], 3) will create 4P3 (4*3*2) permutations.
func LimPerms(toPerm []string, num int) []string {
	subsets := getSubsets(toPerm, num)
	perms := getPerms(subsets)
	return perms
}

// Generates each unique subset of the array that is passed in, with the num being the limiting factor
func getSubsets(locks []string, num int) [][]string {
	var helper func([]string, []string, int)
	res := [][]string{}
	subset := []string{}

	helper = func(arr []string, subset []string, n int) {
		if len(subset) == num {
			tmp := make([]string, len(subset))
			copy(tmp, subset)
			res = append(res, tmp)
			return
		}

		if n == len(arr) {
			return
		}

		x := arr[n]
		subset = append(subset, x)
		// recursion starts here
		helper(locks, subset, n+1)
		subset = subset[:len(subset)-1]
		helper(locks, subset, n+1)
	}
	helper(locks, subset, len(subset))
	return res
}

// Simple go implementation of heap's algorithm
func heapPermutation(arr []string) []string {
	var helper func([]string, int)
	res := []string{}

	helper = func(arr []string, n int) {
		if n == 1 {
			tmp := strings.Join(arr, "")
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

// This is the organizer for heap's algo
// It takes a 2d slice of each subset (generated above) and returns a single array
// With each of heap's permutations as a single string with no delimiter (joined)
func getPerms(perms [][]string) []string {
	res := []string{}
	for i := 0; i < len(perms); i++ {
		tmp := []string{}
		list := perms[i]
		tmp = heapPermutation(list)
		for j := 0; j < len(tmp); j++ {
			res = append(res, tmp[j])
		}
	}
	return res
}
