package main

import (
	"fmt"
	"time"
)

type Token struct {
	Data      string
	Recipient int
	TTL       int
}

func node(tokenChannel <-chan Token, nextChannel chan<- Token, id int) {
	for {

		token := <-tokenChannel   //get token from the channel
		if token.Recipient == 0 { // if message reach recipient
			fmt.Printf("Received message: %s\n", token.Data)
		} else if token.TTL > 0 { // reporting from goroutines
			fmt.Printf("Received message: %s, forwarding to node %d\n", token.Data, id+1)
			token.Recipient--
			token.TTL--
			nextChannel <- token
		} else {
			fmt.Printf("Message TTL expired for message: %s\n", token.Data)
		}
	}
}

func main() {
	var numNodes int
	fmt.Print("Enter the number of nodes: ")
	fmt.Scanln(&numNodes)

	var ptr = make([]chan Token, numNodes+1) //creating slice for channels
	for i := range ptr {
		ptr[i] = make(chan Token) //creating channels
	}

	// creating numNodes goroutines
	for i := 0; i < numNodes; i++ {

		go node(ptr[i], ptr[i+1], i)
	}

	message := Token{Data: "Hello, World!", Recipient: 10, TTL: 100}
	ptr[0] <- message // Sending message from the main flow

	// sleep so message have time to reach recipient
	time.Sleep(time.Second * 5)
}
