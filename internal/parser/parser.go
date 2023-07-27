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
	Signal Signal
}

type Signal struct {
	// 受信信号のバイト数
	SumBytes uint64
	// 受信信号に含まれるサムチェックコードのバイト数
	SumCheckCodeSize uint64
	// 受信信号におけるポイント数
	NumPoints uint64
	// 受信信号における総チャンネル数
	NumChannels uint64
	// 受信信号における有効チャンネルのインデックス
	IndexAvailableChs []int
	// サムチェックに使用するポイントのインデックス
	IndexPntsSumCheck []int
}
