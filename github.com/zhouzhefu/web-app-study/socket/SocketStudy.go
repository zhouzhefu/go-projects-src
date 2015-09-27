package main

import (
	"fmt"
	"net"
	"os"
	"io/ioutil"
	"time"
)

func main() {
	ip := net.ParseIP("192.168.1.254")
	fmt.Println(ip)
	fmt.Println(ip.DefaultMask())

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989") //no ip mean localhost or 127.0.0.1
	checkError(err)
	fmt.Println("tcpAddr:", tcpAddr)

	// basicSingleConn(tcpAddr)
	// singleEventLoopThreadConn(tcpAddr)
	multiConcurrentLongConn(tcpAddr)

	os.Exit(0)
}

func multiConcurrentConn(tcpAddr *net.TCPAddr) {
	ch := make(chan int, 5)

	for i:=0;i<cap(ch);i++ {
		go chanSingleConn(tcpAddr, ch, i)
	}

	for idx := range ch {
		fmt.Println("Goroutine No.", idx, "ending. len(ch):", len(ch))

		// it is NOT the perfect way to make sure all idx have been written into channel, 
		// in which idx=4 maybe later pushed in than the idx=5 one. 
		if idx == 4 {close(ch)}
	}
}

func multiConcurrentLongConn(tcpAddr *net.TCPAddr) {
	ch := make(chan int, 5)

	for i:=0;i<cap(ch);i++ {
		go chanSingleLongConn(tcpAddr, ch, i)
	}

	for idx := range ch {
		fmt.Println("Goroutine No.", idx, "ending. ")

		// it is NOT the perfect way to make sure all idx have been written into channel, 
		// in which idx=4 maybe later pushed in than the idx=5 one. 
		if idx == 4 {close(ch)}
	}
}

func singleEventLoopThreadConn(tcpAddr *net.TCPAddr) {
	for i:=0;i<5;i++ {
		singleConn(tcpAddr)
	}
}

func basicSingleConn(tcpAddr *net.TCPAddr) {
	singleConn(tcpAddr)
}

// channel writing must be placed AFTER the actual network communication, otherwise 
// the close(ch) will not wait till the communication finished. 
func chanSingleConn(tcpAddr *net.TCPAddr, ch chan int, idx int) {
	singleConn(tcpAddr)
	ch <- idx
}

func chanSingleLongConn(tcpAddr *net.TCPAddr, ch chan int, idx int) {
	singleLongConn(tcpAddr)
	ch <- idx
}

func singleConn(tcpAddr *net.TCPAddr) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)

	// cannot use ReadAll() when long connection is required, it will hang the goroutine until 
	// an remote EOF is received (e.g. server side close the connection). 
	result, err := ioutil.ReadAll(conn) 
	checkError(err)
	fmt.Println(string(result))
	conn.Close()
}

func singleLongConn(tcpAddr *net.TCPAddr) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	err = conn.SetKeepAlive(true)
	checkError(err)

	for i:=0; i<5; i++ {
		_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		checkError(err)

		response := make([]byte, 128)
		readLen, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error found in client:", err)
		}
		fmt.Println(string(response), readLen)

		time.Sleep(500 * time.Millisecond)
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}
}