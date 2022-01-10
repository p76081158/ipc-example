package sharedmemory

import (
	"fmt"
	"strconv"
	"strings"
)

func Client(shared <-chan string, ready chan<- string, fin chan<- string) {
	ready  <- "Client3 is ready"
	result := <-shared

	// get mode value
	nums        := toNums(result)
	mode, exist := calcMode(nums)

	// print result and tell main process client is finished
	if exist {
		fmt.Println("Mode is", mode)
	}else {
		fmt.Println("Mode is not existed")
	}
	fin <- "Client3 is finished"
}

// calculate mode value
func calcMode(nums []int) ([]int, bool) {
	var ans []int
	countMap := make(map[int]int)
	max      := 0
	mode     := 0
	freq     := 0
	for _, key := range nums {
		countMap[key] += 1
	}
	for _, key := range nums {
		freq = countMap[key]
		if freq > max {
			mode = key
			max  = freq
		}
	}
	// Pigeonhole principle
	if max > (len(nums) / len(countMap) + 1) {
		// single mode value
		return append(ans, mode), true
	} else {
		if max == (len(nums) / len(countMap)) && (len(nums) % len(countMap)) == 0 {
			// mode value not existed
			return ans, false
		} else {
			// single or multiple mode values
			for key, value := range countMap {
				if max == value {
					ans = append(ans, key)
				}
			}
			return ans, true
		}
	}
}

// split string to nums
func toNums(str string) []int {
	var nums []int
	strs := strings.Split(str, " ")
	for i := 0; i < len(strs); i++ {
		value, err :=strconv.Atoi(strs[i])
		if err == nil {
			nums = append(nums, value)
		}
	}
	return nums
}