package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	mock_parser "github.com/Be3751/MaP1058-socket-client/internal/parser/mock"
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

		sCmdStr := "<SCMD>START:A:300,2023/01/01 12-00-00,,,,,,,,</SCMD>"
		conn.EXPECT().Write(gomock.Any()).DoAndReturn(
			func(buf []byte) (int, error) {
				assert.Equal(t, []byte(sCmdStr), buf)
				return len(sCmdStr), nil
			},
		)
		conn.EXPECT().Read(gomock.Any()).SetArg(0, []byte(sCmdStr)).Return(len(sCmdStr), nil)
		ctx := context.Background()

		txtAdapter := NewTxtAdapter(conn, parser)
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

		sCmdStr := "<SCMD>END:A:,,,,,,,,,</SCMD>"
		conn.EXPECT().Write(gomock.Any()).DoAndReturn(
			func(buf []byte) (int, error) {
				assert.Equal(t, []byte(sCmdStr), buf)
				return len(sCmdStr), nil
			},
		)
		conn.EXPECT().Read(gomock.Any()).SetArg(0, []byte(sCmdStr)).Return(len(sCmdStr), nil)
		ctx := context.Background()

		txtAdapter := NewTxtAdapter(conn, parser)
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
		parser.EXPECT().ToCommand(string(rCmd)).Return(&model.Command{Name: "STATUS", Params: [10]string{"Acq"}}, nil)

		txtAdapter := NewTxtAdapter(conn, parser)
		status, err := txtAdapter.GetStatus(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "Acq", string(status))
	})
}
