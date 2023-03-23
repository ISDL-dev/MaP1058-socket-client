package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Be3751/socket-capture-signals/internal/agent"
	"github.com/Be3751/socket-capture-signals/internal/config"
	"github.com/Be3751/socket-capture-signals/internal/pkg/socket"
)

// TODO: 受信コマンドをParseする関数も必要
// TODO: クライアント側のIPアドレスを自動取得する処理も必要
func main() {
	conf := config.Config{
		ServerIP:         "127.0.0.1",
		ServerPortText:   "3000",
		ServerPortBinary: "2000",
		ClientIP:         "127.0.0.1",
		ClientPort:       1000,
	}
	ctx := context.Background()

	conn, err := socket.NewClient(conf)
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
	err = agent.CaptureSig(ctx, buf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(buf)

	agent.EndRec(ctx)
}
