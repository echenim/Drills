package main

import (
	"fmt"
	"sort"
)

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	sum := append(nums1, nums2...)
	sort.Ints(sum)

	if len(sum)%2 == 0 {
		return float64((sum[(len(sum)/2)-1] + sum[(len(sum)/2)]) / 2)
	}

	return float64(sum[len(sum)/2])
}

func main() {
	nums1 := []int{1, 3}
	nums2 := []int{2, 3}

	k := findMedianSortedArrays(nums1, nums2)
	fmt.Println(k)
}
