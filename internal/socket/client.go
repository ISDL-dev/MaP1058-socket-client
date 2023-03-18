package socket

import "context"

type SocketClient interface {
	GetSignals(ctx context.Context) ([]byte, error)
}

var _ SocketClient = (*client)(nil)

func NewSocketClient() SocketClient {
	return &client{}
}

type client struct {
}

func (c *client) GetSignals(ctx context.Context) ([]byte, error) {
	return nil, nil
}
