package runtime

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
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
	OpIntIsZero

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
	OpFloatIsPlusInf
	OpFloatIsMinusInf
	OpFloatIsInf
	OpFloatIsNan
	OpFloatSin
	OpFloatCos
	OpFloatTan
	OpFloatAsin
	OpFloatAcos
	OpFloatAtan
	OpFloatAtan2
	OpFloatSinh
	OpFloatCosh
	OpFloatTanh
	OpFloatAsinh
	OpFloatAcosh
	OpFloatAtanh
	OpFloatCeil
	OpFloatFloor
	OpFloatSqrt
	OpFloatCbrt
	OpFloatLog
	OpFloatHypot
	OpFloatGamma

	OpStringInt
	OpStringFloat

	OpError
	OpDump
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
	OpIntIsZero: 1,

	OpFloatInt:        1,
	OpFloatString:     1,
	OpFloatNeg:        1,
	OpFloatAbs:        1,
	OpFloatInc:        1,
	OpFloatDec:        1,
	OpFloatAdd:        2,
	OpFloatSub:        2,
	OpFloatMul:        2,
	OpFloatDiv:        2,
	OpFloatMod:        2,
	OpFloatExp:        2,
	OpFloatEq:         2,
	OpFloatNeq:        2,
	OpFloatLess:       2,
	OpFloatLessEq:     2,
	OpFloatMore:       2,
	OpFloatMoreEq:     2,
	OpFloatIsPlusInf:  1,
	OpFloatIsMinusInf: 1,
	OpFloatIsInf:      1,
	OpFloatIsNan:      1,
	OpFloatSin:        1,
	OpFloatCos:        1,
	OpFloatTan:        1,
	OpFloatAsin:       1,
	OpFloatAcos:       1,
	OpFloatAtan:       1,
	OpFloatAtan2:      2,
	OpFloatSinh:       1,
	OpFloatCosh:       1,
	OpFloatTanh:       1,
	OpFloatAsinh:      1,
	OpFloatAcosh:      1,
	OpFloatAtanh:      1,
	OpFloatCeil:       1,
	OpFloatFloor:      1,
	OpFloatSqrt:       1,
	OpFloatCbrt:       1,
	OpFloatLog:        1,
	OpFloatHypot:      2,
	OpFloatGamma:      1,

	OpStringInt:   1,
	OpStringFloat: 1,

	OpError: 1,
	OpDump:  2,
}

var OperatorString = [...]string{
	OpCharInt:    "char->int",
	OpCharInc:    "inc/char",
	OpCharDec:    "dec/char",
	OpCharAdd:    "+/char",
	OpCharSub:    "-/char",
	OpCharEq:     "==/char",
	OpCharNeq:    "!=/char",
	OpCharLess:   "</char",
	OpCharLessEq: "<=/char",
	OpCharMore:   ">/char",
	OpCharMoreEq: ">=/char",

	OpIntChar:   "int->char",
	OpIntFloat:  "int->float",
	OpIntString: "int->string",
	OpIntNeg:    "neg/int",
	OpIntAbs:    "abs/int",
	OpIntInc:    "inc/int",
	OpIntDec:    "dec/int",
	OpIntAdd:    "+/int",
	OpIntSub:    "-/int",
	OpIntMul:    "*/int",
	OpIntDiv:    "//int",
	OpIntMod:    "%/int",
	OpIntExp:    "^/int",
	OpIntEq:     "==/int",
	OpIntNeq:    "!=/int",
	OpIntLess:   "</int",
	OpIntLessEq: "<=/int",
	OpIntMore:   ">/int",
	OpIntMoreEq: ">=/int",
	OpIntIsZero: "zero?/int",

	OpFloatInt:        "float->int",
	OpFloatString:     "float->string",
	OpFloatNeg:        "neg/float",
	OpFloatAbs:        "abs/float",
	OpFloatInc:        "inc/float",
	OpFloatDec:        "dec/float",
	OpFloatAdd:        "+/float",
	OpFloatSub:        "-/float",
	OpFloatMul:        "*/float",
	OpFloatDiv:        "//float",
	OpFloatMod:        "%/float",
	OpFloatExp:        "^/float",
	OpFloatEq:         "==/float",
	OpFloatNeq:        "!=/float",
	OpFloatLess:       "</float",
	OpFloatLessEq:     "<=/float",
	OpFloatMore:       ">/float",
	OpFloatMoreEq:     ">=/float",
	OpFloatIsPlusInf:  "+inf?",
	OpFloatIsMinusInf: "-inf?",
	OpFloatIsInf:      "inf?",
	OpFloatIsNan:      "nan?",
	OpFloatSin:        "sin",
	OpFloatCos:        "cos",
	OpFloatTan:        "tan",
	OpFloatAsin:       "asin",
	OpFloatAcos:       "acos",
	OpFloatAtan:       "atan",
	OpFloatAtan2:      "atan2",
	OpFloatSinh:       "sinh",
	OpFloatCosh:       "cosh",
	OpFloatTanh:       "tanh",
	OpFloatAsinh:      "asinh",
	OpFloatAcosh:      "acosh",
	OpFloatAtanh:      "atanh",
	OpFloatCeil:       "ceil",
	OpFloatFloor:      "floor",
	OpFloatSqrt:       "sqrt",
	OpFloatCbrt:       "cbrt",
	OpFloatLog:        "log",
	OpFloatHypot:      "hypot",
	OpFloatGamma:      "gamma",

	OpStringInt:   "string->int",
	OpStringFloat: "string->float",

	OpError: "error",
	OpDump:  "dump",
}

var bigOne = big.NewInt(1)

func accumString(globals []Value, x Value) string {
	var b strings.Builder
	for str := x.(*Struct); str.Index != 0; str = Reduce(globals, str.Values[0]).(*Struct) {
		b.WriteRune(Reduce(globals, str.Values[1]).(*Char).Value)
	}
	return b.String()
}

func operator1(globals []Value, code int32, x Value) Value {
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
	case OpIntIsZero:
		if x.(*Int).Value.Sign() == 0 {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]

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
	case OpFloatIsMinusInf:
		if math.IsInf(x.(*Float).Value, -1) {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatIsPlusInf:
		if math.IsInf(x.(*Float).Value, +1) {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatIsInf:
		if math.IsInf(x.(*Float).Value, 0) {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatIsNan:
		if math.IsNaN(x.(*Float).Value) {
			return &nullaryStructs[0]
		}
		return &nullaryStructs[1]
	case OpFloatSin:
		return &Float{Value: math.Sin(x.(*Float).Value)}
	case OpFloatCos:
		return &Float{Value: math.Cos(x.(*Float).Value)}
	case OpFloatTan:
		return &Float{Value: math.Tan(x.(*Float).Value)}
	case OpFloatAsin:
		return &Float{Value: math.Asin(x.(*Float).Value)}
	case OpFloatAcos:
		return &Float{Value: math.Acos(x.(*Float).Value)}
	case OpFloatAtan:
		return &Float{Value: math.Atan(x.(*Float).Value)}
	case OpFloatSinh:
		return &Float{Value: math.Sinh(x.(*Float).Value)}
	case OpFloatCosh:
		return &Float{Value: math.Cosh(x.(*Float).Value)}
	case OpFloatTanh:
		return &Float{Value: math.Tanh(x.(*Float).Value)}
	case OpFloatAsinh:
		return &Float{Value: math.Asinh(x.(*Float).Value)}
	case OpFloatAcosh:
		return &Float{Value: math.Acosh(x.(*Float).Value)}
	case OpFloatAtanh:
		return &Float{Value: math.Atanh(x.(*Float).Value)}
	case OpFloatCeil:
		return &Float{Value: math.Ceil(x.(*Float).Value)}
	case OpFloatFloor:
		return &Float{Value: math.Floor(x.(*Float).Value)}
	case OpFloatSqrt:
		return &Float{Value: math.Sqrt(x.(*Float).Value)}
	case OpFloatCbrt:
		return &Float{Value: math.Cbrt(x.(*Float).Value)}
	case OpFloatLog:
		return &Float{Value: math.Log(x.(*Float).Value)}
	case OpFloatGamma:
		return &Float{Value: math.Gamma(x.(*Float).Value)}

	case OpStringInt:
		var i Int
		fmt.Sscanf("%d", accumString(globals, x), &i.Value)
		return &i
	case OpStringFloat:
		var f Float
		fmt.Sscan(accumString(globals, x), &f.Value)
		return &f

	case OpError:
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", accumString(globals, x))
		os.Exit(1)
		return nil

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
	case OpFloatAtan2:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: math.Atan2(xf, yf)}
	case OpFloatHypot:
		xf, yf := x.(*Float).Value, y.(*Float).Value
		return &Float{Value: math.Hypot(xf, yf)}

	case OpDump:
		fmt.Fprintln(os.Stderr, accumString(globals, x))
		return y

	default:
		panic("wrong operator code")
	}
}
