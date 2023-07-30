package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewString(t *testing.T) {
	t.Run("START", func(t *testing.T) {
		cmd := Command{
			Name: "START",
			Params: [10]string{
				"300",
				"2023/01/01 12-00-00",
			},
		}
		cmdStr := cmd.String()
		assert.Equal(t, "<SCMD>START:A:300,2023/01/01 12-00-00,,,,,,,,</SCMD>", cmdStr)
	})

	t.Run("END", func(t *testing.T) {
		cmd := Command{
			Name: "END",
		}
		cmdStr := cmd.String()
		assert.Equal(t, "<SCMD>END:A:,,,,,,,,,</SCMD>", cmdStr)
	})
}
