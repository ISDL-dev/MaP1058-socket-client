//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/socket-capture-signals/$GOPACKAGE
package parser

import (
	"fmt"

	"github.com/Be3751/socket-capture-signals/internal/model"
)

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

// AD値のバイト列を解析してAD値を持つmodel.Signals型のポインタを返す
func (p *parser) ToSignals(adSignals []byte) (*model.Signals, error) {
	if len(adSignals) != int(p.Config.SumBytes) {
		return nil, fmt.Errorf("adSignals' len must be %d", p.Config.SumBytes)
	}
	var result *model.Signals
	var sum uint16
	var chCnt int
	var pntCnt int
	for i := 0; i < int(p.Config.SumBytes); i += 2 {
		// 前半8チャネルのみを使用
		if chCnt != 0 && chCnt%8 == 0 {
			chCnt = 0
		}
		// 16ループに1回、ポイント数をインクリメント
		if i != 0 && i%32 == 0 {
			pntCnt++
		}
		// 未使用のチャネルは無視
		if adSignals[i] == 0 {
			continue
		}
		formerByte := uint16(adSignals[i])
		latterByte := uint16(adSignals[i+1])
		adValue := (formerByte << 8) | latterByte // 8ビットシフトとビット和で2byteのAD値を構成
		result.Channels[chCnt].ADValues[pntCnt] = adValue
		// 先頭10ポイント分の合計でサムチェック
		if pntCnt < 10 {
			sum += adValue
		}
		chCnt++
	}
	bytesSumCheckCode := adSignals[p.Config.SumBytes-p.Config.SumCheckCodeSize:]
	formerByte := uint16(bytesSumCheckCode[0])
	latterByte := uint16(bytesSumCheckCode[1])
	valueSumCheckCode := (formerByte << 8) | latterByte
	if valueSumCheckCode != sum {
		return nil, &FailureSumCheckError{Expected: valueSumCheckCode, Actual: sum}
	}
	return result, nil
}

type FailureSumCheckError struct {
	Expected uint16
	Actual   uint16
}

func (e *FailureSumCheckError) Error() string {
	return fmt.Sprintf("parsed invalid signals doesn't match with the sum-check-code: expected %d but actual %d", e.Expected, e.Actual)
}
