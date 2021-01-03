package transport

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type Handler func([]byte) []byte

type server struct {
	listener net.Listener

	Handler Handler

	sign chan struct{}
}

func NewServerWithListener(listener net.Listener, handler Handler) *server {
	return &server{
		listener: listener,
		Handler:  handler,
	}
}

func (s *server) Start() error {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go func() {
			fmt.Printf("accecpt connection from %v\n", conn.RemoteAddr())
			defer func() {
				fmt.Println("lost connection")
			}()
			scanner := bufio.NewScanner(conn)

			HandleReceive(scanner, func(buf []byte) {
				ch := make(chan *response)

				go func() {
					for {
						select {
						case resp := <-ch:
							pack := NewPack(resp.id, resp.body, RESP)
							conn.Write(pack.Encode())
						}
					}
				}()

				go func() {
					pack := Decode(bytes.NewReader(buf))
					if !pack.IsRequest() {
						return
					}
					resp := s.Handler(pack.body)

					ch <- &response{
						id:   pack.id,
						body: resp,
					}
				}()

			})

		}()
	}
}

func (s *server) Close() error {
	s.listener.Close()
	return nil
}
