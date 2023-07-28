//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
)

// AD値のバイト列を解析してAD値を持つmodel.Signals型のポインタを返す
func (p *parser) ToSignals(adSignals []byte) (*model.Signals, error) {
	if len(adSignals) != int(model.NumTotalBytes) {
		return nil, fmt.Errorf("adSignals' len must be %d", model.NumTotalBytes)
	}

	signalBuf := bytes.NewBuffer(adSignals[:model.NumTotalBytes-model.SumCheckCodeSize])
	result := &model.Signals{}
	var sum uint16
	for pnt := 0; pnt < int(model.NumPoints); pnt++ {
		for ch := 0; ch < model.NumAvailableChs; ch++ {
			var adValue uint16
			err := binary.Read(signalBuf, binary.BigEndian, &adValue)
			if err != nil {
				return nil, fmt.Errorf("failed to read binary: %w", err)
			}
			result.Channels[ch].ADValues[pnt] = adValue
			if pnt < model.NumPntsSumCheck {
				sum += adValue
			}
		}
		signalBuf.Next(int(model.NumChannels-model.NumAvailableChs) * model.NumADBytes) // 後半8個のチャンネルは未使用
	}

	var valueSumCheckCode uint32
	sumCheckBuf := bytes.NewBuffer(adSignals[model.NumTotalBytes-model.SumCheckCodeSize:])
	err := binary.Read(sumCheckBuf, binary.BigEndian, &valueSumCheckCode)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary: %w", err)
	}

	if int(valueSumCheckCode) != int(sum) {
		return nil, &FailureSumCheckError{Expected: uint16(valueSumCheckCode), Actual: sum}
	}
	return result, nil
}

type FailureSumCheckError struct {
	Expected uint16
	Actual   uint16
}

func (e *FailureSumCheckError) Error() string {
	return fmt.Sprintf("parsed invalid signals doesn't match with the sum-check-code: expected %d but actual %d", e.Expected, e.Actual)
}
