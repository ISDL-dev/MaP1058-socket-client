package model

import (
	"fmt"
	"time"
)

const (
	// 1回の受信に含まれるバイト数
	NumTotalBytes = (NumChannels * NumPoints * NumADBytes) + SumCheckCodeSize
	// AD値と解釈するために必要なバイト数
	NumADBytes = 2
	// サムチェックコードを示すバイト列の長さ
	SumCheckCodeSize = 4
	// チャンネル数
	NumChannels = 16
	// 実際に利用可能なチャンネル数（16種類のうち1~8のみが利用可能）
	NumAvailableChs = 8
	// 1回の受信で1つのチャンネルが取得できる信号の数
	NumPoints = 50
	// 先頭の何ポイントまでの値をサムチェックに用いるか
	NumPntsSumCheck = 10
)

func NewSignals() *Signals {
	return &Signals{
		Time: time.Now(),
	}
}

// 送信バイナリーデータ1回分のAD値
type Signals struct {
	Time     time.Time
	Channels [NumAvailableChs]channelSignal
}

type channelSignal struct {
	ADValues     [NumPoints]uint16
	Measurements [NumPoints]float64
}

func (s *Signals) SetMeasurements(cal Cal) error {
	for i, ch := range s.Channels {
		for j, adV := range ch.ADValues {
			chCal := cal[i]
			if chCal.CalAD == 0 {
				return fmt.Errorf("value of CAL_AD must not be 0. CAL_AD at %dch is 0", i+1)
			}
			s.Channels[i].Measurements[j] = (float64(adV)-chCal.BaseAD)*(chCal.EuHi-chCal.EuLo)/chCal.CalAD + chCal.EuLo
		}
	}
	return nil
}
