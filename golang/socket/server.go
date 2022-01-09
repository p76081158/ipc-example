package socket

import (
	"fmt"
	"net"
	"os"
)

// create TCP Server
func TCPServer(addr string, input <-chan string) {
	tcpAddr, err  := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// handle multiple client access
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		
		go handleClient(conn, input)
	}
}

// handle each client access
func handleClient(conn net.Conn, input <-chan string) {
	defer conn.Close()
	output := <-input
	conn.Write([]byte(output))
}

// check error message
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}