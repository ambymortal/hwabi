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

	pkt, err := onHandshake()
	if err != nil {
		log.Printf("failed to build handshake: %v", err)
		return
	}

	// send handshake to client
	_, err2 := conn.Write(pkt)
	if err2 != nil {
		log.Printf("handshake error %s: %v", addr, err2)
		return
	}
	log.Printf("sent handshake to %s", addr)

	buffer := make([]byte, 1024)

	//attempt to read what the client is sending back
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("disconnected from %s read error: %v", addr, err)
		return
	}

	log.Printf("received %d bytes from %s: %s", n, addr, hex.EncodeToString(buffer[:n]))
}

func onHandshake() ([]byte, error) {
	pkt := NewPacket()

	pkt.WriteShort(95)       // version
	pkt.WriteString("1")     // patch
	pkt.WriteInt(0x11223344) // Recv
	pkt.WriteInt(0xAABBCCDD) // Send
	pkt.WriteByte(8)         // locale

	return pkt.ToBytes()
}
