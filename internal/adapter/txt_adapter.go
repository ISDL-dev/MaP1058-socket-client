package adapter

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Be3751/socket-capture-signals/internal/model"
)

// テキストデータでトレンドデータの受信やコマンドの送受信をする
type TxtAdapter interface {
	StartRec(ctx context.Context, recTime time.Duration, recDateTime string) error
	EndRec(ctx context.Context) error
}

func NewTxtAdapter(c *net.TCPConn) TxtAdapter {
	return &binAdapter{
		Conn: c,
	}
}

type binAdapter struct {
	Conn *net.TCPConn
}

func (a *binAdapter) StartRec(ctx context.Context, recSecond time.Duration, recDateTime string) error {
	// recDateTimeのフォーマットを整える
	// var recDateTimeStr string

	sCmd := model.Command{
		Name:   "START",
		Params: []string{recSecond.String(), recDateTime},
	}
	sCmdStr := sCmd.NewString()
	_, err := a.Conn.Write([]byte(sCmdStr))
	if err != nil {
		return fmt.Errorf("failed to send %s: %w", sCmd, err)
	}

	var rCmd []byte
	_, err = a.Conn.Read(rCmd)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	if string(rCmd) != sCmdStr {
		return fmt.Errorf("failed to start recording")
	}
	return nil
}

func (a *binAdapter) EndRec(ctx context.Context) error {
	sCmd := "<SCMD>END:A:,,,,,,,,,</SCMD>"
	_, err := a.Conn.Write([]byte(sCmd))
	if err != nil {
		return fmt.Errorf("failed to send %s: %w", sCmd, err)
	}

	var rCmd []byte
	_, err = a.Conn.Read(rCmd)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	if string(rCmd) != sCmd {
		return fmt.Errorf("failed to end recording")
	}

	return nil
}
