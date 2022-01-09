package pipes

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func Server(pipeFile string, input <-chan string) {
	f, err := os.OpenFile(pipeFile, os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	f.WriteString(fmt.Sprintf("%s\n", <-input))
}

// create named pipes file
func CreateFlie(pipeFile string) {
	os.Remove(pipeFile)
	err := syscall.Mkfifo(pipeFile, 0666)
	if err != nil {
		log.Fatal("create named pipe error:", err)
	}
}
