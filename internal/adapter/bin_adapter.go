package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

// バイナリーデータで生波形データのAD値を受信する
type BinAdapter interface {
	// AD値を受信する
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

func (a *binAdapter) ReceiveADValues(ctx context.Context) (*model.Signals, error) {
	rawBytes := make([]byte, model.NumTotalBytes)
	for {
		n, err := a.Conn.Read(rawBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to receive binary data: %w", err)
		}
		if n == 0 {
			return nil, errors.New("received 0 byte")
		}

		s, err := a.Parser.ToSignals(rawBytes)
		if err == nil {
			if err := a.sendACK(); err != nil {
				return nil, err
			}
			return s, nil
		} else if _, ok := err.(*parser.FailureSumCheckError); ok {
			if err := a.sendNAK(); err != nil {
				return nil, fmt.Errorf("%s, and failed to send NAK to the server", err.Error())
			}
		} else {
			return nil, fmt.Errorf("failed to parse binary data to Signals: %w", err)
		}
		rawBytes = make([]byte, model.NumTotalBytes)
	}
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
