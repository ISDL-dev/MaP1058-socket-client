//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE -self_package=github.com/Be3751/MaP1058-socket-client/$GOPACKAGE
package parser

import "github.com/Be3751/MaP1058-socket-client/internal/model"

type Parser interface {
	ToSignals(adSignals []byte) (*model.Signals, error)
}

func NewParser() Parser {
	return &parser{}
}

type parser struct {
}
