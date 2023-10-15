package model

import "fmt"

const (
	Total BandNum = iota
	Delta
	Theta
	Alpha
	Beta
	Ratio
	NumBands = 6
)

type BandNum uint

// AnalyzedEEG is an array that stores the EEG analysis values for all available channels for 1 second.
type AnalyzedEEG [NumAvailableChs][NumBands]ChannelPower

// ChannelPower is a struct for storing the data of the power of each channel.
type ChannelPower struct {
	// データが得られた時の測定経過時間
	Time uint
	// チャンネル番号
	ChNum uint
	// 帯域の番号
	BandNum BandNum
	// 直前5秒間に対する周波数解析の結果から算出された帯域のパワー
	Power float64
	// 直前の1秒間における生波形データの最大値
	MaxEEG float64
	// 直前の1秒間における生波形データの最小値
	MinEEG float64
}

func (a *AnalyzedEEG) ToCSVRow() []string {
	var row []string
	row = append(row, fmt.Sprint(a[0][Total].Time))
	for _, ch := range a {
		// Timeが0ならば、データが得られていないチャネル
		if ch[Total].Time == 0 {
			continue
		}
		row = append(row, []string{
			fmt.Sprintf("%.5f", ch[Total].Power),
			fmt.Sprintf("%.5f", ch[Delta].Power),
			fmt.Sprintf("%.5f", ch[Theta].Power),
			fmt.Sprintf("%.5f", ch[Alpha].Power),
			fmt.Sprintf("%.5f", ch[Beta].Power),
			fmt.Sprintf("%.5f", ch[Ratio].Power),
		}...)
	}
	return row
}

func (a *AnalyzedEEG) ToCSVHeader(at AnalysisType) []string {
	header := []string{""}
	for i, t := range at {
		if t == EEGCh {
			header = append(header, []string{
				fmt.Sprintf("total_ch%d", i+1),
				fmt.Sprintf("delta_ch%d", i+1),
				fmt.Sprintf("theta_ch%d", i+1),
				fmt.Sprintf("alpha_ch%d", i+1),
				fmt.Sprintf("beta_ch%d", i+1),
				fmt.Sprintf("ratio_ch%d", i+1),
			}...)
		}
	}
	return header
}
