//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import "github.com/Be3751/MaP1058-socket-client/internal/model"

type Parser interface {
	// AD値のバイト列を解析してAD値を持つmodel.Signals型のポインタを返す
	ToSignals(b []byte, s *model.Signals) error
}

func NewParser() Parser {
	return &parser{}
}

type parser struct {
}
