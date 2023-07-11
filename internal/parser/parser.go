//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import "github.com/Be3751/MaP1058-socket-client/internal/model"

type Parser interface {
	ToSignals(adSignals []byte) (*model.Signals, error)
}

func NewParser(c ParseConfig) Parser {
	return &parser{
		Config: c,
	}
}

type parser struct {
	// 解析に必要な設定値
	Config ParseConfig
}

type ParseConfig struct {
	// 1受信あたりに得られる信号のバイト数
	SumBytes uint64
	// 1受信あたりに得られる信号に含まれるサムチェックコードのバイト数
	SumCheckCodeSize uint64
}
