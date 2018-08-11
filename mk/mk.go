package mk

import (
	"math/big"

	"github.com/faiface/crux"
	"github.com/faiface/crux/runtime"
)

var operators map[string]int32

func init() {
	operators = make(map[string]int32)
	for code, str := range runtime.OperatorString {
		operators[str] = int32(code)
	}
}

func Char(c rune) *crux.Char {
	return &crux.Char{Value: c}
}

func Int(i int) *crux.Int {
	var bi big.Int
	bi.SetInt64(int64(i))
	return &crux.Int{Value: bi}
}

func BigInt(i *big.Int) *crux.Int {
	var bi big.Int
	bi.Set(i)
	return &crux.Int{Value: bi}
}

func Float(f float64) *crux.Float {
	return &crux.Float{Value: f}
}

func Operator(code int32) *crux.Operator {
	return &crux.Operator{Code: code}
}

func Op(s string) *crux.Operator {
	return &crux.Operator{Code: operators[s]}
}

func Make(index int32) *crux.Make {
	return &crux.Make{Index: index}
}

func Var(name string, index int32) *crux.Var {
	return &crux.Var{Name: name, Index: index}
}

func Abst(bound ...string) func(body crux.Expr) *crux.Abst {
	return func(body crux.Expr) *crux.Abst {
		return &crux.Abst{Bound: bound, Body: body}
	}
}

func Appl(rator crux.Expr, rands ...crux.Expr) *crux.Appl {
	return &crux.Appl{Rator: rator, Rands: rands}
}

func Strict(expr crux.Expr) *crux.Strict {
	return &crux.Strict{Expr: expr}
}

func Switch(expr crux.Expr, cases ...crux.Expr) *crux.Switch {
	return &crux.Switch{Expr: expr, Cases: cases}
}
