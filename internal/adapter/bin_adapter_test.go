package adapter

import (
	"context"
	"testing"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
	my_parser "github.com/Be3751/MaP1058-socket-client/internal/parser"
	mock_parser "github.com/Be3751/MaP1058-socket-client/internal/parser/mock"
	mock_socket "github.com/Be3751/MaP1058-socket-client/internal/socket/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReceiveADValues(t *testing.T) {
	t.Run("エラーなくAD値を受信する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		socket := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		ctx := context.Background()
		binAdapter := NewBinAdapter(socket, parser)

		buf := make([]byte, 1604)
		socket.EXPECT().Read(buf).Return(1604, nil)
		parser.EXPECT().ToSignals(buf).Return(&model.Signals{}, nil)
		socket.EXPECT().Write([]byte{0x06}).Return(1, nil)

		signals, err := binAdapter.ReceiveADValues(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, signals)
	})

	t.Run("サムチェックに1度失敗してリトライする", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		socket := mock_socket.NewMockConn(ctrl)
		parser := mock_parser.NewMockParser(ctrl)
		ctx := context.Background()
		binAdapter := NewBinAdapter(socket, parser)

		buf := make([]byte, 1604)
		socket.EXPECT().Read(buf).Return(1604, nil)
		parser.EXPECT().ToSignals(buf).Return(nil, &my_parser.FailureSumCheckError{Expected: 100, Actual: 10})
		socket.EXPECT().Write([]byte{0x15}).Return(1, nil)
		socket.EXPECT().Read(buf).Return(1604, nil)
		parser.EXPECT().ToSignals(buf).Return(&model.Signals{}, nil)
		socket.EXPECT().Write([]byte{0x6}).Return(1, nil)

		signals, err := binAdapter.ReceiveADValues(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, signals)
	})
}
