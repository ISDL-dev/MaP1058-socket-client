package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/adapter"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// TODO: クライアント側のIPアドレスを自動取得する処理も必要
func main() {
	txtAdConf := socket.SocketConfig{
		ServerIP:   "192.168.86.24",
		ServerPort: "3000",
		ClientIP:   "192.168.86.21",
		ClientPort: 1000,
	}
	txtAdConn, err := socket.Connect(txtAdConf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	binAdConf := socket.SocketConfig{
		ServerIP:   "192.168.86.24",
		ServerPort: "2200",
		ClientIP:   "192.168.86.21",
		ClientPort: 1000,
	}
	binAdConn, err := socket.Connect(binAdConf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	ctx := context.Background()
	txtAdapter := adapter.NewTxtAdapter(txtAdConn)
	parser := parser.NewParser(parser.ParseConfig{
		SumBytes:         1604,
		SumCheckCodeSize: 4,
	})
	binAdapter := adapter.NewBinAdapter(binAdConn, parser)

	err = txtAdapter.StartRec(ctx, time.Second*60, time.Now())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for i := 0; i < 50; i++ {
		signals, err := binAdapter.ReceiveADValues(ctx)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println(signals)
	}

	err = txtAdapter.EndRec(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}
