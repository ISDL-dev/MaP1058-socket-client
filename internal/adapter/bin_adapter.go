package adapter

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
	"io"
	"time"
)

type BinAdapter interface {
	// WriteRawSignal receives raw signals and writes them to a csv file.
	WriteRawSignal(ctx context.Context, rcvSuccess <-chan bool, stg *model.Setting) error

	// Pause keeps receiving signals but does not write them to a csv file.
	Pause()

	// Resume restarts writing signals to a csv file. It is called after Pause.
	Resume()
}

func NewBinAdapter(c socket.Conn, p parser.Parser, w io.Writer) BinAdapter {
	return &binAdapter{
		Conn:   c,
		Parser: p,
		File:   w,
	}
}

type binAdapter struct {
	Conn   socket.Conn
	Parser parser.Parser
	File   io.Writer
	// noWrite is a flag to prevent writing signals to a csv file. If true, it receives signals but does not write them.
	noWrite bool
}

const (
	bufferSize = 10
)

func (a *binAdapter) WriteRawSignal(ctx context.Context, rcvSuccess <-chan bool, stg *model.Setting) (err error) {
	defer func() {
		err := a.Conn.Close()
		if err != nil {
			err = fmt.Errorf("failed to close connection: %w", err)
			panic(err)
		}
	}()

	csvWriter := csv.NewWriter(a.File)
	defer func() {
		csvWriter.Flush()
		if hereErr := csvWriter.Error(); hereErr != nil {
			err = fmt.Errorf("%s: failed to flush csv writer: %w", err.Error(), hereErr)
		}
	}()
	if err = csvWriter.Write(model.SignalHeader()); err != nil {
		return fmt.Errorf("failed to write header to csv: %w", err)
	}

	var buf [][]string
	var timeReceived int
LOOP:
	for {
		select {
		case <-rcvSuccess: // the receiving process is complete.
			break LOOP
		case <-ctx.Done():
			break LOOP
		default:
			signals, err := a.receiveAD()
			if err != nil {
				var targetErr *parser.FailureSumCheckError
				if errors.As(err, &targetErr) {
					if err := a.sendNAK(); err != nil {
						return fmt.Errorf("%s, and failed to send NAK to the server", err.Error())
					}
					continue
				}
				return fmt.Errorf("failed to receive AD values: %w", err)
			}
			if a.noWrite {
				continue
			}

			if err = signals.SetMeasurements(stg.Calibration, stg.AnalysisType); err != nil {
				return fmt.Errorf("failed to set measurements: %w", err)
			}

			buf = append(buf, signals.ToRecords(timeReceived)...)
			if err = a.sendACK(); err != nil {
				return err
			}
			timeReceived++
			// write records to csv when buffer is full
			if timeReceived%bufferSize == 0 {
				for _, record := range buf {
					if err = csvWriter.Write(record); err != nil {
						return fmt.Errorf("failed to write raw signal records to csv: %w", err)
					}
				}
				buf = [][]string{}
			}
		}
		// prevent busy loop
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}

func (a *binAdapter) Pause() {
	a.noWrite = true
}

func (a *binAdapter) Resume() {
	a.noWrite = false
}

func (a *binAdapter) receiveAD() (*model.Signals, error) {
	var rawBytes []byte
	rawBytes, err := a.read(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary data: %w", err)
	}

	s, err := a.Parser.ToSignals(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse binary data to Signals: %w", err)
	}
	return s, nil
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

// 再帰的にmodel.NumTotalBytesのバイト数になるまで受信する
func (a *binAdapter) read(rawBytes []byte) ([]byte, error) {
	b := make([]byte, model.NumTotalBytes)
	n, err := a.Conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to receive binary data: %w", err)
	}
	if n == 0 {
		return nil, errors.New("tried receiving but got 0 byte")
	}
	rawBytes = append(rawBytes, b[:n]...)
	if len(rawBytes) == model.NumTotalBytes {
		return rawBytes, nil
	} else {
		return a.read(rawBytes)
	}
}
