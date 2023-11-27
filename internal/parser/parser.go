//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import "github.com/Be3751/MaP1058-socket-client/internal/model"

type Parser interface {
	// AD値のバイト列を解析してAD値を持つmodel.Signals型のポインタを返す
	ToSignals(b []byte) (*model.Signals, error)
	// <SCMD>...</SCMD>の文字列を解析してmodel.Commandを返す
	ToCommand(s string) (model.Command, error)
	// model.Commandを解析してTrendRangeを返す
	ToTrendRange(c model.Command) (model.TrendRange, error)
	// model.Commandを解析してAnalysisTypeを返す
	ToAnalysis(c model.Command) (model.AnalysisType, error)
	// model.Commandを解析してCalibrationを返す
	ToChannelCal(c model.Command) (model.ChannelCal, error)
	// model.Commandを解析してEEGPowerを返す
	ToChannelPower(c model.Command) (model.ChannelPower, error)
}

func NewParser() Parser {
	return &parser{}
}

type parser struct {
}
