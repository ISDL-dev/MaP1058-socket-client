//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
)

func (p *parser) ToSignals(b []byte) (*model.Signals, error) {
	s := model.NewSignals()
	if len(b) != int(model.NumTotalBytes) {
		return nil, NewInvalidLenError(int(model.NumTotalBytes), len(b))
	}
	err := sumCheck(b)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b[:model.NumTotalBytes-4])
	for pnt := 0; pnt < model.NumPoints; pnt++ {
		for ch := 0; ch < model.NumAvailableChs; ch++ {
			var ad uint16
			err := binary.Read(
				buf,
				binary.LittleEndian,
				&ad,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to read binary: %w", err)
			}
			s.Channels[ch].ADValues[pnt] = ad
		}
		_ = buf.Next(int(model.NumChannels-model.NumAvailableChs) * model.NumADBytes) // don't use the 8 channels that are not available
	}
	return s, nil
}

func sumCheck(b []byte) error {
	var actual uint16
	buf := bytes.NewBuffer(b[:1600])
	for i := 0; i < 160; i++ {
		var ad uint16
		err := binary.Read(buf, binary.LittleEndian, &ad)
		if err != nil {
			return fmt.Errorf("failed to read AD value in Little Endian order: %w", err)
		}
		actual += ad
	}
	var expected uint16
	err := binary.Read(bytes.NewBuffer(b[1600:1602]), binary.LittleEndian, &expected)
	if err != nil {
		return fmt.Errorf("failed to read AD value in Little Endian order: %w", err)
	}
	if actual != expected {
		return &FailureSumCheckError{Expected: expected, Actual: actual}
	}
	return nil
}

type FailureSumCheckError struct {
	Expected uint16
	Actual   uint16
}

func (e *FailureSumCheckError) Error() string {
	return fmt.Sprintf("parsed invalid signals doesn't match with the sum-check-code: expected %d but actual %d", e.Expected, e.Actual)
}

func NewFailureSumCheckError(expected uint16, actual uint16) *FailureSumCheckError {
	return &FailureSumCheckError{Expected: expected, Actual: actual}
}

type InvalidLenError struct {
	Expected int
	Actual   int
}

func (e *InvalidLenError) Error() string {
	return fmt.Sprintf("the arg b's len must be %d, the actual len was %d", e.Expected, e.Actual)
}

func NewInvalidLenError(expected int, actual int) *InvalidLenError {
	return &InvalidLenError{Expected: expected, Actual: actual}
}