package adapter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	mock_parser "github.com/Be3751/MaP1058-socket-client/internal/parser/mock"
	mock_scanner "github.com/Be3751/MaP1058-socket-client/internal/scanner/mock"
	mock_socket "github.com/Be3751/MaP1058-socket-client/internal/socket/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStartRec(t *testing.T) {
	t.Run("収録を開始する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		conn := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		scanner := mock_scanner.NewMockCustomScanner(ctrl)

		sCmdStr := "<SCMD>START:A:300,2023/01/01 12-00-00,,,,,,,,</SCMD>"
		conn.EXPECT().Write(gomock.Any()).DoAndReturn(
			func(buf []byte) (int, error) {
				assert.Equal(t, []byte(sCmdStr), buf)
				return len(sCmdStr), nil
			},
		)
		conn.EXPECT().Read(gomock.Any()).SetArg(0, []byte(sCmdStr)).Return(len(sCmdStr), nil)
		ctx := context.Background()

		txtAdapter := NewTxtAdapter(conn, scanner, parser)
		err := txtAdapter.StartRec(ctx, time.Second*300, time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local))
		assert.NoError(t, err)
	})
}

func TestEndRec(t *testing.T) {
	t.Run("収録を終了する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		conn := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		scanner := mock_scanner.NewMockCustomScanner(ctrl)

		sCmdStr := "<SCMD>END:A:,,,,,,,,,</SCMD>"
		conn.EXPECT().Write(gomock.Any()).DoAndReturn(
			func(buf []byte) (int, error) {
				assert.Equal(t, []byte(sCmdStr), buf)
				return len(sCmdStr), nil
			},
		)
		conn.EXPECT().Read(gomock.Any()).SetArg(0, []byte(sCmdStr)).Return(len(sCmdStr), nil)
		ctx := context.Background()

		txtAdapter := NewTxtAdapter(conn, scanner, parser)
		err := txtAdapter.EndRec(ctx)
		assert.NoError(t, err)
	})
}

func TestGetStatus(t *testing.T) {
	t.Run("ステータスを取得する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		conn := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		scanner := mock_scanner.NewMockCustomScanner(ctrl)

		sCmdStr := "<SCMD>STATUS:A:,,,,,,,,,</SCMD>"
		conn.EXPECT().Write(gomock.Any()).DoAndReturn(
			func(buf []byte) (int, error) {
				assert.Equal(t, sCmdStr, string(buf))
				return len(sCmdStr), nil
			},
		)
		rCmd := []byte("<SCMD>STATUS:A:Acq,,,,,,,,,</SCMD>")
		conn.EXPECT().Read(gomock.Any()).SetArg(0, rCmd).Return(len(rCmd), nil)
		ctx := context.Background()
		parser.EXPECT().ToCommand(string(rCmd)).Return(model.Command{Name: "STATUS", Params: [10]string{"Acq"}}, nil)

		txtAdapter := NewTxtAdapter(conn, scanner, parser)
		status, err := txtAdapter.GetStatus(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "Acq", string(status))
	})
}

func TestGetSetting(t *testing.T) {
	t.Run("設定を取得する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		conn := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		scanner := mock_scanner.NewMockCustomScanner(ctrl)

		var order []*gomock.Call
		order = append(order, scanner.EXPECT().Scan().Return(true))
		order = append(order, scanner.EXPECT().Text().Return("<SCMD>RANGE:A:1200;600,100;0,300;0,2;0,200;0,100;0,6;0,5;0,,</SCMD>"))
		order = append(order, parser.EXPECT().ToCommand("<SCMD>RANGE:A:1200;600,100;0,300;0,2;0,200;0,100;0,6;0,5;0,,</SCMD>").Return(model.Command{
			Name:   "RANGE",
			Params: [10]string{"1200;600", "100;0", "300;0", "2;0", "200;0", "100;0", "6;0", "5;0", "", ""},
		}, nil))
		order = append(order, parser.EXPECT().ToTrendRange(model.Command{
			Name:   "RANGE",
			Params: [10]string{"1200;600", "100;0", "300;0", "2;0", "200;0", "100;0", "6;0", "5;0", "", ""},
		}).Return(
			model.TrendRange{
				{Upper: 1200, Lower: 600},
				{Upper: 100, Lower: 0},
				{Upper: 300, Lower: 0},
				{Upper: 2, Lower: 0},
				{Upper: 200, Lower: 0},
				{Upper: 100, Lower: 0},
				{Upper: 6, Lower: 0},
				{Upper: 5, Lower: 0},
			},
			nil))
		order = append(order, scanner.EXPECT().Scan().Return(true))
		order = append(order, scanner.EXPECT().Text().Return("<SCMD>ANALYSIS:A:5,4,4,4,4,3,11,9,,</SCMD>"))
		order = append(order, parser.EXPECT().ToCommand("<SCMD>ANALYSIS:A:5,4,4,4,4,3,11,9,,</SCMD>").Return(model.Command{
			Name:   "ANALYSIS",
			Params: [10]string{"5", "4", "4", "4", "4", "3", "11", "9", "", ""},
		}, nil))
		order = append(order, parser.EXPECT().ToAnalysis(model.Command{
			Name:   "ANALYSIS",
			Params: [10]string{"5", "4", "4", "4", "4", "3", "11", "9", "", ""},
		}).Return(
			model.Analysis{5, 4, 4, 4, 4, 3, 11, 9},
			nil,
		))

		cal := []model.ChannelCal{
			{BaseAD: 0, CalAD: 409, EuHi: 0.05, EuLo: 0},
			{BaseAD: 1, CalAD: 317, EuHi: 0.03, EuLo: 1},
			{BaseAD: 2, CalAD: 298, EuHi: 0.13, EuLo: 2},
			{BaseAD: 3, CalAD: 101, EuHi: 0.05, EuLo: 3},
			{BaseAD: 4, CalAD: 201, EuHi: 0.03, EuLo: 4},
			{BaseAD: 5, CalAD: 102, EuHi: 0.13, EuLo: 5},
			{BaseAD: 6, CalAD: 202, EuHi: 0.05, EuLo: 6},
			{BaseAD: 7, CalAD: 302, EuHi: 0.03, EuLo: 7},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
			{BaseAD: 0, CalAD: 0, EuHi: 0, EuLo: 0},
		}
		for i, c := range cal {
			stgCmdStr := fmt.Sprintf(
				"<SCMD>GETSETTING:A:\"CH%d,BASE_AD=%d,CAL_AD=%d,EU_HI=%f,EU_LO=%f,\",,,,,,,,,</SCMD>",
				i+1, c.BaseAD, c.CalAD, c.EuHi, c.EuLo)
			stgCmd := model.Command{
				Name: "GETSETTING",
				Params: [10]string{
					fmt.Sprintf("CH%d", i+1),
					fmt.Sprintf("BASE_AD=%d", c.BaseAD),
					fmt.Sprintf("CAL_AD=%d", c.CalAD),
					fmt.Sprintf("EU_HI=%f", c.EuHi),
					fmt.Sprintf("EU_LO=%f", c.EuLo),
					"", "", "", "", "",
				},
			}
			order = append(order, scanner.EXPECT().Scan().Return(true))
			order = append(order, scanner.EXPECT().Text().Return(stgCmdStr))
			order = append(order, parser.EXPECT().ToCommand(stgCmdStr).Return(stgCmd, nil))
			if i < 8 {
				order = append(order, parser.EXPECT().ToChannelCal(stgCmd).Return(c, nil))
			}
		}
		order = append(order, scanner.EXPECT().Err().Return(nil))
		gomock.InOrder(order...)

		txtAdapter := NewTxtAdapter(conn, scanner, parser)
		setting, err := txtAdapter.GetSetting()
		assert.NoError(t, err)
		assert.Equal(t, model.Setting{
			TrendRange: model.TrendRange{
				{Upper: 1200, Lower: 600},
				{Upper: 100, Lower: 0},
				{Upper: 300, Lower: 0},
				{Upper: 2, Lower: 0},
				{Upper: 200, Lower: 0},
				{Upper: 100, Lower: 0},
				{Upper: 6, Lower: 0},
				{Upper: 5, Lower: 0},
			},
			Analysis: model.Analysis{5, 4, 4, 4, 4, 3, 11, 9},
			Calibration: [8]model.ChannelCal{
				{BaseAD: 0, CalAD: 409, EuHi: 0.05, EuLo: 0},
				{BaseAD: 1, CalAD: 317, EuHi: 0.03, EuLo: 1},
				{BaseAD: 2, CalAD: 298, EuHi: 0.13, EuLo: 2},
				{BaseAD: 3, CalAD: 101, EuHi: 0.05, EuLo: 3},
				{BaseAD: 4, CalAD: 201, EuHi: 0.03, EuLo: 4},
				{BaseAD: 5, CalAD: 102, EuHi: 0.13, EuLo: 5},
				{BaseAD: 6, CalAD: 202, EuHi: 0.05, EuLo: 6},
				{BaseAD: 7, CalAD: 302, EuHi: 0.03, EuLo: 7}},
		}, *setting)
	})

	t.Run("スキャンが可能なトークンを受信できずエラー", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		conn := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		scanner := mock_scanner.NewMockCustomScanner(ctrl)

		scanner.EXPECT().Scan().Return(false)
		scanner.EXPECT().Err().Return(fmt.Errorf("error"))

		txtAdapter := NewTxtAdapter(conn, scanner, parser)
		_, err := txtAdapter.GetSetting()
		assert.Error(t, err)
	})

	t.Run("ANALYSISを受信できずエラー", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		conn := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		scanner := mock_scanner.NewMockCustomScanner(ctrl)

		scanner.EXPECT().Scan().Return(true)
		scanner.EXPECT().Text().Return("<SCMD>RANGE:A:1200;600,100;0,300;0,2;0,200;0,100;0,6;0,5;0,,</SCMD>")
		parser.EXPECT().
			ToCommand("<SCMD>RANGE:A:1200;600,100;0,300;0,2;0,200;0,100;0,6;0,5;0,,</SCMD>").
			Return(model.Command{
				Name: "RANGE",
				Params: [10]string{
					"1200;600", "100;0", "300;0", "2;0", "200;0", "100;0", "6;0", "5;0", "", "",
				},
			}, nil)
		parser.EXPECT().ToTrendRange(model.Command{
			Name: "RANGE",
			Params: [10]string{
				"1200;600", "100;0", "300;0", "2;0", "200;0", "100;0", "6;0", "5;0", "", "",
			},
		}).Return(model.TrendRange{
			{Upper: 1200, Lower: 600},
			{Upper: 100, Lower: 0},
			{Upper: 300, Lower: 0},
			{Upper: 2, Lower: 0},
			{Upper: 200, Lower: 0},
			{Upper: 100, Lower: 0},
			{Upper: 6, Lower: 0},
			{Upper: 5, Lower: 0},
		}, nil)
		scanner.EXPECT().Scan().Return(true)
		scanner.EXPECT().Text().Return("<ERR>error</ERR>")
		parser.EXPECT().ToCommand("<ERR>error</ERR>").Return(model.Command{}, fmt.Errorf("error"))

		txtAdapter := NewTxtAdapter(conn, scanner, parser)
		_, err := txtAdapter.GetSetting()
		assert.Error(t, err)
	})
}
