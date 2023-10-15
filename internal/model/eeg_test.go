package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyzedEEG_ToCSVHeader(t *testing.T) {
	type args struct {
		at AnalysisType
	}
	tests := []struct {
		name string
		a    AnalyzedEEG
		args args
		want []string
	}{
		{
			name: "ch1~4がEEGの場合",
			a:    AnalyzedEEG{},
			args: args{
				at: AnalysisType{
					EEGCh, EEGCh, EEGCh, EEGCh,
					NoAnalysis, NoAnalysis, NoAnalysis, NoAnalysis,
				},
			},
			want: []string{
				"", "total_ch1", "delta_ch1", "theta_ch1", "alpha_ch1", "beta_ch1", "ratio_ch1",
				"total_ch2", "delta_ch2", "theta_ch2", "alpha_ch2", "beta_ch2", "ratio_ch2",
				"total_ch3", "delta_ch3", "theta_ch3", "alpha_ch3", "beta_ch3", "ratio_ch3",
				"total_ch4", "delta_ch4", "theta_ch4", "alpha_ch4", "beta_ch4", "ratio_ch4",
			},
		},
		{
			name: "ch4~6, ch8がEEGの場合",
			a:    AnalyzedEEG{},
			args: args{
				at: AnalysisType{
					NoAnalysis, NoAnalysis, NoAnalysis, EEGCh,
					EEGCh, EEGCh, NoAnalysis, EEGCh,
				},
			},
			want: []string{
				"", "total_ch4", "delta_ch4", "theta_ch4", "alpha_ch4", "beta_ch4", "ratio_ch4",
				"total_ch5", "delta_ch5", "theta_ch5", "alpha_ch5", "beta_ch5", "ratio_ch5",
				"total_ch6", "delta_ch6", "theta_ch6", "alpha_ch6", "beta_ch6", "ratio_ch6",
				"total_ch8", "delta_ch8", "theta_ch8", "alpha_ch8", "beta_ch8", "ratio_ch8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.ToCSVHeader(tt.args.at), "ToCSVHeader(%v)", tt.args.at)
		})
	}
}

func TestAnalyzedEEG_ToCSVRow(t *testing.T) {
	tests := []struct {
		name string
		a    AnalyzedEEG
		want []string
	}{
		{
			name: "ch1~4がEEGの場合",
			a: AnalyzedEEG{
				{
					{Time: 1, ChNum: 1, BandNum: Total, Power: 1.1, MaxEEG: 1.2, MinEEG: 1.3},
					{Time: 1, ChNum: 1, BandNum: Delta, Power: 2.1, MaxEEG: 2.2, MinEEG: 2.3},
					{Time: 1, ChNum: 1, BandNum: Theta, Power: 3.1, MaxEEG: 3.2, MinEEG: 3.3},
					{Time: 1, ChNum: 1, BandNum: Alpha, Power: 4.1, MaxEEG: 4.2, MinEEG: 4.3},
					{Time: 1, ChNum: 1, BandNum: Beta, Power: 5.1, MaxEEG: 5.2, MinEEG: 5.3},
					{Time: 1, ChNum: 1, BandNum: Ratio, Power: 6.1, MaxEEG: 6.2, MinEEG: 6.3},
				},
				{
					{Time: 1, ChNum: 2, BandNum: Total, Power: 1.1, MaxEEG: 1.2, MinEEG: 1.3},
					{Time: 1, ChNum: 2, BandNum: Delta, Power: 2.1, MaxEEG: 2.2, MinEEG: 2.3},
					{Time: 1, ChNum: 2, BandNum: Theta, Power: 3.1, MaxEEG: 3.2, MinEEG: 3.3},
					{Time: 1, ChNum: 2, BandNum: Alpha, Power: 4.1, MaxEEG: 4.2, MinEEG: 4.3},
					{Time: 1, ChNum: 2, BandNum: Beta, Power: 5.1, MaxEEG: 5.2, MinEEG: 5.3},
					{Time: 1, ChNum: 2, BandNum: Ratio, Power: 6.1, MaxEEG: 6.2, MinEEG: 6.3},
				},
				{
					{Time: 1, ChNum: 3, BandNum: Total, Power: 1.1, MaxEEG: 1.2, MinEEG: 1.3},
					{Time: 1, ChNum: 3, BandNum: Delta, Power: 2.1, MaxEEG: 2.2, MinEEG: 2.3},
					{Time: 1, ChNum: 3, BandNum: Theta, Power: 3.1, MaxEEG: 3.2, MinEEG: 3.3},
					{Time: 1, ChNum: 3, BandNum: Alpha, Power: 4.1, MaxEEG: 4.2, MinEEG: 4.3},
					{Time: 1, ChNum: 3, BandNum: Beta, Power: 5.1, MaxEEG: 5.2, MinEEG: 5.3},
					{Time: 1, ChNum: 3, BandNum: Ratio, Power: 6.1, MaxEEG: 6.2, MinEEG: 6.3},
				},
				{
					{Time: 1, ChNum: 4, BandNum: Total, Power: 1.1, MaxEEG: 1.2, MinEEG: 1.3},
					{Time: 1, ChNum: 4, BandNum: Delta, Power: 2.1, MaxEEG: 2.2, MinEEG: 2.3},
					{Time: 1, ChNum: 4, BandNum: Theta, Power: 3.1, MaxEEG: 3.2, MinEEG: 3.3},
					{Time: 1, ChNum: 4, BandNum: Alpha, Power: 4.1, MaxEEG: 4.2, MinEEG: 4.3},
					{Time: 1, ChNum: 4, BandNum: Beta, Power: 5.1, MaxEEG: 5.2, MinEEG: 5.3},
					{Time: 1, ChNum: 4, BandNum: Ratio, Power: 6.1, MaxEEG: 6.2, MinEEG: 6.3},
				},
				{
					{0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0},
				},
			},
			want: []string{
				"1", "1.10000", "2.10000", "3.10000", "4.10000", "5.10000", "6.10000",
				"1.10000", "2.10000", "3.10000", "4.10000", "5.10000", "6.10000",
				"1.10000", "2.10000", "3.10000", "4.10000", "5.10000", "6.10000",
				"1.10000", "2.10000", "3.10000", "4.10000", "5.10000", "6.10000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.ToCSVRow(), "ToCSVRow()")
		})
	}
}
