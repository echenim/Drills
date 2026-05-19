package main

import "fmt"

func lengthOfLongestSubstring(s string) int {
	lenght := 0

	for i, c := range s {
		checker := make(map[rune]int)
		checker[c] = 1
		sl := 1
		for j := i + 1; j < len(s); j++ {
			if _, exist := checker[rune(s[j])]; exist {
				break
			}
			checker[rune(s[j])] = 1
			sl++
		}
		// fmt.Println(i, " : ", string(c))
		if sl > lenght {
			lenght = sl
		}
	}

	return lenght
}

func main() {
	s := "pwwkew" //"bbbbb" //"abcabcbb"
	k := lengthOfLongestSubstring(s)
	fmt.Println(k) // 3
}
