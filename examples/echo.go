package main

import (
	"fmt"
	"logx"
	"net"
	"time"
	"tron"
)

func main() {
	s := tron.NewServer("localhost:8080", serverPacketHandler)
	s.ListenAndServe()

	conn, err := dial("localhost:8080")
	if err != nil {
		logx.Error(err)
		return
	}

	cli := tron.NewClient(conn, clientPacketHandler)
	cli.Run()

	pack := tron.NewPacket(1, []byte("ping"))
	if err = cli.DirectWrite(pack); err != nil {
		logx.Error(err)
		return
	}

	time.Sleep(1 * time.Second)
}

func serverPacketHandler(worker *tron.Client, p *tron.Packet) {
	fmt.Printf("[client %s] -> [server %s]: %s\n", worker.RemoteAddr(), worker.LocalAddr(), string(p.Data))

	var data []byte
	if string(p.Data) == "ping" {
		data = []byte("pong")
	}
	resp := tron.NewPacket(p.Cate, data)
	worker.DirectWrite(resp)
}

func clientPacketHandler(cli *tron.Client, p *tron.Packet) {
	fmt.Printf("[server %s] -> [client %s]: %s\n", cli.RemoteAddr(), cli.LocalAddr(), string(p.Data))
}

func dial(address string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return conn, nil
}
