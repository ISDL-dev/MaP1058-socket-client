//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/ISDL-dev/MaP1058-socket-client/$GOPACKAGE
package scanner

import (
	"bufio"
	"bytes"
	"net"
)

type CustomScanner interface {
	Scan() bool
	Text() string
	Err() error
}

type customScanner struct {
	origin *bufio.Scanner
}

func NewCustomScanner(c *net.TCPConn) CustomScanner {
	s := bufio.NewScanner(c)
	onSCMD := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, []byte("</SCMD>")); i >= 0 {
			return i + 7, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return
	}
	s.Split(onSCMD)
	return &customScanner{origin: s}
}

func (s *customScanner) Scan() bool {
	return s.origin.Scan()
}

func (s *customScanner) Text() string {
	return s.origin.Text() + "</SCMD>"
}

func (s *customScanner) Err() error {
	return s.origin.Err()
}
