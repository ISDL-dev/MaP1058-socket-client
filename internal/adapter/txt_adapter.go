package adapter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/scanner"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// TxtAdapter テキストデータでトレンドデータの受信やコマンドの送受信をする
type TxtAdapter interface {
	StartRec(ctx context.Context, recTime time.Duration, recDateTime time.Time) error
	EndRec(ctx context.Context) error
	GetStatus(ctx context.Context) (model.Status, error)
	GetTrendData(ctx context.Context, w CSVWriterGroup, at model.AnalysisType) error
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
	var (
		rangeCnt    int
		analysisCnt int
		calCnt      int
	)
	for a.Scanner.Scan() {
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
			s.AnalysisType = as
			analysisCnt++
		case "GETSETTING":
			// 値を含む受信コマンドは前半8チャネルのみ
			if calCnt < model.NumChannels-model.NumAvailableChs {
				chc, err := a.Parser.ToChannelCal(cmd)
				if err != nil {
					return nil, fmt.Errorf("failed to convert %s to Calibration: %w", cmdStr, err)
				}
				s.Calibration[calCnt] = chc
			}
			calCnt++
		default:
			return nil, fmt.Errorf("invalid command: %s", cmdStr)
		}
		if rangeCnt == 1 && analysisCnt == 1 && calCnt == model.NumChannels {
			break
		}
	}
	if err := a.Scanner.Err(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	return &s, nil
}

type CSVWriterGroup struct {
	EEGWriter   csv.Writer
	HRWriter    io.Writer
	AnsWriter   io.Writer
	ExplsWriter io.Writer
	RespWriter  io.Writer
}

func (a *txtAdapter) GetTrendData(ctx context.Context, w CSVWriterGroup, at model.AnalysisType) error {
	var analyzedEEG model.AnalyzedEEG

	err := w.EEGWriter.Write(analyzedEEG.ToCSVHeader(at))
	if err != nil {
		return fmt.Errorf("failed to write AnalyzedEEG header to csv: %w", err)
	}

	for a.Scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
			cmdStr := a.Scanner.Text()
			cmd, err := a.Parser.ToCommand(cmdStr)
			if err != nil {
				return fmt.Errorf("failed to convert %s to Command: %w", cmdStr, err)
			}
			switch cmd.Name {
			case "DATA":
				continue
			case "DATA_HR":
				continue
			case "DATA_ANS":
				continue
			case "DATA_EXPLS":
				continue
			case "DATA_RESP":
				continue
			case "DATA_RESP2":
				continue
			case "DATA_RESP2UP":
				continue
			case "DATA_RESP2DP":
				continue
			case "DATA_EEG":
				power, err := a.Parser.ToChannelPower(cmd)
				if err != nil {
					return fmt.Errorf("failed to convert %s to AnalyzedEEG: %w", cmdStr, err)
				}
				analyzedEEG[power.ChNum][power.BandNum] = power
			case "EVENT_SEC":
				err := w.EEGWriter.Write(analyzedEEG.ToCSVRow())
				if err != nil {
					return fmt.Errorf("failed to write AnalyzedEEG to csv: %w", err)
				}
			case "STATUS":
				continue
			case "EVENT_MARK":
				continue
			case "EVENT_MARKCANCEL":
				continue
			case "GUIDANCE":
				continue
			default:
				return fmt.Errorf("invalid command: %s", cmdStr)
			}
		}
	}
	if err := a.Scanner.Err(); err != nil {
		return fmt.Errorf("invalid input: %w", err)
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
