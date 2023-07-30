package adapter

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// テキストデータでトレンドデータの受信やコマンドの送受信をする
type TxtAdapter interface {
	StartRec(ctx context.Context, recTime time.Duration, recDateTime time.Time) error
	EndRec(ctx context.Context) error
}

func NewTxtAdapter(c socket.Conn) TxtAdapter {
	return &binAdapter{
		Conn: c,
	}
}

type binAdapter struct {
	Conn socket.Conn
}

func (a *binAdapter) StartRec(ctx context.Context, recSecond time.Duration, recDateTime time.Time) error {
	recDateTimeParam := recDateTime.Format("2006/01/02 15-04-05")
	recSecondParam := strSecond(recSecond)
	sCmd := model.Command{
		Name:   "START",
		Params: [model.NumSeparator + 1]string{recSecondParam, recDateTimeParam},
	}
	sCmdStr := sCmd.String()
	_, err := a.Conn.Write([]byte(sCmdStr))
	if err != nil {
		return fmt.Errorf("failed to send %s: %w", sCmd, err)
	}
	buf := make([]byte, 128)
	readLen, err := a.Conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	rCmdStr := string(buf[:readLen])
	if rCmdStr != sCmdStr {
		return fmt.Errorf("failed to start recording because %s doesn't match with %s", rCmdStr, sCmdStr)
	}
	return nil
}

func (a *binAdapter) EndRec(ctx context.Context) error {
	sCmd := model.Command{
		Name: "END",
	}
	sCmdStr := sCmd.String()
	_, err := a.Conn.Write([]byte(sCmdStr))
	if err != nil {
		return fmt.Errorf("failed to send %s: %w", sCmd, err)
	}
	buf := make([]byte, 128)
	readLen, err := a.Conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	rCmdStr := string(buf[:readLen])
	if string(rCmdStr) != sCmdStr {
		return fmt.Errorf("failed to end recording")
	}
	return nil
}

func strSecond(d time.Duration) string {
	var secondStr string
	if d >= time.Minute {
		second := int64(d) / int64(time.Second)
		secondStr = strconv.Itoa(int(second))
	} else {
		secondStr = strings.Replace(d.String(), "s", "", -1)
	}
	return secondStr
}
