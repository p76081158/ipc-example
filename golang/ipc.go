package main

import (
	"github.com/p76081158/ipc-example/golang/socket"
	"github.com/p76081158/ipc-example/golang/pipes"
	"github.com/p76081158/ipc-example/golang/sharedmemory"
	"bufio"
	"fmt"
	"os"
)

// send stdin to Servers
func producer(channel chan string, client_num int, input string) {
	for i := 0; i < client_num; i++ {
		channel <- input
	}
}

// wait client process finish
func waitClient(channels []chan string) {
	for i := 0; i < len(channels); i++ {
		fmt.Println(<-channels[i])
	}
}

// close channel
func closeChannel(channels ...chan string) {
	for _, channel := range channels {
		close(channel)
	}
}

func main() {
	pipeFile     := "pipes/namedpipe.ipc"
	mess         := make(chan string)
	shared       := make(chan string)
	client_num   := 3
	rdy_channels := make([]chan string,client_num)
	out_channels := make([]chan string,client_num)
	for i := 0; i < client_num; i++ {
		rdy_channels[i] = make(chan string)
		out_channels[i] = make(chan string)
	}
	pipes.CreateFlie(pipeFile)

	// create Servers and Clients goroutines ( T = 0 )
	go socket.TCPServer(":7070", mess)
	go socket.TCPClient(":7070", rdy_channels[0], out_channels[0])
	go pipes.Server(pipeFile, mess)
	go pipes.Client(pipeFile, rdy_channels[1], out_channels[1])
	go sharedmemory.Server(shared, mess)
	go sharedmemory.Client(shared, rdy_channels[2], out_channels[2])
	
	// wait for clients ready ( T = 1 )
	waitClient(rdy_channels)
	fmt.Println("Server is ready. You can type intergers and then click [ENTER].  Clients will show the mean, median, and mode of the input values.")
	// stdin
	consolescanner := bufio.NewScanner(os.Stdin)
	consolescanner.Scan()
	input := consolescanner.Text()
	producer(mess, client_num, input)
	
	// wait for clients finish and close channels ( T = 2 )
	waitClient(out_channels)
	closeChannel(mess)
	closeChannel(rdy_channels...)
	closeChannel(out_channels...)
}