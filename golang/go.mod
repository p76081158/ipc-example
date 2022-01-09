module github.com/p76081158/ipc-example/golang

go 1.14

replace (
	github.com/p76081158/ipc-example/golang/socket => ./socket
	github.com/p76081158/ipc-example/golang/pipes => ./pipes
	github.com/p76081158/ipc-example/golang/sharedmemory => ./sharedmemory
)