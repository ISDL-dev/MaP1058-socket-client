package socket

import (
	"context"
	"fmt"
	"net"
	"time"
)

type SocketClient interface {
	StartRec(ctx context.Context, recSec time.Duration, recDateTime string) error
	EndRec(ctx context.Context) error
	CaptureSig(ctx context.Context, buf []byte) error
	// GetStat(ctx context.Context) ([]byte, error)
	// Mark(ctx context.Context) error
	// UnMark(ctx context.Context) error
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

func (c *client) StartRec(ctx context.Context, recSec time.Duration, recDateTime string) error {
	sCmd := fmt.Sprintf("<SCMD>START:A:%d,\"%s\",,,,,,,,,</SCMD>", recSec, recDateTime)
	err := c.sendCmd(ctx, sCmd)
	if err != nil {
		return fmt.Errorf("failed to send command %w", err)
	}

	rCmd, err := c.receiveCmd(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	if rCmd != sCmd {
		return fmt.Errorf("failed to start recording")
	}

	return nil
}

func (c *client) EndRec(ctx context.Context) error {
	sCmd := "<SCMD>END:A:,,,,,,,,,</SCMD>"
	err := c.sendCmd(ctx, sCmd)
	if err != nil {
		return fmt.Errorf("failed to start recording %w", err)
	}

	rCmd, err := c.receiveCmd(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	if rCmd != sCmd {
		return fmt.Errorf("failed to end recording")
	}

	return nil
}

func (c *client) CaptureSig(ctx context.Context, buf []byte) error {
	sCmd := "<SCMD>DATA_EEG:A:,,,,,,,,,</SCMD>"
	err := c.sendCmd(ctx, sCmd)
	if err != nil {
		return fmt.Errorf("failed to capture signals %w", err)
	}
	return nil
}

func (c *client) sendCmd(ctx context.Context, cmd string) error {
	_, err := c.conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send \"%s\": %w", cmd, err)
	}
	return nil
}

func (c *client) receiveCmd(ctx context.Context) (string, error) {
	buf := make([]byte, 256)
	_, err := c.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to receive command %w", err)
	}
	return string(buf), nil
}

// TODO: 受信コマンドをParseする関数も必要
