package model

import "fmt"

type Command struct {
	Name   string
	Params [NumSeparator + 1]string
}

const (
	// パラメータを分割するカンマの数
	NumSeparator = 9
)

func (c *Command) String() string {
	var paramsStr string
	for i, p := range c.Params {
		if i == NumSeparator {
			paramsStr += p
		} else {
			paramsStr += p + ","
		}
	}
	return fmt.Sprintf("<SCMD>%s:A:%s</SCMD>", c.Name, paramsStr)
}

// c.Paramsのうち空の文字列でない要素の数を返す
func (c *Command) NumValueParams() int {
	var num int
	for _, p := range c.Params {
		if p != "" {
			num++
		}
	}
	return num
}
