package hwabinet

import (
	"encoding/hex"
	"log"
	"net"
)

type ClientConnection struct {
	conn net.Conn
}

func NewClientConnection(c net.Conn) *ClientConnection {
	cc := &ClientConnection{
		conn: c,
	}

	return cc
}

func (cc *ClientConnection) HandleConnection() {
	defer cc.conn.Close()

	addr := cc.conn.RemoteAddr().String()
	log.Printf("new client connected from: %s", addr)

	var handshake = clientHandshake()

	_, err := cc.conn.Write(handshake)
	if err != nil {
		log.Printf("handshake error %s: %v", addr, err)
		return
	}

	log.Printf("sent handshake to %s", addr)

	buffer := make([]byte, 1024)
	n, err := cc.conn.Read(buffer)
	if err != nil {
		log.Printf("disconnected from %s read error: %v", addr, err)
		return
	}

	log.Printf("received %d bytes from %s: %s", n, addr, hex.EncodeToString(buffer[:n]))
}

func clientHandshake() Packet {
	p := NewPacket()

	p.WriteInt16(14)
	p.WriteInt16(95)                         // version
	p.WriteString("1")                       // patch
	p.Append([]byte{0x11, 0x22, 0x33, 0x44}) // recv
	p.Append([]byte{0xAA, 0xBB, 0xCC, 0xDD}) // send
	p.WriteByte(8)                           // locale

	return p
}
