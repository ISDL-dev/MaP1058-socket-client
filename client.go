package MaP1058_socket_client

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ISDL-dev/MaP1058-socket-client/internal/model"
	"github.com/ISDL-dev/MaP1058-socket-client/internal/utils/net"
	"os"
	"sync"
	"time"

	"github.com/ISDL-dev/MaP1058-socket-client/internal/adapter"
	"github.com/ISDL-dev/MaP1058-socket-client/internal/parser"
	"github.com/ISDL-dev/MaP1058-socket-client/internal/scanner"
	"github.com/ISDL-dev/MaP1058-socket-client/internal/socket"
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
		return nil, fmt.Errorf("failed to get the client IP address: %w", err)
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
		return nil, fmt.Errorf("failed to make a TCP connection for raw wave: %w", err)
	}

	txtAdConf := socket.SocketConfig{
		ServerIP:   c.ServerIP,
		ServerPort: "3000",
		ClientIP:   clientIP,
		ClientPort: 1100,
	}
	txtAdConn, err := socket.Connect(txtAdConf)
	if err != nil {
		return nil, fmt.Errorf("failed to make a TCP connection for trend data: %w", err)
	}

	// make file for saving raw signal
	sgFilePath := fmt.Sprintf("%s/rawwave_%s.csv", c.SaveDir, time.Now().Format("20060102150405"))
	sgFile, err := os.Create(sgFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create a file of experiment setting: %w", err)
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
		return fmt.Errorf("failed to start recording: %w", err)
	}

	setting, err := c.txt.GetSetting()
	if err != nil {
		return fmt.Errorf("failed to get setting: %w", err)
	}
	if err := WriteSetting(setting, c.config.SaveDir); err != nil {
		return fmt.Errorf("failed to write setting: %w", err)
	}

	//ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), rec)
	defer cancel()

	// via rcvSuccess, the goroutine for trend data tells the other whether it successfully receives data.
	rcvSuccess := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if bErr := c.bin.WriteRawSignal(ctx, rcvSuccess, setting); bErr != nil {
			err = fmt.Errorf("failed to write raw signal: %w", bErr)
			cancel()
		}
		fmt.Println("finish WriteRawSignal")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if tErr := c.txt.WriteTrendData(ctx, rcvSuccess, adapter.CSVWriterGroup{}, setting.AnalysisType); tErr != nil {
			err = fmt.Errorf("failed to write trend data: %w", tErr)
			cancel()
		}
		fmt.Println("finish WriteTrendData")
	}()

	fmt.Println("Now receiving...")
	wg.Wait()
	return err
}

func (c *client) Stop() error {
	if err := c.txt.EndRec(); err != nil {
		return fmt.Errorf("failed to end recording: %w", err)
	}
	return nil
}

func WriteSetting(stg *model.Setting, dir string) error {
	stgFilePath := fmt.Sprintf("%s/setting_%s.csv", dir, time.Now().Format("20060102150405"))
	stgFile, err := os.Create(stgFilePath)
	if err != nil {
		return fmt.Errorf("failed to create setting file: %w", err)
	}
	defer stgFile.Close()

	stgWriter := csv.NewWriter(stgFile)
	defer stgWriter.Flush()

	if err := stgWriter.Write(stg.ToCSVHeader()); err != nil {
		return fmt.Errorf("failed to write setting header: %w", err)
	}
	if err := stgWriter.WriteAll(stg.ToCSVRows()); err != nil {
		return fmt.Errorf("failed to write setting: %w", err)
	}
	return nil
}
