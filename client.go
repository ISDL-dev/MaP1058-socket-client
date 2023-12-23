package MaP1058_socket_client

import (
	"context"
	"fmt"
	"github.com/Be3751/MaP1058-socket-client/internal/utils/net"
	"os"
	"sync"
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

	sc := scanner.NewCustomScanner(txtAdConn)
	txtAdapter := adapter.NewTxtAdapter(txtAdConn, sc, ps)

	return &client{
		binAdapter,
		txtAdapter,
		sgFile,
		c,
	}, nil
}

func (c *client) Start(rec time.Duration) error {
	var err error
	err = c.txt.StartRec(rec, time.Now())
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// via rcvSuccess, the goroutine for trend data tells the other whether it successfully receives data.
	var rcvSuccess chan bool
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		wg.Add(1)
		if bErr := c.bin.WriteRawSignal(ctx, rcvSuccess, setting); bErr != nil {
			err = fmt.Errorf("failed to write raw signal: %w", bErr)
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		wg.Add(1)
		if tErr := c.txt.WriteTrendData(ctx, rcvSuccess, adapter.CSVWriterGroup{}, setting.AnalysisType); tErr != nil {
			err = fmt.Errorf("failed to write trend data: %w", tErr)
			cancel()
		}
	}()

	wg.Wait()
	return err
}

func (c *client) Stop() error {
	if err := c.txt.EndRec(); err != nil {
		return fmt.Errorf("failed to end recording: %w", err)
	}
	return nil
}
