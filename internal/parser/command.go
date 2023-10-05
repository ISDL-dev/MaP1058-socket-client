package parser

import (
	"fmt"
	"regexp"
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
	var params [model.NumSeparator + 1]string

	// GETSETTINGコマンドのみ例外的なフォーマット
	if name == "GETSETTING" {
		var calStr string
		if _, err := fmt.Sscanf(paramsStr, "\"%s\"", &calStr); err != nil {
			return nil, fmt.Errorf("%s, failed to scan paramsStr as calibration string: %w", baseErrMsg, err)
		}
		copy(params[:], strings.Split(calStr, ",")[:5])
		return &model.Command{Name: name, Params: params}, nil
	}

	copy(params[:], strings.Split(paramsStr, ","))
	return &model.Command{Name: name, Params: params}, nil
}

func (p *parser) ToTrendRange(c *model.Command) (model.TrendRange, error) {
	var tr model.TrendRange
	if c.Name != "RANGE" {
		return tr, fmt.Errorf("the received command is not RANGE: %s", c.Name)
	} else if c.NumValueParams() != 8 {
		return tr, fmt.Errorf(`the received command has %d with-value parameters, 
			but it should have 8 with-value parameters: %s`, c.NumValueParams(), c.String())
	}
	for i, p := range c.Params {
		if p == "" {
			break
		}
		var cr model.ChannelRange
		if _, err := fmt.Sscanf(p, "%d;%d", &cr.Upper, &cr.Lower); err != nil {
			return tr, fmt.Errorf("failed to scan %dth parameter as ChannelRange: %s", i, err)
		}
		tr[i] = cr
	}
	return tr, nil
}

func (p *parser) ToAnalysis(c *model.Command) (model.Analysis, error) {
	var a model.Analysis
	if c.Name != "ANALYSIS" {
		return a, fmt.Errorf("the received command is not ANALYSIS: %s", c.Name)
	} else if c.NumValueParams() != 8 {
		return a, fmt.Errorf(`the received command has %d with-value parameters, 
			but it should have 8 with-value parameters: %s`, c.NumValueParams(), c.String())
	}
	for i, p := range c.Params {
		if p == "" {
			break
		}
		var ca model.ChannelAnalysis
		if _, err := fmt.Sscanf(p, "%d", &ca); err != nil {
			return a, fmt.Errorf("failed to scan %dth parameter as ChannelAnalysis: %s", i, err)
		}
		a[i] = ca
	}
	return a, nil
}

func (p *parser) ToChannelCal(c *model.Command) (model.ChannelCal, error) {
	var cal model.ChannelCal
	if c.Name != "GETSETTING" {
		return cal, fmt.Errorf("the received command is not GETSETTING: %s", c.Name)
	} else if c.NumValueParams() != 5 {
		return cal, fmt.Errorf(`the received command has %d with-value parameters, 
			but it should have 5 with-value parameters: %s`, c.NumValueParams(), c.String())
	}
	for i, p := range c.Params {
		if p == "" {
			break
		}
		if _, err := fmt.Sscanf(p, "%f,%f,%f,%f", &cal.BaseAD, &cal.CalAD, &cal.EuHi, &cal.EuLo); err != nil {
			return cal, fmt.Errorf("failed to scan %dth parameter as ChannelCal: %s", i, err)
		}
	}
	return cal, nil
}

func FindSCMD(s string) string {
	re := regexp.MustCompile(`<SCMD>[a-zA-Z]+:A:[a-zA-Z]+</SCMD>`)
	return re.FindString(s)
}
