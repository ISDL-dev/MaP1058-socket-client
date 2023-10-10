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
	// TODO: サーバー側のIPアドレスを動的に取得する処理も必要
	serverIP := "192.168.10.105"
	txtAdConf := socket.SocketConfig{
		ServerIP:   serverIP,
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
		ServerIP:   serverIP,
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
		os.Exit(0)
	}()

	setting, err := txtAdapter.GetSetting()
	if err != nil {
		panic(err)
	}
	// TODO: 設定値をファイルに書き込む

	// TODO: 生波形の受信と解析データの受信を並行して行う
	// TODO: ENDコマンドを受信するまで、生波形の受信とファイル書き込みを繰り返す
	for i := 0; i < 10; i++ {
		s, err := binAdapter.ReceiveADValues(ctx)
		if err != nil {
			panic(err)
		}
		err = s.SetMeasurements(setting.Calibration)
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
		// // TODO: 計測値をファイルに書き込む
		// _, err = file.WriteString()
		// if err != nil {
		// 	panic(err)
		// }
	}

	// TODO: ENDコマンドを受信するまで、解析データの受信とファイル書き込みを繰り返す
	// TODO: 計測値をファイルに書き込む
}
