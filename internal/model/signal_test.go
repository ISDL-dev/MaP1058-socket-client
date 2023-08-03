package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetMeasurements(t *testing.T) {
	t.Run("計測値をセット", func(t *testing.T) {
		signals := NewSignals()
		for ch := range signals.Channels {
			for pnt := range signals.Channels[ch].ADValues {
				signals.Channels[ch].ADValues[pnt] = 1
			}
		}
		cal := Cal{
			{1, 1, 0, 2},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
		}
		err := signals.SetMeasurements(cal)
		assert.NoError(t, err)
		for ch := range signals.Channels {
			for pnt := range signals.Channels[ch].Measurements {
				if ch == 0 {
					assert.Equal(t, float64(2), signals.Channels[ch].Measurements[pnt])
				} else {
					assert.Equal(t, float64(0), signals.Channels[ch].Measurements[pnt])
				}
			}
		}
	})

	t.Run("CAL_ADの値が0のCalを受け取ってエラー", func(t *testing.T) {
		signals := NewSignals()
		cal := Cal{
			{1, 0, 0, 2},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
		}
		err := signals.SetMeasurements(cal)
		assert.Error(t, err)
	})
}

func randomRawBytes(signalByteSize int, sumCheckCodeSize int) []byte {
	rawBytes := make([]byte, signalByteSize)
	for i := 0; i < signalByteSize; i++ {
		// TODO: byte型の乱数を生成する
		rawBytes = append(rawBytes)
	}
	return rawBytes
}
