package sharedmemory

func Server(shared chan<- string, input <-chan string) {
	shared <- <-input
}