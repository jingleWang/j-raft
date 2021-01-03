package transport

import (
	"bytes"
	"encoding/binary"
	"io"
)

// |--magic--|---id---|--length--|-flag-|----extra----|-----body-----|
type pack struct {
	magic   uint16
	version byte
	id      uint64
	length  uint32
	flag    byte
	body    []byte
}

const (
	REQ  byte = 0b00000000
	RESP byte = 0b10000000
)

func NewPack(id uint64, body []byte, flag byte) *pack {
	p := &pack{
		magic:   MagicVal,
		version: Version,
		id:      id,
		length:  uint32(len(body)),
		flag:    flag,
		body:    body,
	}
	return p
}

func Decode(reader io.Reader) *pack {
	p := &pack{}
	binary.Read(reader, binary.BigEndian, &p.magic)
	binary.Read(reader, binary.BigEndian, &p.version)
	binary.Read(reader, binary.BigEndian, &p.id)
	binary.Read(reader, binary.BigEndian, &p.length)
	binary.Read(reader, binary.BigEndian, &p.flag)

	p.body = make([]byte, p.length)
	binary.Read(reader, binary.BigEndian, &p.body)

	return p
}

func (p *pack) Encode() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &p.magic)
	binary.Write(buf, binary.BigEndian, &p.version)
	binary.Write(buf, binary.BigEndian, &p.id)
	binary.Write(buf, binary.BigEndian, &p.length)
	binary.Write(buf, binary.BigEndian, &p.flag)
	binary.Write(buf, binary.BigEndian, &p.body)

	return buf.Bytes()
}

func (p *pack) IsRequest() bool {
	return p.flag&RESP == 0
}

func (p *pack) Id() uint64 {
	return p.id
}

func (p *pack) Body() []byte {
	return p.body
}
