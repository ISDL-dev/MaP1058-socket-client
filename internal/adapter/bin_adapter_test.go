package adapter

import (
	"context"
	"testing"

	"github.com/ISDL-dev/MaP1058-socket-client/internal/model"
	my_parser "github.com/ISDL-dev/MaP1058-socket-client/internal/parser"
	mock_parser "github.com/ISDL-dev/MaP1058-socket-client/internal/parser/mock"
	mock_socket "github.com/ISDL-dev/MaP1058-socket-client/internal/socket/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReceiveADValues(t *testing.T) {
	t.Run("エラーなくAD値を受信する", func(t *testing.T) {
		t.Skip("実環境で動作検証する")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		socket := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		ctx, cancel := context.WithCancel(context.Background())
		binAdapter := NewBinAdapter(socket, parser, nil)

		buf := make([]byte, 1604)
		socket.EXPECT().Read(buf).Return(1604, nil)
		parser.EXPECT().ToSignals(buf).Return(&model.Signals{}, nil)
		socket.EXPECT().Write([]byte{0x06}).DoAndReturn(func(p []byte) (int, error) {
			cancel()
			return 1, nil
		})

		var rcvSuccess chan bool
		err := binAdapter.WriteRawSignal(ctx, rcvSuccess, nil)
		assert.NoError(t, err)
	})

	t.Run("サムチェックに1度失敗してリトライする", func(t *testing.T) {
		t.Skip("実環境で動作検証する")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		socket := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		ctx, cancel := context.WithCancel(context.Background())
		binAdapter := NewBinAdapter(socket, parser, nil)

		buf := make([]byte, 1604)
		socket.EXPECT().Read(buf).Return(1604, nil)
		parser.EXPECT().ToSignals(buf).Return(nil, &my_parser.FailureSumCheckError{Expected: 100, Actual: 10})
		socket.EXPECT().Write([]byte{0x15}).Return(1, nil)
		socket.EXPECT().Read(buf).Return(1604, nil)
		parser.EXPECT().ToSignals(buf).Return(&model.Signals{}, nil)
		socket.EXPECT().Write([]byte{0x6}).DoAndReturn(func(p []byte) (int, error) {
			cancel()
			return 1, nil
		})

		var rcvSuccess chan bool
		err := binAdapter.WriteRawSignal(ctx, rcvSuccess, nil)
		assert.NoError(t, err)
	})
}
