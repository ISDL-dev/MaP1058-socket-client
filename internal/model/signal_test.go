package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetMeasurements(t *testing.T) {
	t.Run("計測値をセット", func(t *testing.T) {
		signals := NewSignals()
		for ch := range signals.Channels {
			for pnt := range signals.Channels[ch].ADValues {
				signals.Channels[ch].ADValues[pnt] = 2
			}
		}
		cal := Calibration{
			{0, 1, 2, 2},
			{1, 1, 2, 1},
			{1, 1, 2, 1},
			{1, 1, 2, 1},
			{1, 1, 2, 1},
			{1, 1, 2, 1},
			{1, 1, 2, 1},
			{1, 1, 2, 1},
		}
		at := AnalysisType{
			EEGCh,
			EEGCh,
			EEGCh,
			EEGCh,
			RRIntCh,
			NoAnalysis,
			EEGCh,
			EEGCh,
		}
		err := signals.SetMeasurements(cal, at)
		assert.NoError(t, err)
		for ch := range signals.Channels {
			for pnt := range signals.Channels[ch].Measurements {
				if ch == 0 {
					assert.Equal(t, float64(2), signals.Channels[ch].Measurements[pnt])
				} else if ch == 5 {
					assert.Equal(t, float64(0), signals.Channels[ch].Measurements[pnt])
				} else {
					assert.Equal(t, float64(2), signals.Channels[ch].Measurements[pnt])
				}
			}
		}
	})

	t.Run("CAL_ADの値が0のCalを受け取ってエラー", func(t *testing.T) {
		signals := NewSignals()
		cal := Calibration{
			{1, 0, 0, 2},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
		}
		at := AnalysisType{
			EEGCh,
			EEGCh,
			EEGCh,
			EEGCh,
			RRIntCh,
			NoAnalysis,
			EEGCh,
			EEGCh,
		}
		err := signals.SetMeasurements(cal, at)
		assert.Error(t, err)
	})
}

func TestToRecords(t *testing.T) {
	t.Run("計測値をcsvのレコードに変換", func(t *testing.T) {
		signals := NewSignals()
		for ch := range signals.Channels {
			for pnt := range signals.Channels[ch].Measurements {
				signals.Channels[ch].Measurements[pnt] = float64(1)
			}
		}
		records := signals.ToRecords(0)
		assert.Equal(t, 50, len(records))
		for i, record := range records {
			assert.Equal(t, 9, len(record))
			for j, cell := range record {
				if j == 0 {
					assert.Equal(t, fmt.Sprint(i+1), record[j])
				} else {
					assert.Equal(t, fmt.Sprintf("%.5f", 1.0), cell)
				}
			}
		}
		records = signals.ToRecords(5)
		assert.Equal(t, 50, len(records))
		for i, record := range records {
			assert.Equal(t, 9, len(record))
			for j, cell := range record {
				if j == 0 {
					assert.Equal(t, fmt.Sprint(i+50*5+1), record[j])
				} else {
					assert.Equal(t, fmt.Sprintf("%.5f", 1.0), cell)
				}
			}
		}
	})
}
