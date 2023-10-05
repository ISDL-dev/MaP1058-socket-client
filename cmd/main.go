package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/adapter"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/scanner"
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
	scanner := scanner.NewCustomScanner(txtAdConn)
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, scanner, parser)
	binAdapter := adapter.NewBinAdapter(binAdConn, parser)

	err = txtAdapter.StartRec(ctx, time.Second*60, time.Now())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = txtAdapter.EndRec(ctx)
		if err != nil {
			panic(err)
		}
	}()

	setting, err := txtAdapter.GetSetting()
	if err != nil {
		panic(err)
	}

	file, err := os.Create("temp.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// ENDコマンドを受信するまで受信とファイル書き込みを繰り返す
	for i := 0; i < 10; i++ {
		s, err := binAdapter.ReceiveADValues(ctx)
		if err != nil {
			panic(err)
		}
		err = s.SetMeasurements(setting.Calibration)
		if err != nil {
			panic(err)
		}
		// TODO: 計測値をファイルに書き込む
		_, err = file.WriteString()
		if err != nil {
			panic(err)
		}
	}

}
