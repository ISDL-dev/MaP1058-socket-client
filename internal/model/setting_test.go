package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToCSVHeader(t *testing.T) {
	type fields struct {
		TrendRange   TrendRange
		AnalysisType AnalysisType
		Calibration  Calibration
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "generate CSV header",
			fields: fields{
				TrendRange:   TrendRange{},
				AnalysisType: AnalysisType{},
				Calibration:  Calibration{},
			},
			want: []string{
				"ch_num",
				"analysis_type",
				"trend_range_upper",
				"trend_range_lower",
				"cal_base_ad",
				"cal_ad",
				"cal_eu_hi",
				"cal_eu_lo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{
				TrendRange:   tt.fields.TrendRange,
				AnalysisType: tt.fields.AnalysisType,
				Calibration:  tt.fields.Calibration,
			}
			assert.Equalf(t, tt.want, s.ToCSVHeader(), "ToCSVHeader()")
		})
	}
}

func TestToCSVRow(t *testing.T) {
	type fields struct {
		TrendRange   TrendRange
		AnalysisType AnalysisType
		Calibration  Calibration
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]string
	}{
		{
			name: "generate CSV row",
			fields: fields{
				TrendRange: TrendRange{
					{Upper: 1, Lower: 2},
					{Upper: 3, Lower: 4},
					{Upper: 5, Lower: 6},
					{Upper: 7, Lower: 8},
					{Upper: 9, Lower: 10},
					{Upper: 11, Lower: 12},
					{Upper: 13, Lower: 14},
					{Upper: 15, Lower: 16},
				},
				AnalysisType: AnalysisType{
					EMGCh,
					EOGHCh,
					EOGVCh,
					EEGCh,
					RRIntCh,
					CICh,
					NoAnalysis,
					NoAnalysis,
				},
				Calibration: Calibration{
					{BaseAD: 1, CalAD: 2, EuHi: 3, EuLo: 4},
					{BaseAD: 5, CalAD: 6, EuHi: 7, EuLo: 8},
					{BaseAD: 9, CalAD: 10, EuHi: 11, EuLo: 12},
					{BaseAD: 13, CalAD: 14, EuHi: 15, EuLo: 16},
					{BaseAD: 17, CalAD: 18, EuHi: 19, EuLo: 20},
					{BaseAD: 21, CalAD: 22, EuHi: 23, EuLo: 24},
					{BaseAD: 25, CalAD: 26, EuHi: 27, EuLo: 28},
					{BaseAD: 29, CalAD: 30, EuHi: 31, EuLo: 32},
				},
			},
			want: [][]string{
				{"1", "EMGCh", "1", "2", "1", "2", "3", "4"},
				{"2", "EOGHCh", "3", "4", "5", "6", "7", "8"},
				{"3", "EOGVCh", "5", "6", "9", "10", "11", "12"},
				{"4", "EEGCh", "7", "8", "13", "14", "15", "16"},
				{"5", "RRIntCh", "9", "10", "17", "18", "19", "20"},
				{"6", "CICh", "11", "12", "21", "22", "23", "24"},
				{"7", "NoAnalysis", "13", "14", "25", "26", "27", "28"},
				{"8", "NoAnalysis", "15", "16", "29", "30", "31", "32"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{
				TrendRange:   tt.fields.TrendRange,
				AnalysisType: tt.fields.AnalysisType,
				Calibration:  tt.fields.Calibration,
			}
			assert.Equalf(t, tt.want, s.ToCSVRows(), "ToCSVRow()")
		})
	}
}
