package hwabinet

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"
)

type ClientSession struct {
	conn net.Conn
	send []byte
	recv []byte
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	log.Printf("new client connected from: %s", addr)

	session := &ClientSession{
		conn: conn,
		send: make([]byte, 4),
		recv: make([]byte, 4),
	}

	rand.Read(session.send)
	rand.Read(session.recv)

	pkt, err := onHandshake(session.recv, session.send)
	if err != nil {
		log.Printf("failed to build handshake: %v", err)
		return
	}

	// send handshake to client
	_, err2 := session.conn.Write(pkt)
	if err2 != nil {
		log.Printf("handshake error %s: %v", addr, err2)
		return
	}
	log.Printf("sent handshake to %s", addr)

	buffer := make([]byte, 1024)

	// keep connection open
	for {
		n, err := session.conn.Read(buffer)
		if err != nil {
			log.Printf("disconnected from %s read error: %v", addr, err)
			return
		}

		data := buffer[:n]

		// handle crypto next

		log.Printf("received %d bytes from %s: %s",
			n, addr, hex.EncodeToString(data))
	}
}

func onHandshake(recv, send []byte) ([]byte, error) {
	pkt := NewPacket()

	pkt.WriteShort(95)                             // version
	pkt.WriteString("1")                           // patch
	pkt.WriteInt(binary.LittleEndian.Uint32(recv)) // Recv
	pkt.WriteInt(binary.LittleEndian.Uint32(send)) // Send
	pkt.WriteByte(8)                               // locale

	return pkt.ToBytes()
}
