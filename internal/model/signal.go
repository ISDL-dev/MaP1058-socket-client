package model

const (
	// 1回の受信に含まれるバイト数
	SumBytes = (numCh * numPnt * adLen) + sumCheckCodeSize
	// チャンネル数
	numCh = 16
	// 実際に利用可能なチャンネル数（16種類のうち1~8のみが利用可能）
	numAvailableCh = 8
	// 1回の受信で1つのチャンネルが取得できる信号の数
	numPnt = 50
	// AD値を意味するバイト列の長さ
	adLen = 2
	// サムチェックコードを示すバイト列の長さ
	sumCheckCodeSize = 4
)

// 送信バイナリーデータ1回分のAD値
type Signals struct {
	Channels [numAvailableCh]channelSignal
}

type channelSignal struct {
	ADValues     [numPnt]uint16
	Measurements [numPnt]float64
}

func (s *Signals) SetMeasurements(cal *Cal) error {
	for i, ch := range s.Channels {
		for j, adV := range ch.ADValues {
			chCal := cal.Channels[i]
			m := float64(adV-chCal.BaseAD)*(chCal.EuHi-chCal.EuLo)/float64(chCal.CalAD) + chCal.EuLo
			ch.Measurements[j] = m
		}
	}
	return nil
}

// AD値から測定値に変換するための校正値
type Cal struct {
	Channels [numAvailableCh]channelCal
}

type channelCal struct {
	BaseAD uint16
	CalAD  uint16
	EuHi   float64
	EuLo   float64
}
