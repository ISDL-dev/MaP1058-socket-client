package model

import "fmt"

type Command struct {
	Name   string
	Params []string
}

const (
	// パラメータを分割するカンマの数
	NumSeparator = 9
)

func (c *Command) String() string {
	var paramsStr string
	var paramCnt int
	for _, p := range c.Params {
		paramsStr += p + ","
		paramCnt++
	}
	for i := 0; i < NumSeparator-paramCnt; i++ {
		paramsStr += ","
	}
	return fmt.Sprintf("<SCMD>%s:A:%s</SCMD>", c.Name, paramsStr)
}
