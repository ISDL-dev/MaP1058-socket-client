package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Be3751/socket-capture-signals/internal/adapter"
	"github.com/Be3751/socket-capture-signals/internal/parser"
	"github.com/Be3751/socket-capture-signals/internal/socket"
)

func main() {
	serverIP := "192.168.10.128"
	clientIP, err := getMyLocalIP()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	txtAdConf := socket.SocketConfig{
		ServerIP:   serverIP,
		ServerPort: "3000",
		ClientIP:   clientIP,
		ClientPort: 2000,
	}
	txtAdConn, err := socket.Connect(txtAdConf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer txtAdConn.Close()
	fmt.Printf("successful to connect to the txt server!: connection config=%+v\n", txtAdConf)

	binAdConf := socket.SocketConfig{
		ServerIP:   serverIP,
		ServerPort: "2200",
		ClientIP:   clientIP,
		ClientPort: 2100,
	}
	binAdConn, err := socket.Connect(binAdConf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer binAdConn.Close()
	fmt.Printf("successful to connect to the bin server!: connection config=%+v\n", binAdConf)

	ctx := context.Background()
	txtAdapter := adapter.NewTxtAdapter(txtAdConn)
	parser := parser.NewParser(parser.ParseConfig{
		SumBytes:         1604,
		SumCheckCodeSize: 4,
	})
	binAdapter := adapter.NewBinAdapter(binAdConn, parser)

	err = txtAdapter.StartRec(ctx, time.Second*60, "2023/07/11 15-00-00")
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

func getMyLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("you might not be connected to the network")
}
