package runtime

const (
	OpIntNeg int32 = iota
	OpIntAdd
	OpIntSub
	OpIntMul
	OpIntDiv
	OpIntMod
	OpIntEq
	OpIntNeq
	OpIntLess
	OpIntLessEq
	OpIntMore
	OpIntMoreEq
)

var operatorArity = [...]int{
	OpIntNeg:    1,
	OpIntAdd:    2,
	OpIntSub:    2,
	OpIntMul:    2,
	OpIntDiv:    2,
	OpIntMod:    2,
	OpIntEq:     2,
	OpIntNeq:    2,
	OpIntLess:   2,
	OpIntLessEq: 2,
	OpIntMore:   2,
	OpIntMoreEq: 2,
}

var OperatorString = [...]string{
	OpIntNeg:    "neg",
	OpIntAdd:    "+",
	OpIntSub:    "-",
	OpIntMul:    "*",
	OpIntDiv:    "/",
	OpIntMod:    "%",
	OpIntEq:     "==",
	OpIntNeq:    "!=",
	OpIntLess:   "<",
	OpIntLessEq: "<=",
	OpIntMore:   ">",
	OpIntMoreEq: ">=",
}

func operator1(globals []Value, code int32, x Value) Value {
	switch code {
	case OpIntNeg:
		var y Int
		y.Value.Neg(&x.(*Int).Value)
		return &y
	default:
		panic("wrong operator code")
	}
}

func operator2(globals []Value, code int32, x, y Value) Value {
	switch code {
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
	case OpIntEq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) == 0 {
			return &Struct{Index: 0}
		}
		return &Struct{Index: 1}
	case OpIntNeq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) != 0 {
			return &Struct{Index: 0}
		}
		return &Struct{Index: 1}
	case OpIntLess:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) < 0 {
			return &Struct{Index: 0}
		}
		return &Struct{Index: 1}
	case OpIntLessEq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) <= 0 {
			return &Struct{Index: 0}
		}
		return &Struct{Index: 1}
	case OpIntMore:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) > 0 {
			return &Struct{Index: 0}
		}
		return &Struct{Index: 1}
	case OpIntMoreEq:
		if x.(*Int).Value.Cmp(&y.(*Int).Value) >= 0 {
			return &Struct{Index: 0}
		}
		return &Struct{Index: 1}
	default:
		panic("wrong operator code")
	}
}
