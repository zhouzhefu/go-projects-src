package main

import (
	"fmt"
	"net/rpc"
	"os"
)

// duck typing args type, as long as fields matching (contains of) the remote struct, name is not a problem
type SomeArgs struct {
	A, B int
	C string
}

func main() {
	client, err := rpc.DialHTTP("tcp", ":8989")
	checkError(err)

	someArgs := SomeArgs{A:17, B:4, C:"ok"}
	var reply int
	err = client.Call("Arith1.Multiply", someArgs, &reply)
	checkError(err)

	fmt.Println("Result from server is:", reply)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}