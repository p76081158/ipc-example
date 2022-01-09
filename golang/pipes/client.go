package pipes

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"sort"
	"strings"
	"strconv"
)

func Client(pipeFile string, ready chan<- string, fin chan<- string) {
	ready  <- "Client2 is ready"
	result := read(pipeFile,fin)

	// get median value
	nums   := toNums(string(result))
	median := caclcMedian(nums)
	
	// print result and tell main process client is finished
	fmt.Println("Median is", median)
	fin <- "Client2 is finished"
}

func read(pipeFile string, fin chan<- string) string {
	file, err := os.OpenFile(pipeFile, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	reader    := bufio.NewReader(file)
	line, err := reader.ReadBytes('\n')
	return string(line)
}

func caclcMedian(nums []int) float64 {
	sort.Ints(nums)
	mid := len(nums) / 2
	if len(nums) % 2 == 0 {
		return float64((nums[mid-1] + nums[mid])) / 2.0
	} else {
		return float64(nums[mid])
	}
}

// split string to nums
func toNums(str string) []int {
	var nums []int
	tmp  := strings.Split(str, "\n")
	strs := strings.Split(tmp[0], " ")
	for i := 0; i < len(strs); i++ {
		value, err :=strconv.Atoi(strs[i])
		if err == nil {
			nums = append(nums, value)
		}
	}
	return nums
}