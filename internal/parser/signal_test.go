package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSignals(t *testing.T) {
	const (
		sumBytes         = 1604
		sumCheckCodeSize = 4
	)

	t.Run("信号をパースする", func(t *testing.T) {
		pConf := ParseConfig{
			SumBytes:         sumBytes,
			SumCheckCodeSize: sumCheckCodeSize,
		}
		parser := NewParser(pConf)

		rawBytes := make([]byte, 1604)
		for i := 0; i < sumBytes-sumCheckCodeSize; i += 32 {
			for j := 0; j < 16; j++ {
				if (i+j)%2 == 0 {
					rawBytes[i+j] = 0x00
				} else {
					rawBytes[i+j] = 0x01
				}
			}
		}
		rawBytes[sumBytes-2] = 0x00
		rawBytes[sumBytes-1] = 0x50

		_, err := parser.ToSignals(rawBytes)
		assert.NoError(t, err)
	})

	t.Run("規定の長さでないバイト列を受け取ってエラー", func(t *testing.T) {
		pConf := ParseConfig{
			SumBytes:         sumBytes,
			SumCheckCodeSize: sumCheckCodeSize,
		}
		parser := NewParser(pConf)
		rawBytes := []byte{0x00, 0x01, 0x02}
		_, err := parser.ToSignals(rawBytes)
		assert.Error(t, err)
	})

	t.Run("サムチェックの結果が合わずエラー", func(t *testing.T) {
		pConf := ParseConfig{
			SumBytes:         sumBytes,
			SumCheckCodeSize: sumCheckCodeSize,
		}
		parser := NewParser(pConf)

		rawBytes := make([]byte, 1604)
		for i := 0; i < sumBytes-sumCheckCodeSize; i += 32 {
			for j := 0; j < 16; j++ {
				rawBytes[i+j] = 0x01
			}
		}
		rawBytes[sumBytes-2] = 0x00
		rawBytes[sumBytes-1] = 0x50

		_, err := parser.ToSignals(rawBytes)
		assert.EqualValues(t, &FailureSumCheckError{Expected: 80, Actual: 257 * 80}, err)
	})
}
