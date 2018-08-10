package runtime

var Reductions = 0

var (
	stackPool  [][]Value
	sharesPool [][]*Thunk
)

func getStack() []Value {
	if len(stackPool) == 0 {
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

func Reduce(globals []Value, value Value) (result Value) {
	var (
		stack  = getStack()
		shares = getShares()
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
				switch operatorArity[code.X] {
				case 1:
					result = operator1(globals, code.X, stack[len(stack)-1])
					goto end
				case 2:
					result = operator2(globals, code.X, stack[len(stack)-1], stack[len(stack)-2])
					goto end
				default:
					panic("invalid arity")
				}

			case CodeMake:
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
				pop := int32(len(stack)) - code.X
				data = make([]Value, code.X)
				copy(data, stack[pop:])
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
					default:
						stack = append(stack, &Thunk{Code: &code.Table[i], Data: data})
					}
				}
				code = &code.Table[0]

			case CodeSwitch:
				str := Reduce(globals, &Thunk{Code: &code.Table[0], Data: data}).(*Struct)
				stack = append(stack, str.Values...)
				code = &code.Table[str.Index+1]
			}
		}

	default:
		panic("unreachable")
	}

end:
	for _, share := range shares {
		share.Result = result
	}
	putStack(stack)
	putShares(shares)
	return result
}
