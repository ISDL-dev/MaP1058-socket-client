package adapter

import (
	"context"
	"fmt"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// バイナリーデータで生波形データのAD値を受信する
type BinAdapter interface {
	ReceiveADValues(ctx context.Context) (*model.Signals, error)
}

func NewBinAdapter(c socket.Conn, p parser.Parser) BinAdapter {
	return &adapter{
		Conn:   c,
		Parser: p,
	}
}

type adapter struct {
	Conn   socket.Conn
	Parser parser.Parser
}

// AD値を受信する
func (a *adapter) ReceiveADValues(ctx context.Context) (*model.Signals, error) {
	rawBytes := make([]byte, model.NumTotalBytes)
	_, err := a.Conn.Read(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to receive binary data %w", err)
	}
	signals := model.NewSignals()
	err = a.Parser.ToSignals(rawBytes, &signals)
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
	return &signals, nil
}

func (a *adapter) sendACK() error {
	_, err := a.Conn.Write([]byte("ACK"))
	if err != nil {
		return fmt.Errorf("failed to write connection ACK: %w", err)
	}
	return nil
}

func (a *adapter) sendNAK() error {
	_, err := a.Conn.Write([]byte("NAK"))
	if err != nil {
		return fmt.Errorf("failed to write connection NAK: %w", err)
	}
	return nil
}
