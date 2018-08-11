package runtime

var (
	Reductions = 0
	Stacks     = 0
	Shares     = 0
	Datas      = 0
	Thunks     = 0
	Structs    = 0
)

var (
	stackPool  [][]Value
	sharesPool [][]*Thunk
	thunkPool  []*Thunk

	nullaryStructs [16]Struct
)

func init() {
	for i := range nullaryStructs {
		nullaryStructs[i].Index = int32(i)
	}
}

func getStack() []Value {
	if len(stackPool) == 0 {
		Stacks++
		return nil
	}
	i := len(stackPool) - 1
	stack := stackPool[i]
	stackPool = stackPool[:i]
	return stack[:0]
}

func putStack(stack []Value) {
	stackPool = append(stackPool, stack)
}

func getShares() []*Thunk {
	if len(sharesPool) == 0 {
		Shares++
		return nil
	}
	i := len(sharesPool) - 1
	shares := sharesPool[i]
	sharesPool = sharesPool[:i]
	return shares[:0]
}

func putShares(shares []*Thunk) {
	sharesPool = append(sharesPool, shares)
}

func getThunk() *Thunk {
	if len(thunkPool) == 0 {
		return &Thunk{}
	}
	i := len(thunkPool) - 1
	thunk := thunkPool[i]
	thunkPool = thunkPool[:i]
	return thunk
}

func putThunk(thunk *Thunk) {
	thunkPool = append(thunkPool, thunk)
}

func Reduce(globals []Value, value Value) (result Value) {
	var (
		stack    = getStack()
		fastData = getStack()
		shares   = getShares()
	)

beginning:
	switch v := value.(type) {
	case *Char, *Int, *Float, *Struct:
		if len(stack) > 0 {
			panic("not empty stack")
		}
		result = v
		goto end

	case *Thunk:
		if v.Result != nil {
			result = v.Result
			goto end
		}
		if v.Code == nil {
			panic("infinite reduction")
		}

		code, data := v.Code, v.Data
		if len(stack) == 0 {
			shares = append(shares, v)
			v.Code, v.Data = nil, nil
		}

		for {
			Reductions++

			switch code.Kind {
			case CodeValue:
				result = code.Value
				goto end

			case CodeOperator:
				if len(stack) != operatorArity[code.X] {
					panic("wrong number of operands on stack")
				}
				switch operatorArity[code.X] {
				case 1:
					x := stack[0]
					putStack(stack)
					putStack(fastData)
					result = operator1(globals, code.X, Reduce(globals, x))
					goto operatorEnd
				case 2:
					x, y := stack[1], stack[0]
					putStack(stack)
					putStack(fastData)
					result = operator2(globals, code.X, Reduce(globals, x), Reduce(globals, y))
					goto operatorEnd
				default:
					panic("invalid arity")
				}

			case CodeMake:
				if len(stack) == 0 && code.X < int32(len(nullaryStructs)) {
					result = &nullaryStructs[code.X]
					goto end
				}
				Structs++
				values := make([]Value, len(stack))
				copy(values, stack)
				result = &Struct{Index: code.X, Values: values}
				goto end

			case CodeVar:
				index := int32(len(data)) - code.X - 1
				value = data[index]
				goto beginning

			case CodeGlobal:
				value = globals[code.X]
				goto beginning

			case CodeAbst:
				if int32(len(stack)) < code.X {
					panic("not enough arguments on stack")
				}
				Datas++
				pop := int32(len(stack)) - code.X
				data = make([]Value, code.X)
				copy(data, stack[pop:])
				stack = stack[:pop]
				code = &code.Table[0]

			case CodeFastAbst:
				if int32(len(stack)) < code.X {
					panic("not enough arguments on stack")
				}
				pop := int32(len(stack)) - code.X
				data = append(fastData[:0], stack[pop:]...)
				stack = stack[:pop]
				code = &code.Table[0]

			case CodeAppl:
				for i := len(code.Table) - 1; i >= 1; i-- {
					switch code.Table[i].Kind {
					case CodeValue:
						stack = append(stack, code.Table[i].Value)
					case CodeVar:
						index := int32(len(data)) - code.Table[i].X - 1
						stack = append(stack, data[index])
					case CodeStrict:
						thunk := getThunk()
						thunk.Result = nil
						thunk.Code = &code.Table[i].Table[0]
						thunk.Data = data
						stack = append(stack, Reduce(globals, thunk))
						putThunk(thunk)
					default:
						Thunks++
						stack = append(stack, &Thunk{Code: &code.Table[i], Data: data})
					}
				}
				code = &code.Table[0]

			case CodeStrict:
				code = &code.Table[0]

			case CodeSwitch:
				thunk := getThunk()
				thunk.Result = nil
				thunk.Code = &code.Table[0]
				thunk.Data = data
				str := Reduce(globals, thunk).(*Struct)
				putThunk(thunk)
				stack = append(stack, str.Values...)
				code = &code.Table[str.Index+1]
			}
		}

	default:
		panic("unreachable")
	}

end:
	putStack(stack)
	putStack(fastData)
operatorEnd:
	for _, share := range shares {
		share.Result = result
	}
	putShares(shares)
	return result
}
