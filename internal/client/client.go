package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/adapter"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
	myNet "github.com/Be3751/MaP1058-socket-client/utils/net"
)

type Client interface {
	Start(rec time.Duration) error
	Stop() error
}

type Config struct {
	GetRaw      bool
	GetAnalyzed bool
	SaveDir     string
}

func NewClient(c Config) (Client, error) {
	clientIP, err := myNet.GetMyLocalIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get client IP address: %w", err)
	}
	txtAdConf := socket.SocketConfig{
		ServerIP:   "192.168.10.101",
		ServerPort: "3000",
		ClientIP:   clientIP,
		ClientPort: 1100,
	}
	txtAdConn, err := socket.Connect(txtAdConf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to txt adapter: %w", err)
	}

	binAdConf := socket.SocketConfig{
		ServerIP:   "192.168.10.101",
		ServerPort: "2200",
		ClientIP:   clientIP,
		ClientPort: 1200,
	}
	binAdConn, err := socket.Connect(binAdConf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to bin adapter: %w", err)
	}

	ctx := context.Background()
	parser := parser.NewParser()
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, parser)
	binAdapter := adapter.NewBinAdapter(binAdConn, parser)
	return &client{
		bin:    binAdapter,
		txt:    txtAdapter,
		ctx:    ctx,
		config: c,
	}, nil
}

type client struct {
	bin    adapter.BinAdapter
	txt    adapter.TxtAdapter
	ctx    context.Context
	config Config
}

func (c *client) Start(rec time.Duration) error {
	err := c.txt.StartRec(c.ctx, rec, time.Now())
	if err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}
	defer func() {
		// TODO: 既に終了している場合は、終了処理をスキップする
	}()

	// TODO: 各種設定値を取得する
	setFile, err := os.Create(fmt.Sprintf("%s/settings.csv", c.config.SaveDir))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	// TODO: Stop()が実行されるまで受信値を出力し続ける
	sig, err := c.bin.ReceiveADValues(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to receive AD values: %w", err)
	}

	rawFile, err := os.Create(fmt.Sprintf("%s/raw_data.csv", c.config.SaveDir))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

func (c *client) Stop() error {
	panic("implement me")
}
