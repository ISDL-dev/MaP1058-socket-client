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

func main() {
	clientIP, err := net.GetMyLocalIP()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

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
	defer func() {
		err = txtAdapter.EndRec(ctx)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("end recording")

		err = txtAdConn.Close()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("close txt connection")

		err = binAdConn.Close()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("close bin connection")
		os.Exit(0)
	}()

	// TODO: parseする前の受信コマンドを実際に確認する（parserにバグがある可能性）
	setting, err := txtAdapter.GetSetting()
	if err != nil {
		panic(err)
	}
	fmt.Println(setting)
	for {
		tmp := make([]byte, 128)
		n, err := txtAdConn.Read(tmp)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(tmp[:n]))
	}

}
