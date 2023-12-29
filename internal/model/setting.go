package model

import (
	"fmt"
)

type Setting struct {
	TrendRange   TrendRange
	AnalysisType AnalysisType
	Calibration  Calibration
}

func NewSetting() *Setting {
	return &Setting{
		TrendRange:   TrendRange{},
		AnalysisType: AnalysisType{},
		Calibration:  Calibration{},
	}
}

func (s *Setting) ToCSVHeader() []string {
	var header []string
	header = append(header, []string{
		"ch_num",
		"analysis_type",
		"trend_range_upper",
		"trend_range_lower",
		"cal_base_ad",
		"cal_ad",
		"cal_eu_hi",
		"cal_eu_lo"}...)
	return header
}

func (s *Setting) ToCSVRows() [][]string {
	var row [][]string
	for i, tr := range s.TrendRange {
		row = append(row, []string{
			fmt.Sprint(i + 1),
			fmt.Sprint(s.AnalysisType[i].String()),
			fmt.Sprint(tr.Upper),
			fmt.Sprint(tr.Lower),
			fmt.Sprint(s.Calibration[i].BaseAD),
			fmt.Sprint(s.Calibration[i].CalAD),
			fmt.Sprint(s.Calibration[i].EuHi),
			fmt.Sprint(s.Calibration[i].EuLo)})
	}
	return row
}

// TrendRange is an array that stores the upper and lower limits of the trend graph display range for all available channels.
type TrendRange [NumAvailableChs]ChannelRange

// ChannelRange is a struct for storing the upper and lower limits of the trend graph display range for each channel.
type ChannelRange struct {
	Upper uint
	Lower uint
}

// AnalysisType stores the analysis type for all available channels.
type AnalysisType [NumAvailableChs]ChannelType

// ChannelType shows the analysis type for each channel.
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

// String returns the name of the analysis type. If the analysis type is not defined, it returns an empty string.
func (c *ChannelType) String() string {
	switch *c {
	case NoAnalysis:
		return "NoAnalysis"
	case EMGCh:
		return "EMGCh"
	case EOGHCh:
		return "EOGHCh"
	case EOGVCh:
		return "EOGVCh"
	case EEGCh:
		return "EEGCh"
	case RRIntCh:
		return "RRIntCh"
	case CICh:
		return "CICh"
	case RRCVCh:
		return "RRCVCh"
	case HELFCh:
		return "HELFCh"
	case WAVECh:
		return "WAVECh"
	case ExtPlsCh:
		return "ExtPlsCh"
	case RespCh:
		return "RespCh"
	case Resp2Ch:
		return "Resp2Ch"
	default:
		return ""
	}
}

// Calibration stores the calibration values for all available channels.
type Calibration [NumAvailableChs]ChannelCal

// ChannelCal stores the calibration values for each channel.
type ChannelCal struct {
	BaseAD uint16
	CalAD  uint16
	EuHi   float64
	EuLo   float64
}
