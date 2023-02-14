package main

import (
	"fmt"
	"time"
)

var (
	ch    = make(chan int)
	done  = make(chan struct{}, 1)
	token = make(chan struct{}, 1)
	wait  = make(chan struct{}, 1)
)

func main() {
	setupTokens()

	go processor()

	send()
	<-wait
	fmt.Println("main finsihed")
}

func send() {
	fmt.Println("send send started...")

	ch <- 10
	fmt.Println("send waiting...")
	<-token

	fmt.Println("send finished...")
	close(done)
	close(ch)
	close(wait)
}

func processor() {
	for {
		select {
		case n := <-ch:

			time.Sleep(3 * time.Second)
			fmt.Println("processed number: ", n*10)
			token <- struct{}{}
		case <-done:
			return
		}
	}
}

func setupTokens() {
	for i := 0; i < len(token); i++ {
		token <- struct{}{}
	}
}
