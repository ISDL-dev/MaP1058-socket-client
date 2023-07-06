// an echo server to check the connection with a client
package main

import (
	"fmt"
	"net"
	"os"
	"time"

	utilsNet "github.com/Be3751/socket-capture-signals/utils/net"
)

func main() {
	myIP, err := utilsNet.GetMyLocalIP()
	checkWithExit(err)
	fmt.Println("my ip address: ", myIP)

	port := ":3000"
	protocol := "tcp"
	tcpAddr, err := net.ResolveTCPAddr(protocol, port)
	checkWithExit(err)
	listner, err := net.ListenTCP(protocol, tcpAddr)
	checkWithExit(err)
	fmt.Println("waiting for the connection ...")
	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		} else {
			fmt.Println("connected by .. ", conn.RemoteAddr().String())
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	messageBuf := make([]byte, 1024)
	messageLen, err := conn.Read(messageBuf)
	checkWithExit(err)
	message := string(messageBuf[:messageLen])
	fmt.Println("received message: ", message)

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.Write([]byte(message))
}

func checkWithExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}
