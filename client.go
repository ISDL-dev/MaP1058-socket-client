package MaP1058_socket_client

import (
	"context"
	"fmt"
	"github.com/Be3751/MaP1058-socket-client/internal/utils/net"
	"os"
	"time"

	"github.com/Be3751/MaP1058-socket-client/internal/adapter"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/scanner"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
)

type Client interface {
	// Start sends a command to start recording and receives some commands
	// containing signals and measurement conditions.
	Start(rec time.Duration) error
	// Stop sends a command to stop recording.
	Stop() error
}

type client struct {
	bin    adapter.BinAdapter
	txt    adapter.TxtAdapter
	raw    *os.File
	ctx    context.Context
	cancel context.CancelFunc
	config Config
}

var _ Client = (*client)(nil)

// Config is the configuration for MaP1058 client
type Config struct {
	// ServerIP is the IP address of MaP1058.
	ServerIP string
	// SaveDir is the directory to save received signals and configuration.
	SaveDir string
}

func NewClient(c Config) (Client, error) {
	clientIP, err := net.GetMyLocalIP()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// make TCP/IP connection for binary data and text data
	binAdConf := socket.SocketConfig{
		ServerIP:   c.ServerIP,
		ServerPort: "2200",
		ClientIP:   clientIP,
		ClientPort: 1200,
	}
	binAdConn, err := socket.Connect(binAdConf)
	if err != nil {
		panic(err)
	}

	txtAdConf := socket.SocketConfig{
		ServerIP:   c.ServerIP,
		ServerPort: "3000",
		ClientIP:   clientIP,
		ClientPort: 1100,
	}
	txtAdConn, err := socket.Connect(txtAdConf)
	if err != nil {
		panic(err)
	}

	// make file for saving raw signal
	sgFilePath := fmt.Sprintf("%s/rawwave_%s.csv", c.SaveDir, time.Now().Format("20060102150405"))
	sgFile, err := os.Create(sgFilePath)
	if err != nil {
		panic(err)
	}

	ps := parser.NewParser()
	binAdapter := adapter.NewBinAdapter(binAdConn, ps, sgFile)

	ctx, cancel := context.WithCancel(context.Background())
	sc := scanner.NewCustomScanner(txtAdConn)
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, sc, ps)

	return &client{
		binAdapter,
		txtAdapter,
		sgFile,
		ctx,
		cancel,
		c,
	}, nil
}

func (c *client) Start(rec time.Duration) error {
	err := c.txt.StartRec(rec, time.Now())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = c.txt.EndRec()
		if err != nil {
			panic(err)
		}
	}()

	setting, err := c.txt.GetSetting()
	if err != nil {
		return fmt.Errorf("failed to get setting: %w", err)
	}
	// TODO: write setting to file

	var bErrChan chan error
	go func() {
		err := c.bin.WriteRawSignal(c.ctx, setting)
		if err != nil {
			bErrChan <- fmt.Errorf("failed to write raw signal: %w", err)
		}
	}()

	var tErrChan chan error
	go func() {
		err = c.txt.WriteTrendData(c.ctx, adapter.CSVWriterGroup{}, setting.AnalysisType)
		if err != nil {
			tErrChan <- fmt.Errorf("failed to write trend data: %w", err)
		}
	}()

	for i := 0; i < 2; i++ {
		select {
		case err := <-bErrChan:
			if err != nil {
				return err
			}
			if err = c.raw.Close(); err != nil {
				return fmt.Errorf("failed to close raw file: %w", err)
			}
		case err := <-tErrChan:
			if err != nil {
				return err
			}
			c.cancel()
		}
	}

	return nil
}

func (c *client) Stop() error {
	if err := c.txt.EndRec(); err != nil {
		return fmt.Errorf("failed to end recording: %w", err)
	}
	return nil
}
