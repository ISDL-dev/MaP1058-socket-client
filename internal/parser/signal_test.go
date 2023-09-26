package parser

import (
	"testing"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestToSignals(t *testing.T) {
	const (
		sumBytes         = 1604
		sumCheckCodeSize = 4
	)

	t.Run("信号をパースする", func(t *testing.T) {
		parser := NewParser()
		rawBytes := make([]byte, 1604)
		for i := 0; i < sumBytes-sumCheckCodeSize; i += 32 {
			for j := 0; j < 16; j += 2 {
				rawBytes[i+j] = 0x01
				rawBytes[i+j+1] = 0x00
			}
		}
		rawBytes[sumBytes-4] = 0x50
		rawBytes[sumBytes-3] = 0x00
		rawBytes[sumBytes-2] = 0x00
		rawBytes[sumBytes-1] = 0x00
		s, err := parser.ToSignals(rawBytes)
		assert.NoError(t, err)
		for c := 0; c < model.NumAvailableChs; c++ {
			for p := 0; p < model.NumPoints; p++ {
				assert.Equal(t, uint16(1), s.Channels[c].ADValues[p])
			}
		}
	})

	t.Run("規定の長さでないバイト列を受け取ってエラー", func(t *testing.T) {
		parser := NewParser()
		rawBytes := []byte{0x00, 0x01, 0x02}
		_, err := parser.ToSignals(rawBytes)
		assert.Error(t, err)
	})

	t.Run("サムチェックの結果が合わずエラー", func(t *testing.T) {
		parser := NewParser()
		rawBytes := make([]byte, 1604)
		for i := 0; i < sumBytes-sumCheckCodeSize; i += 32 {
			for j := 0; j < 32; j += 2 {
				rawBytes[i+j] = 0x01
				rawBytes[i+j+1] = 0x00
			}
		}
		rawBytes[sumBytes-4] = 0x50
		rawBytes[sumBytes-3] = 0x00
		rawBytes[sumBytes-2] = 0x00
		rawBytes[sumBytes-1] = 0x00
		_, err := parser.ToSignals(rawBytes)
		assert.EqualValues(t, &FailureSumCheckError{Expected: 80, Actual: 160}, err)
	})
}
