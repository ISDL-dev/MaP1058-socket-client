package model

// 測定における全てのチャネルにおける校正値
type Cal [NumAvailableChs]ChannelCal

// AD値から測定値に変換するための校正値
type ChannelCal struct {
	BaseAD float64
	CalAD  float64
	EuHi   float64
	EuLo   float64
}
