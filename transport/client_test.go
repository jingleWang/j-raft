package transport

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"testing"
	"time"
)

func TestClient_Send(t *testing.T) {

	conn, _ := net.Dial("tcp", "127.0.0.1:2532")

	//time.Sleep(time.Second * 5)
	c := NewClient(&conn)
	c.Init()
	conn.Close()
	resp, err := c.Send([]byte("hello wrold1"))
	if err == net.ErrWriteToConnected {

		fmt.Println(err)
	}
	resp, _ = c.Send([]byte("hello wrold2"))
	fmt.Println(err)
	resp, _ = c.Send([]byte("hello wrold3"))
	fmt.Println(err)
	resp, _ = c.Send([]byte("hello wrold4"))
	fmt.Println(resp)

	time.Sleep(time.Second * 10)
}

func TestHandleReceive(t *testing.T) {
	pack := NewPack(math.MaxUint64, []byte("hello world12345678"), RESP)

	writer := new(bytes.Buffer)

	binary.Write(writer, binary.BigEndian, pack.Encode())
	binary.Write(writer, binary.BigEndian, pack.Encode())
	binary.Write(writer, binary.BigEndian, pack.Encode())

	reader := bytes.NewReader(writer.Bytes())
	HandleReceive(bufio.NewScanner(reader), func(buf []byte) {

		p := Decode(bytes.NewReader(buf))
		fmt.Println(p)

	})

}
