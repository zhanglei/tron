package main

import (
	"fmt"
	"net"
	"time"
	"tron"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	clientConf := tron.NewDefaultConf(1 * time.Minute)
	r := tron.NewReconnectTaskManager(5*time.Second, 3)
	manager := tron.NewClientsManager(r)
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		panic(err)
	}

	codec := tron.NewDefaultCodec()
	cli := tron.NewClient(conn, clientConf, codec, packHandler)
	cli.ReadWriteAndHandle()
	g := tron.NewClientsGroup("add-service", "add-service")
	manager.Add(g, cli)

	go func() {
		for i := 0; i < 5; i++ {
			pingPack := tron.NewReqPacket([]byte("ping"))

			// 异步写
			// cli.AsyncWrite(pingPack)

			// 同步写
			_, err := cli.SyncWrite(pingPack, 2*time.Second)
			if err != nil {
				fmt.Printf("client sync write %v failed: %v\n", pingPack, err)
				return
			}
			// fmt.Println("res", string(res.([]byte)))
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(100000 * time.Second)
}

func packHandler(cli *tron.Client, p *tron.Packet) {
	fmt.Printf("[server:%s] -> [client:%s]: %s\n",
		tron.SplitPort(cli.RemoteAddr()),
		tron.SplitPort(cli.LocalAddr()),
		string(p.Data))
	// fmt.Println("clientPackHandler, resp: ", string(p.Data))
	cli.NotifyReceived(p.Header.Seq, p.Data) // ok
}
