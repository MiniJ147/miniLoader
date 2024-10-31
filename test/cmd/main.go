package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	mk := []int{}
	fmt.Println("working dir", mk)
	fmt.Println(os.Getwd())
	x := 5
	fmt.Printf("hello from test program time is %v and this is an example int value %d", time.Now().Unix(), x)

	y := 0
	// fmt.Scanf("%d", &y)
	fmt.Println(y)
	b, e := os.ReadFile("cmd/example.txt")
	fmt.Println(string(b), e)
	// mk[0]
	for {
		fmt.Printf("alive fromtest the live rseloaer =)\n")
		time.Sleep(time.Second)
		// fmt.Println(mk[0])
	}
}
