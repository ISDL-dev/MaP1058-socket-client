package parser

import (
	"fmt"
	"github.com/ISDL-dev/MaP1058-socket-client/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCommand(t *testing.T) {
	t.Run("stringからCommandに変換", func(t *testing.T) {
		rCmdStr := "<SCMD>START:A:300,2023/01/01 12-00-00,,,,,,,,</SCMD>"
		parser := NewParser()
		cmd, err := parser.ToCommand(rCmdStr)
		assert.NoError(t, err)
		assert.Equal(t, "START", cmd.Name)
		expectedParams := [10]string{}
		expectedParams[0] = "300"
		expectedParams[1] = "2023/01/01 12-00-00"
		assert.Equal(t, expectedParams, cmd.Params)
	})

	t.Run("GETSETTINGコマンドのパース", func(t *testing.T) {
		rCmdStr := "<SCMD>GETSETTING:A:\"CH1,BASE_AD=0,CAL_AD=2048,EU_HI=2.5,EU_LO=0,\",,,,,,,,,</SCMD>"
		parser := NewParser()
		cmd, err := parser.ToCommand(rCmdStr)
		assert.NoError(t, err)
		assert.Equal(t, "GETSETTING", cmd.Name)
		expectedParams := [10]string{}
		expectedParams[0] = "CH1"
		expectedParams[1] = "BASE_AD=0"
		expectedParams[2] = "CAL_AD=2048"
		expectedParams[3] = "EU_HI=2.5"
		expectedParams[4] = "EU_LO=0"
		assert.Equal(t, expectedParams, cmd.Params)
	})

	t.Run("<SCMD>...</SCMD>の構造になっていない", func(t *testing.T) {
		rCmdStr := "<S>START:A:300,2023/01/01 12-00-00,,,,,,,,</S>"
		parser := NewParser()
		_, err := parser.ToCommand(rCmdStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "<SCMD> and </SCMD> on both sides")
	})

	t.Run("規定の数のカンマがない", func(t *testing.T) {
		rCmdStr := "<SCMD>START:A:300,2023/01/01 12-00-00,,,,</SCMD>"
		parser := NewParser()
		_, err := parser.ToCommand(rCmdStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "9 commas")
	})

	t.Run("コマンド名とパラメータを分ける:A:が含まれていない", func(t *testing.T) {
		rCmdStr := "<SCMD>START:300,2023/01/01 12-00-00,,,,,,,,</SCMD>"
		parser := NewParser()
		_, err := parser.ToCommand(rCmdStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ":A:")
	})
}

func Test_parser_ToChannelPower(t *testing.T) {
	type args struct {
		c model.Command
	}
	tests := []struct {
		name    string
		args    args
		want    model.ChannelPower
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "正常系",
			args: args{
				c: model.Command{
					Name: "DATA_EEG",
					Params: [model.NumSeparator + 1]string{
						"4", "215", "1", "3", "41.16424", "16.70227", "-17.6422", "", "", "",
					},
				},
			},
			want: model.ChannelPower{
				Time:    215,
				ChNum:   1,
				BandNum: model.Alpha,
				Power:   41.16424,
				MaxEEG:  16.70227,
				MinEEG:  -17.6422,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, err := p.ToChannelPower(tt.args.c)
			if !tt.wantErr(t, err, fmt.Sprintf("ToChannelPower(%v)", tt.args.c)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ToChannelPower(%v)", tt.args.c)
		})
	}
}

func TestToChannelCal(t *testing.T) {
	type args struct {
		c model.Command
	}
	tests := []struct {
		name    string
		args    args
		want    model.ChannelCal
		wantErr error
	}{
		{
			"受信に成功する",
			args{c: model.Command{
				Name:   "GETSETTING",
				Params: [10]string{"CH1", "BASE_AD=0", "CAL_AD=2048", "EU_HI=2.5", "EU_LO=0", "", "", "", "", ""},
			}},
			model.ChannelCal{BaseAD: 0, CalAD: 2048, EuHi: 2.5, EuLo: 0},
			nil,
		},
		{
			"不適切なコマンド名でエラー",
			args{c: model.Command{
				Name:   "GET-SETTING",
				Params: [10]string{"CH1", "BASE_AD=0", "CAL_AD=2048", "EU_HI=2.5", "EU_LO=0", "", "", "", "", ""},
			}},
			model.ChannelCal{},
			fmt.Errorf("the received command is not GETSETTING: GET-SETTING: %s", "GET-SETTING"),
		},
		{
			"パラメータ数が不適切でエラー",
			args{c: model.Command{
				Name:   "GETSETTING",
				Params: [10]string{"CH1", "BASE_AD=0", "CAL_AD=2048", "EU_HI=2.5", "EU_LO=0", "hoge", "", "", ""},
			}},
			model.ChannelCal{},
			fmt.Errorf("hoge"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{}
			got, err := p.ToChannelCal(tt.args.c)
			if err != nil {
				assert.Equalf(t, tt.want, got, "ToChannelCal(%v)", tt.args.c)
			}
			assert.Equalf(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}
