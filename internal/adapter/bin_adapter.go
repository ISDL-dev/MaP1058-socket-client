package adapter

import (
	"context"
	"fmt"
	"io"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// バイナリーデータで生波形データのAD値を受信する
type BinAdapter interface {
	ReceiveADValues(ctx context.Context) (*model.Signals, error)
}

func NewBinAdapter(c socket.Conn, p parser.Parser) BinAdapter {
	return &binAdapter{
		Conn:   c,
		Parser: p,
	}
}

type binAdapter struct {
	Conn   socket.Conn
	Parser parser.Parser
}

// AD値を受信する
func (a *binAdapter) ReceiveADValues(ctx context.Context) (*model.Signals, error) {
	rawBytes := make([]byte, model.NumTotalBytes)
	n, err := a.Conn.Read(rawBytes)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to receive binary data %w", err)
	}
	if n == 0 {
		return nil, nil
	}

	signals := model.NewSignals()
	err = a.Parser.ToSignals(rawBytes[:n], signals)
	if err != nil {
		if e, ok := err.(*parser.FailureSumCheckError); ok {
			if err := a.sendNAK(); err != nil {
				return nil, fmt.Errorf("%s, and failed to send NAK to the server", e.Error())
			}
		}
		return nil, fmt.Errorf("failed to capture valid signals: %w", err)
	}
	err = a.sendACK()
	if err != nil {
		return nil, err
	}
	return signals, nil
}

func (a *binAdapter) sendACK() error {
	_, err := a.Conn.Write([]byte{0x06})
	if err != nil {
		return fmt.Errorf("failed to write connection ACK: %w", err)
	}
	return nil
}

func (a *binAdapter) sendNAK() error {
	_, err := a.Conn.Write([]byte{0x15})
	if err != nil {
		return fmt.Errorf("failed to write connection NAK: %w", err)
	}
	return nil
}
