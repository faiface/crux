package crux

import (
	"fmt"
	"strings"

	"github.com/faiface/crux/runtime"
)

func (c *Char) String() string  { return fmt.Sprint(c.Value) }
func (i *Int) String() string   { return fmt.Sprint(&i.Value) }
func (f *Float) String() string { return fmt.Sprint(f.Value) }

func (o *Operator) String() string { return runtime.OperatorString[o.Code] }
func (m *Make) String() string     { return fmt.Sprintf("#make/%d", m.Index) }

func (v *Var) String() string {
	if v.Index < 0 {
		return v.Name
	}
	return fmt.Sprintf("%s/%d", v.Name, v.Index)
}

func (a *Abst) String() string {
	var b strings.Builder
	b.WriteString("(\\")
	for _, bound := range a.Bound {
		b.WriteString(bound)
		b.WriteByte(' ')
	}
	b.WriteString("-> ")
	b.WriteString(a.Body.String())
	b.WriteByte(')')
	return b.String()
}

func (a *FastAbst) String() string {
	var b strings.Builder
	b.WriteString("(\\")
	for _, bound := range a.Bound {
		b.WriteString(bound)
		b.WriteByte(' ')
	}
	b.WriteString("=> ")
	b.WriteString(a.Body.String())
	b.WriteByte(')')
	return b.String()
}

func (a *Appl) String() string {
	var b strings.Builder
	b.WriteByte('(')
	b.WriteString(a.Rator.String())
	for _, rand := range a.Rands {
		b.WriteByte(' ')
		b.WriteString(rand.String())
	}
	b.WriteByte(')')
	return b.String()
}

func (s *Switch) String() string {
	var b strings.Builder
	b.WriteString("(#switch ")
	b.WriteString(s.Expr.String())
	for _, cas := range s.Cases {
		b.WriteByte(' ')
		b.WriteString(cas.String())
	}
	b.WriteByte(')')
	return b.String()
}
