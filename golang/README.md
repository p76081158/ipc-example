# IPC - Golang

IPC via Socket、Pipes、Shared Memory
* [Inter-Process Communication](https://vcx1127.notion.site/Inter-Process-Communication-4be174067b024db69160f362f3f31668)

## Program Flow Chart

![](https://i.imgur.com/KRHQqfY.jpg)

* [rawImage](https://i.imgur.com/KRHQqfY.jpg)
* 總共 7 個 processes (Main、3 Servers、3 Clients)
* 透過 channel 達到同步 (T = 1 and 2)

## Main

```go=
func main() {
    pipeFile     := "pipes/namedpipe.ipc"
    mess         := make(chan string)
    shared       := make(chan string)
    rdy_channels := make([]chan string,client_num)
    out_channels := make([]chan string,client_num)
    
    ...
    
}
```
* **pipeFile :** Linux 中透過 FIFO 存取 File System 中的檔案 ( namedpipe.ipc ) 來實作 Named Pipes。
* **mess :** 為 Main() 與 Servers 之間溝通之 Channel，Main 將 stdin 放入 mess 中通知 Servers。
* **rdy_channels :** 為告知 Main()，Clients 是否 Ready。
* **out_channels :** 通知 Main() 所有 Clients 皆執行完畢。
```go=11
func main() {
    
    ...
    
    // create Servers and Clients goroutines ( T = 0 )
    go socket.TCPServer(":7070", mess)
    go socket.TCPClient(":7070", rdy_channels[0], out_channels[0])
    go pipes.Server(pipeFile, mess)
    go pipes.Client(pipeFile, rdy_channels[1], out_channels[1])
    go sharedmemory.Server(shared, mess)
    go sharedmemory.Client(shared, rdy_channels[2], out_channels[2])
    
    ...
    
}
```
* 在 T = 0 時間點開啟 Server 與 Client goroutines
```go=26
func main() {
    
    ...
    
    // wait for clients ready ( T = 1 )
    waitClient(rdy_channels)
    fmt.Println("Server is ready. You can type intergers and then click [ENTER].  Clients will show the mean, median, and mode of the input values.")
    // stdin
    consolescanner := bufio.NewScanner(os.Stdin)
    consolescanner.Scan()
    input := consolescanner.Text()
    producer(mess, client_num, input)
    
    ...
    
}
```
* rdy_channels 為同步作用 ( 在 T = 1 時同步)，表示所有 Client 皆 Ready
* stdin 完後，透過 producer 通知每個 Server 的 mess Channel
```go=42
func main() {
    
    ...
    
    // wait for clients finish and close channels ( T = 2 )
    waitClient(out_channels)
    closeChannel(mess)
    closeChannel(rdy_channels...)
    closeChannel(out_channels...)
}
```
* out_channels 為同步作用 ( 在 T = 2 時同步)，表示所有 Client 皆 Finished
* Main() 最後要將 Channels 都 Close

## Socket

### TCPServer

```go=
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
```
* Blocked until receiving input ( main function stdin )
* Send input to Client by conn.Wirite
* Close connection

### TCPClient

```go=
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
```
* 透過 ready 通知 Main() Client 1 Ready
* Read input from ioutil.ReadAll(conn)
* 最後藉由 fin 通知 Main() Client 1 Finished

## Pipes

### Server

```go=
func Server(pipeFile string, input <-chan string) {
    f, err := os.OpenFile(pipeFile, os.O_RDWR, 0777)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    f.WriteString(fmt.Sprintf("%s\n", <-input))
}
```
* Blocked until receiving input ( main function stdin )
* Write input to pipeFile

### Client

```go=
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
```
* 透過 ready 通知 Main() Client 2 Ready
* Read input from pipeFile
* 最後藉由 fin 通知 Main() Client 2 Finished

## Shared Memory

* Server 與 Client 之間透過 Shared Memory 溝通 (shared chan)
* channel 有分方向
    * **chan<- :** Write Only
    * **<-chan :** Read Only

### Server

```go=
func Server(shared chan<- string, input <-chan string) {
    shared <- <-input
}
```
* Blocked until receiving input ( main function stdin )

### Client

```go=
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
```
* 透過 ready 通知 Main() Client 3 Ready
* Blocked until receiving shared ( main function stdin which is passed by Server )
* 最後藉由 fin 通知 Main() Client 3 Finished

## How to Use

```bash
git clone https://github.com/p76081158/ipc-example.git
cd ipc-example/golang
./golang
```
![](https://i.imgur.com/zWtJpDi.gif)
