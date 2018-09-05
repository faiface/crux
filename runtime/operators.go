package runtime

import (
	"fmt"
	"math"
	"math/big"
)

const (
	OpCharInt int32 = iota
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
	OpIntString
	OpIntNeg
	OpIntAbs
	OpIntInc
	OpIntDec
	OpIntAdd
	OpIntSub
	OpIntMul
	OpIntDiv
	OpIntMod
	OpIntExp
	OpIntEq
	OpIntNeq
	OpIntLess
	OpIntLessEq
	OpIntMore
	OpIntMoreEq

	OpFloatInt
	OpFloatString
	OpFloatNeg
	OpFloatAbs
	OpFloatInc
	OpFloatDec
	OpFloatAdd
	OpFloatSub
	OpFloatMul
	OpFloatDiv
	OpFloatMod
	OpFloatExp
	OpFloatEq
	OpFloatNeq
	OpFloatLess
	OpFloatLessEq
	OpFloatMore
	OpFloatMoreEq
)

var operatorArity = [...]int{
	OpCharInt:    1,
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
	OpIntString: 1,
	OpIntNeg:    1,
	OpIntAbs:    1,
	OpIntInc:    1,
	OpIntDec:    1,
	OpIntAdd:    2,
	OpIntSub:    2,
	OpIntMul:    2,
	OpIntDiv:    2,
	OpIntMod:    2,
	OpIntExp:    2,
	OpIntEq:     2,
	OpIntNeq:    2,
	OpIntLess:   2,
	OpIntLessEq: 2,
	OpIntMore:   2,
	OpIntMoreEq: 2,

	OpFloatInt:    1,
	OpFloatString: 1,
	OpFloatNeg:    1,
	OpFloatAbs:    1,
	OpFloatInc:    1,
	OpFloatDec:    1,
	OpFloatAdd:    2,
	OpFloatSub:    2,
	OpFloatMul:    2,
	OpFloatDiv:    2,
	OpFloatMod:    2,
	OpFloatExp:    2,
	OpFloatEq:     2,
	OpFloatNeq:    2,
	OpFloatLess:   2,
	OpFloatLessEq: 2,
	OpFloatMore:   2,
	OpFloatMoreEq: 2,
}

var OperatorString = [...]string{
	OpCharInt:    "int",
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
	OpIntString: "string",
	OpIntNeg:    "neg",
	OpIntAbs:    "abs",
	OpIntInc:    "inc",
	OpIntDec:    "dec",
	OpIntAdd:    "+",
	OpIntSub:    "-",
	OpIntMul:    "*",
	OpIntDiv:    "/",
	OpIntMod:    "%",
	OpIntExp:    "^",
	OpIntEq:     "==",
	OpIntNeq:    "!=",
	OpIntLess:   "<",
	OpIntLessEq: "<=",
	OpIntMore:   ">",
	OpIntMoreEq: ">=",

	OpFloatInt:    "int",
	OpFloatString: "string",
	OpFloatNeg:    "neg",
	OpFloatAbs:    "abs",
	OpFloatInc:    "inc",
	OpFloatDec:    "dec",
	OpFloatAdd:    "+",
	OpFloatSub:    "-",
	OpFloatMul:    "*",
	OpFloatDiv:    "/",
	OpFloatMod:    "%",
	OpFloatExp:    "^",
	OpFloatEq:     "==",
	OpFloatNeq:    "!=",
	OpFloatLess:   "<",
	OpFloatLessEq: "<=",
	OpFloatMore:   ">",
	OpFloatMoreEq: ">=",
}

var bigOne = big.NewInt(1)

func operator1(code int32, x Value) Value {
	switch code {
	case OpCharInt:
		var y Int
		y.Value.SetInt64(int64(x.(*Char).Value))
		return &y
	case OpCharInc:
		return &Char{Value: x.(*Char).Value + 1}
	case OpCharDec:
		return &Char{Value: x.(*Char).Value - 1}

	case OpIntChar:
		return &Char{Value: rune(x.(*Int).Value.Int64())}
	case OpIntFloat:
		f, _ := new(big.Float).SetInt(&x.(*Int).Value).Float64()
		return &Float{Value: f}
	case OpIntString:
		runes := []rune(x.(*Int).Value.Text(10))
		chars := make([]Char, len(runes))
		for i := range chars {
			chars[i].Value = runes[i]
		}
		str := &Struct{Index: 0}
		for i := len(runes) - 1; i >= 0; i-- {
			str = &Struct{Index: 1, Values: []Value{str, &chars[i]}}
		}
		return str
	case OpIntNeg:
		var y Int
		y.Value.Neg(&x.(*Int).Value)
		return &y
	case OpIntAbs:
		var y Int
		y.Value.Abs(&x.(*Int).Value)
		return &y
	case OpIntInc:
		var y Int
		y.Value.Add(&x.(*Int).Value, bigOne)
		return &y
	case OpIntDec:
		var y Int
		y.Value.Sub(&x.(*Int).Value, bigOne)
		return &y

	case OpFloatInt:
		var y Int
		big.NewFloat(math.Floor(x.(*Float).Value)).Int(&y.Value)
		return &y
	case OpFloatString:
		runes := []rune(fmt.Sprint(x.(*Float).Value))
		chars := make([]Char, len(runes))
		for i := range chars {
			chars[i].Value = runes[i]
		}
		str := &Struct{Index: 0}
		for i := len(runes) - 1; i >= 0; i-- {
			str = &Struct{Index: 1, Values: []Value{str, &chars[i]}}
		}
		return str
	case OpFloatNeg:
		return &Float{Value: -x.(*Float).Value}
	case OpFloatAbs:
		return &Float{Value: math.Abs(x.(*Float).Value)}
	case OpFloatInc:
		return &Float{Value: x.(*Float).Value + 1}
	case OpFloatDec:
		return &Float{Value: x.(*Float).Value - 1}

	default:
		panic("wrong operator code")
	}
}

func operator2(code int32, x, y Value) Value {
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

	case OpFloatAdd:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: xf + yf}
	case OpFloatSub:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: xf - yf}
	case OpFloatMul:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: xf * yf}
	case OpFloatDiv:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: xf / yf}
	case OpFloatMod:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: xf - yf*math.Floor(xf/yf)}
	case OpFloatExp:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: math.Pow(xf, yf)}
	case OpFloatEq:
		if x.(*Float).Value == y.(*Float).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatNeq:
		if x.(*Float).Value != y.(*Float).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatLess:
		if x.(*Float).Value < y.(*Float).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatLessEq:
		if x.(*Float).Value <= y.(*Float).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatMore:
		if x.(*Float).Value > y.(*Float).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatMoreEq:
		if x.(*Float).Value >= y.(*Float).Value {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]

	default:
		panic("wrong operator code")
	}
}
