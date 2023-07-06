package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Be3751/socket-capture-signals/internal/agent"
	"github.com/Be3751/socket-capture-signals/internal/pkg/socket"
)

// TODO: 受信コマンドをParseする関数も必要
// TODO: クライアント側のIPアドレスを自動取得する処理も必要
func main() {
	conf := socket.Config{
		ServerIP:   "192.168.86.24",
		ServerPort: "3000",
		ClientIP:   "192.168.86.21",
		ClientPort: 1000,
	}
	ctx := context.Background()

	conn, err := socket.Connect(conf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	client := socket.NewSocketClient(conn)
	agent := agent.NewAgent(client)

	err = agent.StartRec(ctx, time.Second, "2023/01/01")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	buf := make([]byte, 1024)
	err = agent.StartRec(ctx, time.Minute, "2023/04/17 12:12:12")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(buf)

	agent.EndRec(ctx)
}
