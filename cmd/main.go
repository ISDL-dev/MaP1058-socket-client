package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/adapter"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
	"github.com/Be3751/MaP1058-socket-client/utils/net"
)

// TODO: クライアント側のIPアドレスを自動取得する処理も必要
func main() {
	clientIP, err := net.GetMyLocalIP()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	txtAdConf := socket.SocketConfig{
		ServerIP:   "192.168.10.101",
		ServerPort: "3000",
		ClientIP:   clientIP,
		ClientPort: 1100,
	}
	txtAdConn, err := socket.Connect(txtAdConf)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = txtAdConn.Close()
		if err != nil {
			panic(err)
		}
	}()
	binAdConf := socket.SocketConfig{
		ServerIP:   "192.168.10.101",
		ServerPort: "2200",
		ClientIP:   clientIP,
		ClientPort: 1200,
	}
	binAdConn, err := socket.Connect(binAdConf)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = binAdConn.Close()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	parser := parser.NewParser()
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, parser)
	binAdapter := adapter.NewBinAdapter(binAdConn, parser)

	err = txtAdapter.StartRec(ctx, time.Second*60, time.Now())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer func() {
		err = txtAdapter.EndRec(ctx)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()

	// 10 x 50 ポイントのAD値を受信する
	for i := 0; i < 10; i++ {
		s, err := binAdapter.ReceiveADValues(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
	}

}
