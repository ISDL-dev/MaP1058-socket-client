package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCommand(t *testing.T) {
	t.Run("stringからCommandに変換", func(t *testing.T) {
		rCmdStr := "<SCMD>START:A:300,2023/01/01 12-00-00,,,,,,,,</SCMD>"
		parser := NewParser()
		cmd, err := parser.ToCommand(rCmdStr)
		assert.NoError(t, err)
		assert.Equal(t, "START", cmd.Name)
		expectedParams := [10]string{}
		expectedParams[0] = "300"
		expectedParams[1] = "2023/01/01 12-00-00"
		assert.Equal(t, expectedParams, cmd.Params)
	})

	t.Run("<SCMD>...</SCMD>の構造になっていない", func(t *testing.T) {
		rCmdStr := "<S>START:A:300,2023/01/01 12-00-00,,,,,,,,</S>"
		parser := NewParser()
		_, err := parser.ToCommand(rCmdStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "<SCMD> and </SCMD> on both sides")
	})

	t.Run("規定の数のカンマがない", func(t *testing.T) {
		rCmdStr := "<SCMD>START:A:300,2023/01/01 12-00-00,,,,</SCMD>"
		parser := NewParser()
		_, err := parser.ToCommand(rCmdStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "9 commas")
	})

	t.Run("コマンド名とパラメータを分ける:A:が含まれていない", func(t *testing.T) {
		rCmdStr := "<SCMD>START:300,2023/01/01 12-00-00,,,,,,,,</SCMD>"
		parser := NewParser()
		_, err := parser.ToCommand(rCmdStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ":A:")
	})
}
