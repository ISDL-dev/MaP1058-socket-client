package main

import (
	"context"
	"encoding/csv"
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
	serverIP := "192.168.10.101"
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

	ctx, cancel := context.WithCancel(context.Background())
	parser := parser.NewParser()
	scanner := scanner.NewCustomScanner(txtAdConn)
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, scanner, parser)
	// binAdapter := adapter.NewBinAdapter(binAdConn, parser)
	eegWriter := csv.NewWriter(os.Stdout)
	writerGroup := adapter.CSVWriterGroup{
		EEGWriter: *eegWriter,
	}
	defer func() {
		eegWriter.Flush()
	}()

	err = txtAdapter.StartRec(ctx, time.Second*60, time.Now())
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	err = txtAdapter.EndRec(ctx)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	os.Exit(0)
	// }()

	// TODO: parseする前の受信コマンドを実際に確認する（parserにバグがある可能性）
	setting, err := txtAdapter.GetSetting()
	if err != nil {
		panic(err)
	}
	// TODO: 設定値をファイルに書き込む
	for {
		tmp := make([]byte, 128)
		_, err := txtAdConn.Read(tmp)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(tmp))
	}

	// TODO: 生波形の受信と解析データの受信を並行して行う
	// TODO: ENDコマンドを受信するまで、生波形の受信とファイル書き込みを繰り返す
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	go func() {
		fmt.Println("start to get trend data")
		err := txtAdapter.GetTrendData(ctx, writerGroup, setting.AnalysisType)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second * 5)

	// TODO: ENDコマンドを受信するまで、解析データの受信とファイル書き込みを繰り返す
	// TODO: 計測値をファイルに書き込む
}
