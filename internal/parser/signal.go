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
	if len(adSignals) != int(p.Config.SumBytes) {
		return nil, fmt.Errorf("adSignals' len must be %d", p.Config.SumBytes)
	}
	result := &model.Signals{}
	var sum uint16
	for i := 0; i < int(p.Config.SumBytes)-int(p.Config.SumCheckCodeSize); i += 32 {
		for j := 0; j < 16; j += 2 {
			var adValue uint16
			buf := bytes.NewBuffer(adSignals[i+j : i+j+2])
			err := binary.Read(buf, binary.BigEndian, &adValue)
			if err != nil {
				return nil, fmt.Errorf("failed to read binary: %w", err)
			}

			ch := int(j / 2)
			pnt := int(i / 32)
			result.Channels[ch].ADValues[pnt] = adValue

			if pnt < 10 {
				sum += adValue
			}
		}
	}
	var valueSumCheckCode uint16
	buf := bytes.NewBuffer(adSignals[p.Config.SumBytes-p.Config.SumCheckCodeSize+2:])
	err := binary.Read(buf, binary.BigEndian, &valueSumCheckCode)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary: %w", err)
	}

	if valueSumCheckCode != sum {
		return nil, &FailureSumCheckError{Expected: valueSumCheckCode, Actual: sum}
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
