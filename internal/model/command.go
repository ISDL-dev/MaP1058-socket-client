package model

import "fmt"

type Command struct {
	Name   string
	Params []string
}

func (c *Command) String() string {
	var paramsStr string
	var paramCnt int
	for _, p := range c.Params {
		paramsStr += p + ","
		paramCnt++
	}
	for i := 0; i < (10-paramCnt)-1; i++ {
		paramsStr += ","
	}
	return fmt.Sprintf("<SCMD>%s:A:%s</SCMD>", c.Name, paramsStr)
}
