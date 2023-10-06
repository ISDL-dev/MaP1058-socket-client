package adapter

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/scanner"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// テキストデータでトレンドデータの受信やコマンドの送受信をする
type TxtAdapter interface {
	StartRec(ctx context.Context, recTime time.Duration, recDateTime time.Time) error
	EndRec(ctx context.Context) error
	GetStatus(ctx context.Context) (model.Status, error)
	GetSetting() (*model.Setting, error)
}

func NewTxtAdapter(c socket.Conn, s scanner.CustomScanner, p parser.Parser) TxtAdapter {
	return &txtAdapter{
		Conn:    c,
		Scanner: s,
		Parser:  p,
	}
}

type txtAdapter struct {
	Conn    socket.Conn
	Scanner scanner.CustomScanner
	Parser  parser.Parser
}

func (a *txtAdapter) StartRec(ctx context.Context, recSecond time.Duration, recDateTime time.Time) error {
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

func (a *txtAdapter) EndRec(ctx context.Context) error {
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
		return fmt.Errorf("the received cmd `%s` dosen't match with the sent one `%s`", rCmdStr, sCmdStr)
	}
	return nil
}

func (a *txtAdapter) GetStatus(ctx context.Context) (model.Status, error) {
	sCmd := model.Command{
		Name: "STATUS",
	}
	sCmdStr := sCmd.String()
	_, err := a.Conn.Write([]byte(sCmdStr))
	if err != nil {
		return "", fmt.Errorf("failed to send %s: %w", sCmd, err)
	}

	buf := make([]byte, 128)
	readLen, err := a.Conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to send %s: %w", sCmd, err)
	}
	rCmdStr := string(buf[:readLen])
	rCmd, err := a.Parser.ToCommand(rCmdStr)
	if err != nil {
		return "", fmt.Errorf("failed to convert %s to Command: %w", rCmdStr, err)
	}
	status := rCmd.Params[0]
	return model.Status(status), nil
}

func (a *txtAdapter) GetSetting() (*model.Setting, error) {
	var s model.Setting
	for a.Scanner.Scan() {
		var (
			rangeCnt    int
			analysisCnt int
			calCnt      int
		)
		cmdStr := a.Scanner.Text()
		cmd, err := a.Parser.ToCommand(cmdStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to Command: %w", cmdStr, err)
		}
		switch cmd.Name {
		case "RANGE":
			tr, err := a.Parser.ToTrendRange(cmd)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to TrendRange: %w", cmdStr, err)
			}
			s.TrendRange = tr
			rangeCnt++
		case "ANALYSIS":
			as, err := a.Parser.ToAnalysis(cmd)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to Analysis: %w", cmdStr, err)
			}
			s.Analysis = as
			analysisCnt++
		case "GETSETTING":
			// 後半8チャネルは収録されていない空のデータが送られてくるため無視する
			if calCnt >= model.NumChannels-model.NumAvailableChs {
				calCnt++
				continue
			}
			chc, err := a.Parser.ToChannelCal(cmd)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to Calibration: %w", cmdStr, err)
			}
			s.Calibration[calCnt] = chc
			calCnt++
		default:
			return nil, fmt.Errorf("invalid command: %s", cmdStr)
		}
		if rangeCnt == 1 && analysisCnt == 1 && calCnt == model.NumChannels {
			break
		}
	}
	if err := a.Scanner.Err(); err != nil {
		return nil, fmt.Errorf("Invalid input: %w", err)
	}
	return &s, nil
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
