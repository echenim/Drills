package main

import "fmt"

func twoSum(nums []int, target int) []int {
	mp := make(map[int]int)

	var sum []int

	for index, value := range nums {
		mp[value] = index
	}

	for j, num := range nums {
		compliment := target - num
		i, exist := mp[compliment]
		if exist && i != j {
			sum = append(sum, j, i)
			fmt.Println(sum)
			break
		}

	}

	return sum
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	fmt.Println(twoSum(nums, target)) // [2, 7]
}
