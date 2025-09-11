package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8484")
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

	defer listener.Close()
	fmt.Println("accepting client connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}

		c := NewClientConnection(conn)
		go c.handleConnection()
	}
}
