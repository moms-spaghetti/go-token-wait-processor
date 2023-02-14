package main

import (
	"fmt"
	"time"
)

var (
	data  = make(chan int)         // data used to send to processor routine from main
	done  = make(chan struct{}, 1) // empty chan if closed kills processor
	token = make(chan struct{}, 1) // holds tokens to be read when needing to stop executing a function
	wait  = make(chan struct{}, 1) // empty chan if close kills main routine
)

func main() {
	setupTokens() // add tokens to token chan

	go processor() // start processor separate routine

	send()
	<-wait // blocks main
	fmt.Println("main finsihed")
}

func send() {
	fmt.Println("send send started...")

	data <- 10 // send int to data chan
	fmt.Println("send waiting...")
	<-token // blocks send() processing as struct has been read and token chan is empty (reading from empty chan = block)

	fmt.Println("send finished...")
	close(done) // kill processor first so not reading from still open data chan
	close(data) // now close data
	close(wait) // close wait to allow main routine to continue processing
}

func processor() {
	for {
		select {
		case n := <-data: // receives into processor from send

			time.Sleep(3 * time.Second)             // simulate processing
			fmt.Println("processed number: ", n*10) // print something n * 10
			token <- struct{}{}                     // put struct back in token chan to allow send() to continue processing
		case <-done:
			return
		}
	}
}

func setupTokens() { // adds tokens to token chan based on token chan length
	for i := 0; i < len(token); i++ {
		token <- struct{}{}
	}
}
