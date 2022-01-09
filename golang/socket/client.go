package socket

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func TCPClient(addr string, ready chan<- string, fin chan<- string) {
	var conn net.Conn
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)
	
	// keep connecting until connection is established
	for {
		conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err == nil {
			break
		}
	}

	// tell main process client1 is ready and read from server
	ready       <- "Client1 is ready"
	result, err := ioutil.ReadAll(conn)
	checkError(err)

	// get mean value
	nums := toNums(string(result))
	mean := calcMean(nums)

	// print result and tell main process client is finished
	fmt.Println("Mean is", mean)
	fin <- "Client1 is finished"
}

// calculate mean value
func calcMean(nums []int) float64 {
	sum := 0
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
	}
	return (float64(sum)) / (float64(len(nums)))
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