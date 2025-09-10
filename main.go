package main

import (
	"encoding/hex"
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	log.Printf("new client connected from: %s", addr)

	var handshake = []byte{
		0x0D, 0x00, // length
		0x5F, 0x00, // version
		0x01,                   // patch string
		0x31,                   // "1"
		0x11, 0x22, 0x33, 0x44, // send IV
		0xAA, 0xBB, 0xCC, 0xDD, // recv IV
		0x08, // locale
	}

	_, err := conn.Write(handshake)
	if err != nil {
		log.Printf("handshake error %s: %v", addr, err)
		return
	}

	log.Printf("sent handshake to %s", addr)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("disconnected from %s read error: %v", addr, err)
		return
	}

	log.Printf("received %d bytes from %s: %s", n, addr, hex.EncodeToString(buffer[:n]))
}
