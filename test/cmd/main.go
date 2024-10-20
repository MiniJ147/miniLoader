package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println(os.Getwd())
	x := 5
	fmt.Printf("hello from test program time is %v and this is an example int value %d", time.Now().Unix(), x)
	for {
		///
		fmt.Printf("alive\n")
		time.Sleep(2 * time.Second)
	}
}
