package socket

import (
	"context"
	"fmt"
	"net"
)

type SocketClient interface {
	SendCmd(ctx context.Context, cmd string) error
	ReceiveCmd(ctx context.Context) (string, error)
	// ReceiveBin(ctx context.Context) ([]byte, error)
}

var _ SocketClient = (*client)(nil)

func NewSocketClient(c *net.TCPConn) SocketClient {
	return &client{
		conn: c,
	}
}

type client struct {
	conn *net.TCPConn
}

func NewClient(conf config.Config) (*net.TCPConn, error) {
	serverAdd, err := net.ResolveTCPAddr("tcp", conf.ServerIP+":"+conf.ServerPortText)
	if err != nil {
		return nil, err
	}
	clientAdd := &net.TCPAddr{
		IP:   net.IP(conf.ClientIP),
		Port: conf.ClientPort,
	}

	// TCPネットワークの接続
	conn, err := net.DialTCP("tcp", clientAdd, serverAdd)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *client) SendCmd(ctx context.Context, cmd string) error {
	_, err := c.conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send \"%s\": %w", cmd, err)
	}
	return nil
}

func (c *client) ReceiveCmd(ctx context.Context) (string, error) {
	buf := make([]byte, 256)
	_, err := c.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to receive command %w", err)
	}
	return string(buf), nil
}
