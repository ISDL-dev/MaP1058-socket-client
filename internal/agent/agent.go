package agent

import (
	"context"
	"time"

	"github.com/Be3751/socket-capture-signals/internal/socket"
)

type Agent interface {
	GetSignals(ctx context.Context) ([]byte, error)
}

var _ Agent = (*agent)(nil)

type agent struct {
	SocketClient socket.SocketClient
}

func NewAgent(s socket.SocketClient) Agent {
	return &agent{
		SocketClient: s,
	}
}

func (a *agent) GetSignals(ctx context.Context) ([]byte, error) {
	// TODO: YYYY/MM/DD形式が
	currentTime := time.Now()
	err := a.SocketClient.StartRec(ctx, time.Second*10, currentTime.String())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
