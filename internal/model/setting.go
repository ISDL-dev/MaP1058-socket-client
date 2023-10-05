package model

type Setting struct {
	TrendRange  TrendRange
	Analysis    Analysis
	Calibration Calibration
}

// 全てのチャネルにおける表示範囲
type TrendRange [NumAvailableChs]ChannelRange

// トレンドグラフの表示範囲
type ChannelRange struct {
	Upper uint
	Lower uint
}

// 　全てのチャネルにおける解析種別
type Analysis []ChannelAnalysis

// 解析種別
type ChannelAnalysis uint

const (
	NoAnalysis ChannelAnalysis = iota
	EMG
	EOG_H
	EOG_V
	EEG
	RRInt
	CI
	RRCV
	HELF
	WAVE
	ExtPls
	Resp
	Resp2
)

// 全てのチャネルにおける校正値
type Calibration []ChannelCal

// AD値から測定値に変換するための校正値
type ChannelCal struct {
	BaseAD uint16
	CalAD  uint16
	EuHi   float64
	EuLo   float64
}
