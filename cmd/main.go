package main

import (
	"context"
	"fmt"
	"net"

	"github.com/Be3751/socket-capture-signals/internal/agent"
	"github.com/Be3751/socket-capture-signals/internal/config"
	"github.com/Be3751/socket-capture-signals/internal/socket"
)

func main() {
	conf := config.Config{
		ServerIP:         "127.0.0.1",
		ServerPortText:   "3000",
		ServerPortBinary: "2000",
		ClientIP:         "127.0.0.1",
		ClientPort:       1000,
	}
	ctx := context.Background()

	conn, err := connect(ctx, conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := socket.NewSocketClient(conn)
	_ = agent.NewAgent(client)
}

func connect(ctx context.Context, conf config.Config) (*net.TCPConn, error) {
	serverAdd, err := net.ResolveTCPAddr("tcp", conf.ServerIP+":"+conf.ServerPortText)
	if err != nil {
		return nil, err
	}
	clientAdd := &net.TCPAddr{
		IP:   net.IP(conf.ClientIP),
		Port: conf.ClientPort,
	}

	// TCPネットワークの接続
	conn, err := net.DialTCP("tcp", clientAdd, serverAdd)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
