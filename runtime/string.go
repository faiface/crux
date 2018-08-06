package runtime

import (
	"fmt"
	"strings"
)

func (c *Code) String() string {
	var b strings.Builder
	indented(0, c, &b)
	return b.String()
}

func indented(level int, c *Code, b *strings.Builder) {
	for i := 0; i < level*2; i++ {
		b.WriteByte(' ')
	}
	fmt.Fprintf(b, "%s %6d %v", codeNames[c.Kind], c.X, c.Value)
	for i := range c.Table {
		indented(level+1, &c.Table[i], b)
	}
}

var codeNames = [...]string{
	CodeValue:    "VALUE   ",
	CodeOperator: "OPERATOR",
	CodeMake:     "MAKE    ",
	CodeVar:      "VAR     ",
	CodeGlobal:   "GLOBAL  ",
	CodeAbst:     "ABST    ",
	CodeAppl:     "APPL    ",
	CodeSwitch:   "SWITCH  ",
}
