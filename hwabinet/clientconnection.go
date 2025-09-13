package hwabinet

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
)

type ClientConnection struct {
	conn net.Conn
	recv []byte
	send []byte
}

func HandleConnection(c net.Conn) {
	defer c.Close()

	addr := c.RemoteAddr().String()
	log.Printf("new client connected from: %s", addr)

	sendIv := make([]byte, 4)
	recvIv := make([]byte, 4)
	rand.Read(sendIv)
	rand.Read(recvIv)

	cc := &ClientConnection{
		conn: c,
		recv: recvIv,
		send: sendIv,
	}

	var handshake = clientHandshake(cc.recv, cc.send)

	// write back to client without any applied crypto
	err := cc.sendPacket(handshake)
	if err != nil {
		log.Printf("handshake error %s: %v", addr, err)
		return
	}

	log.Printf("sent handshake to %s", addr)

	// continiously try to read response from client
	for {
		buffer := make([]byte, 1024)
		n, err := cc.conn.Read(buffer) //replace with readPacket
		if err != nil {
			log.Printf("disconnected from %s read error: %v", addr, err)
			return
		}

		log.Printf("received %d bytes from %s: %s", n, addr, hex.EncodeToString(buffer[:n]))
	}
}

func (cc *ClientConnection) writePacket() {
	// write and handle encryption
}

func (cc *ClientConnection) readPacket() {
	// read and handle encryption
}

// use when sending a packet without aes and shanda encryption applied
// really just to send the handshake packet (afaik)
func (cc *ClientConnection) sendPacket(p Packet) error {
	n, err := cc.conn.Write(p)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	log.Printf("wrote %d bytes to client: %s", n, hex.EncodeToString(p))
	return nil
}

func clientHandshake(r, s []byte) Packet {
	p := NewPacket()

	p.WriteInt16(14)
	p.WriteInt16(95)   // version
	p.WriteString("1") // patch
	p.Append(r)        // recv
	p.Append(s)        // send
	p.WriteByte(8)     // locale

	return p
}
