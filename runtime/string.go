package runtime

import (
	"fmt"
	"strings"
)

func (c *Char) String() string  { return fmt.Sprint(c.Value) }
func (i *Int) String() string   { return fmt.Sprint(&i.Value) }
func (f *Float) String() string { return fmt.Sprint(f.Value) }

func (s *Struct) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "{/%d", s.Index)
	for i := len(s.Values) - 1; i >= 0; i-- {
		b.WriteByte(' ')
		b.WriteString(s.Values[i].String())
	}
	b.WriteByte('}')
	return b.String()
}

func (t *Thunk) String() string {
	if t.Result != nil {
		return t.Result.String()
	}
	return "{...}"
}

func (c *Code) String() string {
	var b strings.Builder
	indented(0, c, &b)
	return b.String()
}

func indented(level int, c *Code, b *strings.Builder) {
	for i := 0; i < level*2; i++ {
		b.WriteByte(' ')
	}
	fmt.Fprintf(b, "%s  %d  %v\n", codeNames[c.Kind], c.X, c.Value)
	for i := range c.Table {
		indented(level+1, &c.Table[i], b)
	}
}

var codeNames = [...]string{
	CodeValue:    "VALUE",
	CodeOperator: "OPERATOR",
	CodeMake:     "MAKE",
	CodeVar:      "VAR",
	CodeGlobal:   "GLOBAL",
	CodeAbst:     "ABST",
	CodeFastAbst: "FASTABST",
	CodeAppl:     "APPL",
	CodeSwitch:   "SWITCH",
}
