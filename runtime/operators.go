package runtime

import "math/big"

const (
	OpCharInt int32 = iota
	OpCharFloat
	OpCharInc
	OpCharDec
	OpCharAdd
	OpCharSub
	OpCharEq
	OpCharNeq
	OpCharLess
	OpCharLessEq
	OpCharMore
	OpCharMoreEq

	OpIntChar
	OpIntFloat
	OpIntNeg
	OpIntInc
	OpIntDec
	OpIntAdd
	OpIntSub
	OpIntMul
	OpIntDiv
	OpIntMod
	OpIntExp
	OpIntExpMod
	OpIntEq
	OpIntNeq
	OpIntLess
	OpIntLessEq
	OpIntMore
	OpIntMoreEq
)

var operatorArity = [...]int{
	OpCharInt:    1,
	OpCharFloat:  1,
	OpCharInc:    1,
	OpCharDec:    1,
	OpCharAdd:    2,
	OpCharSub:    2,
	OpCharEq:     2,
	OpCharNeq:    2,
	OpCharLess:   2,
	OpCharLessEq: 2,
	OpCharMore:   2,
	OpCharMoreEq: 2,

	OpIntChar:   1,
	OpIntFloat:  1,
	OpIntNeg:    1,
	OpIntInc:    1,
	OpIntDec:    1,
	OpIntAdd:    2,
	OpIntSub:    2,
	OpIntMul:    2,
	OpIntDiv:    2,
	OpIntMod:    2,
	OpIntExp:    2,
	OpIntExpMod: 3,
	OpIntEq:     2,
	OpIntNeq:    2,
	OpIntLess:   2,
	OpIntLessEq: 2,
	OpIntMore:   2,
	OpIntMoreEq: 2,
}

var OperatorString = [...]string{
	OpCharInt:    "int",
	OpCharFloat:  "float",
	OpCharInc:    "inc",
	OpCharDec:    "dec",
	OpCharAdd:    "+",
	OpCharSub:    "-",
	OpCharEq:     "==",
	OpCharNeq:    "!=",
	OpCharLess:   "<",
	OpCharLessEq: "<=",
	OpCharMore:   ">",
	OpCharMoreEq: ">=",

	OpIntChar:   "char",
	OpIntFloat:  "float",
	OpIntNeg:    "neg",
	OpIntInc:    "inc",
	OpIntDec:    "dec",
	OpIntAdd:    "+",
	OpIntSub:    "-",
	OpIntMul:    "*",
	OpIntDiv:    "/",
	OpIntMod:    "%",
	OpIntExp:    "^",
	OpIntExpMod: "^%",
	OpIntEq:     "==",
	OpIntNeq:    "!=",
	OpIntLess:   "<",
	OpIntLessEq: "<=",
	OpIntMore:   ">",
	OpIntMoreEq: ">=",
}

var bigOne = big.NewInt(1)

func operator1(globals []Value, code int32, x Value) Value {
	switch code {
	case OpCharInt:
		var y Int
		y.Value.SetInt64(int64(x.(*Char).Value))
		return &y
	case OpCharFloat:
		return &Float{Value: float64(x.(*Char).Value)}
	case OpCharInc:
		return &Char{Value: x.(*Char).Value + 1}
	case OpCharDec:
		return &Char{Value: x.(*Char).Value - 1}

	case OpIntChar:
		return &Char{Value: rune(x.(*Int).Value.Int64())}
	case OpIntFloat:
		f, _ := new(big.Float).SetInt(&x.(*Int).Value).Float64()
		return &Float{Value: f}
	case OpIntNeg:
		var y Int
		y.Value.Neg(&x.(*Int).Value)
		return &y
	case OpIntInc:
		var y Int
		y.Value.Add(&x.(*Int).Value, bigOne)
		return &y
	case OpIntDec:
		var y Int
		y.Value.Sub(&x.(*Int).Value, bigOne)
		return &y
	default:
		panic("wrong operator code")
	}
}

func operator2(globals []Value, code int32, x, y Value) Value {
	switch code {
	case OpCharAdd:
		delta := rune(y.(*Int).Value.Int64())
		return &Char{Value: x.(*Char).Value + delta}
	case OpCharSub:
		delta := rune(y.(*Int).Value.Int64())
		return &Char{Value: x.(*Char).Value - delta}
	case OpCharEq:
		if x.(*Char).Value == y.(*Char).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpCharNeq:
		if x.(*Char).Value != y.(*Char).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpCharLess:
		if x.(*Char).Value < y.(*Char).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpCharLessEq:
		if x.(*Char).Value <= y.(*Char).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpCharMore:
		if x.(*Char).Value > y.(*Char).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpCharMoreEq:
		if x.(*Char).Value >= y.(*Char).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]

	case OpIntAdd:
		var z Int
		z.Value.Add(&x.(*Int).Value, &y.(*Int).Value)
		return &z
	case OpIntSub:
		var z Int
		z.Value.Sub(&x.(*Int).Value, &y.(*Int).Value)
		return &z
	case OpIntMul:
		var z Int
		z.Value.Mul(&x.(*Int).Value, &y.(*Int).Value)
		return &z
	case OpIntDiv:
		var z Int
		z.Value.Div(&x.(*Int).Value, &y.(*Int).Value)
		return &z
	case OpIntMod:
		var z Int
		z.Value.Mod(&x.(*Int).Value, &y.(*Int).Value)
		return &z
	case OpIntExp:
		var z Int
		z.Value.Exp(&x.(*Int).Value, &y.(*Int).Value, nil)
		return &z
	case OpIntEq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) == 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpIntNeq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) != 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpIntLess:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) < 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpIntLessEq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) <= 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpIntMore:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) > 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpIntMoreEq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) >= 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	default:
		panic("wrong operator code")
	}
}

func operator3(globals []Value, code int32, x, y, z Value) Value {
	switch code {
	case OpIntExpMod:
		var w Int
		w.Value.Exp(&x.(*Int).Value, &y.(*Int).Value, &z.(*Int).Value)
		return &w
	default:
		panic("wrong operator code")
	}
}
