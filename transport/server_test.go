package transport

import (
	"fmt"
	"net"
	"testing"
)

func TestNewServerWithListener(t *testing.T) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", 2532))
	if err != nil {
		t.Fatal(err)
	}
	sev := NewServerWithListener(ln, func(bytes []byte) []byte {
		fmt.Printf("recv: %v\n", string(bytes))
		return []byte(fmt.Sprintf("serv recv: %v", string(bytes)))
	})

	sev.Start()
}
