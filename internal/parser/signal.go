//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
)

func (p *parser) ToSignals(b []byte, s *model.Signals) error {
	if len(b) != int(model.NumTotalBytes) {
		return fmt.Errorf("the arg b's len must be %d", model.NumTotalBytes)
	}

	signalBuf := bytes.NewBuffer(b[:model.NumTotalBytes-model.SumCheckCodeSize])
	var sum uint16
	for pnt := 0; pnt < int(model.NumPoints); pnt++ {
		for ch := 0; ch < model.NumAvailableChs; ch++ {
			var adValue uint16
			err := binary.Read(signalBuf, binary.BigEndian, &adValue)
			if err != nil {
				return fmt.Errorf("failed to read binary: %w", err)
			}
			s.Channels[ch].ADValues[pnt] = adValue
			if pnt < model.NumPntsSumCheck {
				sum += adValue
			}
		}
		signalBuf.Next(int(model.NumChannels-model.NumAvailableChs) * model.NumADBytes) // 後半8個のチャンネルは未使用
	}

	var valueSumCheckCode uint32
	sumCheckBuf := bytes.NewBuffer(b[model.NumTotalBytes-model.SumCheckCodeSize:])
	err := binary.Read(sumCheckBuf, binary.BigEndian, &valueSumCheckCode)
	if err != nil {
		return fmt.Errorf("failed to read binary: %w", err)
	}

	if int(valueSumCheckCode) != int(sum) {
		return &FailureSumCheckError{Expected: uint16(valueSumCheckCode), Actual: sum}
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
