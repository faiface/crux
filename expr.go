package crux

import "math/big"

type Expr interface {
	String() string
}

type (
	Char  struct{ Value rune }
	Int   struct{ Value big.Int }
	Float struct{ Value float64 }

	Operator struct{ Code int32 }
	Make     struct{ Index int32 }

	Var struct {
		Name  string
		Index int
	}

	Abst struct {
		Bound []string
		Body  Expr
	}

	Appl struct {
		Rator Expr
		Rands []Expr
	}

	Switch struct {
		Expr  Expr
		Cases []Expr
	}
)
