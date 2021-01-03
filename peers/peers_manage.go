package peers

import (
	"fmt"
	"net"
	"sync"
)

type peersManage struct {
	mutex sync.Mutex
	conf  PeerManageConf
	peers map[int64]*peer
}

func (pm *peersManage) Init() error {
	err := pm.startServ()
	if err != nil {
		return err
	}

	pm.tryConnectMembers()

	return nil
}

func (pm *peersManage) startServ() error {
	ln, err := net.Listen("tcp", fmt.Sprint(":%v", pm.conf.Port))
	if err != nil {
		return err
	}

	//// todo handle conn
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		pm.handleConn(conn)
	}
	return nil
}

func (pm *peersManage) handleConn(conn net.Conn) {

}

func (pm *peersManage) tryConnectMembers() {

	//for id, addr := range pm.conf.Members {
	//	conn, err := net.Dial("tcp", addr)
	//	if err != nil {
	//		continue
	//	}
	//
	//	conn.
	//}

}

type PeerManageConf struct {
	Id      int64
	Port    int64
	Members map[int64]string
}

type peer struct {
	Id   int64
	addr *net.Addr
	Conn *net.Conn
}

//func (p *peer)
