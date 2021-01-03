package transport

import (
	"bufio"
	"encoding/binary"
)

func HandleReceive(scanner *bufio.Scanner, f func(buf []byte)) {
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && len(data) >= 16 && binary.BigEndian.Uint16(data[:2]) == MagicVal {
			var l = binary.BigEndian.Uint32(data[11:15])
			pl := int(HeaderLength) + int(l)
			if pl <= len(data) {
				return pl, data[:pl], nil
			}
		}
		return
	})

	for scanner.Scan() {
		buf := scanner.Bytes()
		f(buf)
	}
}

type request struct {
	id   uint64
	body []byte
}

type response struct {
	id   uint64
	body []byte
	err  error
}
