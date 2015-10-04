package main 

import (
	"fmt"
	"time"
)

func counting(c chan<- int) {
	for i:=0; i<4; i++ {
		time.Sleep(1 * time.Second)
		c <- i
	}
	close(c)
}

func main() {
	msg := "Starting main"
	fmt.Println(msg)

	slice := []int{3, 5, 7, 9}
	for _, v := range slice {
		fmt.Println(v)
	}
	bus := make(chan int)
	msg = "starting a goroutine"
	go counting(bus)
	for count := range bus {
		fmt.Println("count:", count)
	}
}
