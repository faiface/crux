package runtime

var Reductions = 0

func Reduce(globals []Value, value Value, stack ...Value) Value {
	switch v := value.(type) {
	case *Char, *Int, *Float, *Struct:
		if len(stack) > 0 {
			panic("not empty stack")
		}
		return v

	case *Thunk:
		if v.Result != nil {
			return v.Result
		}
		if v.Code == nil && v.Data == nil {
			panic("infinite reduction")
		}

		code, data := v.Code, v.Data

		result := Value(nil)
		share := false
		if len(stack) == 0 {
			v.Code, v.Data = nil, nil
			share = true
		}

	loop:
		for {
			Reductions++

			switch code.Kind {
			case CodeValue:
				result = code.Value
				break loop

			case CodeOperator:
				switch operatorArity[code.X] {
				case 1:
					result = operator1(globals, code.X, stack[len(stack)-1])
					break loop
				case 2:
					result = operator2(globals, code.X, stack[len(stack)-1], stack[len(stack)-2])
					break loop
				}

			case CodeMake:
				result = &Struct{Index: code.X, Values: stack}
				break loop

			case CodeVar:
				index := int32(len(data)) - code.X - 1
				result = Reduce(globals, data[index], stack...)
				break loop

			case CodeGlobal:
				result = Reduce(globals, globals[code.X], stack...)
				break loop

			case CodeAbst:
				if int32(len(stack)) < code.X {
					panic("not enough arguments on the stack")
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

		if share {
			v.Result = result
		}
		return result
	}
	panic("unreachable")
}
