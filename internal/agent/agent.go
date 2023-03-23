package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/Be3751/socket-capture-signals/internal/pkg/socket"
)

type Agent interface {
	StartRec(ctx context.Context, recTime time.Duration, recDateTime string) error
	EndRec(ctx context.Context) error
	CaptureSig(ctx context.Context, buf []byte) error
	// GetStat(ctx context.Context) ([]byte, error)
	// Mark(ctx context.Context) error
	// UnMark(ctx context.Context) error
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

func (a *agent) StartRec(ctx context.Context, recTime time.Duration, recDateTime string) error {
	sCmd := fmt.Sprintf("<SCMD>START:A:%d,\"%s\",,,,,,,,,</SCMD>", recTime, recDateTime)
	err := a.SocketClient.SendCmd(ctx, sCmd)
	if err != nil {
		return fmt.Errorf("failed to send command %w", err)
	}

	rCmd, err := a.SocketClient.ReceiveCmd(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	if rCmd != sCmd {
		return fmt.Errorf("failed to start recording")
	}

	return nil
}

func (a *agent) EndRec(ctx context.Context) error {
	sCmd := "<SCMD>END:A:,,,,,,,,,</SCMD>"
	err := a.SocketClient.SendCmd(ctx, sCmd)
	if err != nil {
		return fmt.Errorf("failed to start recording %w", err)
	}

	rCmd, err := a.SocketClient.ReceiveCmd(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive command %w", err)
	}
	if rCmd != sCmd {
		return fmt.Errorf("failed to end recording")
	}

	return nil
}

func (a *agent) CaptureSig(ctx context.Context, buf []byte) error {
	sCmd := "<SCMD>DATA_EEG:A:,,,,,,,,,</SCMD>"
	err := a.SocketClient.SendCmd(ctx, sCmd)
	if err != nil {
		return fmt.Errorf("failed to capture signals %w", err)
	}
	return nil
}
