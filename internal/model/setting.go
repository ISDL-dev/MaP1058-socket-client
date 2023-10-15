package model

type Setting struct {
	TrendRange   TrendRange
	AnalysisType AnalysisType
	Calibration  Calibration
}

// 全てのチャネルにおける表示範囲
type TrendRange [NumAvailableChs]ChannelRange

// トレンドグラフの表示範囲
type ChannelRange struct {
	Upper uint
	Lower uint
}

// 全てのチャネルにおける解析種別
type AnalysisType [NumAvailableChs]ChannelType

// 解析種別
type ChannelType uint

const (
	NoAnalysis ChannelType = iota
	EMGCh
	EOGHCh
	EOGVCh
	EEGCh
	RRIntCh
	CICh
	RRCVCh
	HELFCh
	WAVECh
	ExtPlsCh
	RespCh
	Resp2Ch
)

// 全てのチャネルにおける校正値
type Calibration [NumAvailableChs]ChannelCal

// AD値から測定値に変換するための校正値
type ChannelCal struct {
	BaseAD uint16
	CalAD  uint16
	EuHi   float64
	EuLo   float64
}
