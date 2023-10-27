package main

import (
	"context"
	"fmt"
	"os"
	"sync"
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

	ctx, cancel := context.WithCancel(context.Background())
	parser := parser.NewParser()
	scanner := scanner.NewCustomScanner(txtAdConn)
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, scanner, parser)
	binAdapter := adapter.NewBinAdapter(binAdConn, parser)
	csvWriterGroup := adapter.CSVWriterGroup{}

	err = txtAdapter.StartRec(time.Second*60, time.Now())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = txtAdapter.EndRec()
		if err != nil {
			panic(err)
		}
	}()

	setting, err := txtAdapter.GetSetting()
	if err != nil {
		panic(err)
	}
	// TODO: 設定値をファイルに書き込む

	go func() {
		err := binAdapter.WriteRawSignal(ctx, os.Stdout)
		if err != nil {
			panic(err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := txtAdapter.WriteTrendData(ctx, csvWriterGroup, setting.AnalysisType)
		if err != nil {
			panic(err)
		}
	}()
	// ENDコマンドを送信するまで待機（ENDコマンドを送信する処理が別途必要）
	wg.Wait()
	cancel()
}
