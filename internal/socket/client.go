//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/socket-capture-signals/$GOPACKAGE
package socket

import (
	"fmt"
	"net"
)

type Conn interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

func Connect(conf SocketConfig) (*net.TCPConn, error) {
	serverAdd, err := net.ResolveTCPAddr("tcp", conf.ServerIP+":"+conf.ServerPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get server address: %w", err)
	}
	clientAdd := &net.TCPAddr{
		IP:   net.ParseIP(conf.ClientIP),
		Port: conf.ClientPort,
	}
	conn, err := net.DialTCP("tcp", clientAdd, serverAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to make connection: %w", err)
	}
	return conn, nil
}
