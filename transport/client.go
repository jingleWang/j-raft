package transport

import (
	"bufio"
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type client struct {
	conn *net.Conn

	timeout time.Duration

	requestChn chan *request

	//request
	rid *uint64

	mutex        sync.Mutex
	responseChns map[uint64]chan *response
}

func NewClient(conn *net.Conn) *client {
	startReqId := uint64(1)
	c := &client{
		conn:         conn,
		timeout:      DefaultTimeOut,
		requestChn:   make(chan *request),
		rid:          &startReqId,
		responseChns: make(map[uint64]chan *response),
	}

	return c
}

func (c *client) Init() error {
	send := func(req *request) {

		pack := NewPack(req.id, req.body, REQ)

		defer func() {
			err := recover()
			if err != nil {
				resp := &response{
					id:   pack.Id(),
					body: nil,
					err:  err.(error),
				}

				c.responseChns[pack.Id()] <- resp
			}
		}()

		//fmt.Println(pack)
		_, err := (*c.conn).Write(pack.Encode())
		if err != nil {
			resp := &response{
				id:   pack.Id(),
				body: nil,
				err:  err,
			}
			c.responseChns[pack.Id()] <- resp
		}
	}

	go func() {
		for {
			select {
			case r := <-c.requestChn:
				send(r)
			}
		}

	}()

	go func() {
		scanner := bufio.NewScanner(*c.conn)

		HandleReceive(scanner, func(buf []byte) {
			pack := Decode(bytes.NewReader(buf))
			if pack.IsRequest() {
				return
			}
			resp := &response{
				id:   pack.Id(),
				body: pack.Body(),
			}
			c.responseChns[pack.Id()] <- resp
		})

	}()

	return nil
}

func (c *client) Send(body []byte) ([]byte, error) {
	req := &request{
		id:   c.getNextReqId(),
		body: body,
	}
	ch := make(chan *response)
	c.addRspChn(req.id, ch)
	defer c.rmRspChn(req.id)

	//tick := time.Tick(c.timeout)

	c.requestChn <- req

	select {
	case rsp := <-ch:
		return rsp.body, rsp.err
		//case <-tick:
		//return nil, errors.New("timed out")
	}
}

func (c *client) getNextReqId() uint64 {
	return atomic.AddUint64(c.rid, 1)
}

func (c *client) addRspChn(id uint64, ch chan *response) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.responseChns[id] = ch
}

func (c *client) rmRspChn(id uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.responseChns, id)
}
