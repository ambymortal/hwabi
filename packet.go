package main

import (
	"bytes"
	"encoding/binary"
)

// error tracking here could be best further in the project
type Packet struct {
	buff *bytes.Buffer
	err  error
}

func NewPacket() *Packet {
	return &Packet{buff: new(bytes.Buffer)}
}

// WriteByte writes a single byte.
func (p *Packet) WriteByte(val byte) {
	if p.err != nil {
		return
	}
	p.err = p.buff.WriteByte(val)
}

func (p *Packet) WriteShort(val uint16) {
	if p.err != nil {
		return
	}
	p.err = binary.Write(p.buff, binary.LittleEndian, val)
}

func (p *Packet) WriteInt(val uint32) {
	if p.err != nil {
		return
	}
	p.err = binary.Write(p.buff, binary.LittleEndian, val)
}

// length prefixed string
func (p *Packet) WriteString(val string) {
	if p.err != nil {
		return
	}
	p.WriteShort(uint16(len(val)))
	if _, err := p.buff.WriteString(val); err != nil && p.err == nil {
		p.err = err
	}
}

func (p *Packet) WriteBytes(val []byte) {
	if p.err != nil {
		return
	}
	if _, err := p.buff.Write(val); err != nil && p.err == nil {
		p.err = err
	}
}

// returns the final byte slice and error
func (p *Packet) ToBytes() ([]byte, error) {
	return p.buff.Bytes(), p.err
}
