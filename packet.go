package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

// credits to valhalla (Hucaru)

// Packet -
type Packet []byte

// NewPacket -
func NewPacket() Packet {
	return make(Packet, 0)
}

// Append -
func (p *Packet) Append(data []byte) {
	*p = append(*p, data...)
}

// Size -
func (p *Packet) Size() int {
	return int(len(*p))
}

// String -
func (p Packet) String() string {
	return fmt.Sprintf("[Packet] (%d) : % X", len(p), string(p))
}

// WriteByte -
func (p *Packet) WriteByte(data byte) {
	*p = append(*p, data)
}

// WriteInt8 - byte
func (p *Packet) WriteInt8(data int8) {
	*p = append(*p, byte(data))
}

// WriteInt -
func (p *Packet) WriteInt(data int) {
	*p = append(*p, byte(data), byte(data>>8), byte(data>>16), byte(data>>24))
}

// WriteBool -
func (p *Packet) WriteBool(data bool) {
	if data {
		*p = append(*p, 0x1)
	} else {
		*p = append(*p, 0x0)
	}
}

// WriteUint16 - short
func (p *Packet) WriteUint16(data uint16) {
	*p = append(*p, byte(data), byte(data>>8))
}

// WriteUint32 - int
func (p *Packet) WriteUint32(data uint32) {
	*p = append(*p, byte(data), byte(data>>8), byte(data>>16), byte(data>>24))
}

// WriteUint64 - long
func (p *Packet) WriteUint64(data uint64) {
	*p = append(*p, byte(data), byte(data>>8), byte(data>>16), byte(data>>24),
		byte(data>>32), byte(data>>40), byte(data>>48), byte(data>>56))
}

// WriteFloat32 -
func (p *Packet) WriteFloat32(data float32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(data))
	*p = append(*p, b[:]...)
}

// WriteBytes -
func (p *Packet) WriteBytes(data []byte) {
	p.Append(data)
}

// WriteString -
func (p *Packet) WriteString(str string) {
	p.WriteUint16(uint16(len(str)))
	p.WriteBytes([]byte(str))
}

// WritePaddedString -
func (p *Packet) WritePaddedString(str string, number int) {
	if len(str) > number {
		p.WriteBytes([]byte(str)[:number])
	} else {
		p.WriteBytes([]byte(str))
		p.WriteBytes(make([]byte, number-len(str)))
	}
}

// WriteInt16 -
func (p *Packet) WriteInt16(data int16) { p.WriteUint16(uint16(data)) }

// WriteInt32 -
func (p *Packet) WriteInt32(data int32) { p.WriteUint32(uint32(data)) }

// WriteInt64 -
func (p *Packet) WriteInt64(data int64) { p.WriteUint64(uint64(data)) }
