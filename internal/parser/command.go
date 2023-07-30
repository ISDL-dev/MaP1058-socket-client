package parser

import (
	"fmt"
	"strings"

	"github.com/Be3751/MaP1058-socket-client/internal/model"
)

func (p *parser) ToCommand(s string) (*model.Command, error) {
	baseErrMsg := "the received variable is an unexpected pattern"
	if !strings.HasPrefix(s, "<SCMD>") || !strings.HasSuffix(s, "</SCMD>") {
		return nil, fmt.Errorf("%s, s string must contain <SCMD> and </SCMD> on both sides: %s", baseErrMsg, s)
	} else if !strings.Contains(s, ":A:") {
		return nil, fmt.Errorf("%s, s string must contain \":A:\": %s", baseErrMsg, s)
	} else if strings.Count(s, ",") != model.NumSeparator {
		return nil, fmt.Errorf("%s, s string must contain %d commas: %s", baseErrMsg, model.NumSeparator, s)
	}

	s = strings.TrimPrefix(s, "<SCMD>")
	s = strings.TrimSuffix(s, "</SCMD>")
	nameAndParams := strings.Split(s, ":A:")
	name := nameAndParams[0]
	paramsStr := nameAndParams[1]
	sliceParams := strings.Split(paramsStr, ",")
	var params [model.NumSeparator + 1]string
	for i, p := range sliceParams {
		params[i] = p
	}
	return &model.Command{Name: name, Params: params}, nil
}
