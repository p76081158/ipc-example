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
	nums   := toNums(result)
	mode   := calcMode(nums)

	// print result and tell main process client is finished
	fmt.Println("Mode is", mode)
	fin <- "Client3 is finished"
}

// calculate mode value
func calcMode(nums []int) int {
	countMap := make(map[int]int)
	max      := 0
	mode     := 0

	for _, key := range nums {
		countMap[key] += 1
	}
	for _, key := range nums {
		freq := countMap[key]
		if freq > max {
			mode = key
			max  = freq
		}
	}
	return mode
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