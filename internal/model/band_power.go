package model

type EEG struct {
	Time    uint
	ChNum   uint
	BandNum BandNum
	Power   float64
	MaxEEG  float64
	MinEEG  float64
}

type BandNum uint

const (
	Total BandNum = iota
	Delta
	Theta
	Alpha
	Beta
	Ratio
)
