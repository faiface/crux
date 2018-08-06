package runtime

func Reduce(globals []Value, value Value) Value {
beginning:
	switch v := value.(type) {
	case *Char, *Int, *Float, *Struct:
		return v

	case *Thunk:
		if v.Result != nil {
			return v.Result
		}
		if v.Code == nil && v.Data == nil {
			panic("infinite reduction")
		}

		code, data := v.Code, v.Data
		v.Code, v.Data = nil, nil

		stack := []Value(nil)
		result := Value(nil)

	loop:
		for {
			switch code.Kind {
			case CodeValue:
				result = code.Value
				break loop

			case CodeOperator:
				panic("TODO")

			case CodeMake:
				result = &Struct{Index: code.X, Values: stack}
				break loop

			case CodeVar:
				index := int32(len(data)) - code.X - 1
				value = data[index]
				goto beginning

			case CodeGlobal:
				result = Reduce(globals, globals[code.X])
				globals[code.X] = result
				break loop

			case CodeAbst:
				pop := int32(len(stack)) - code.X
				data = make([]Value, code.X)
				copy(data, stack[:pop])
				stack = stack[pop:]
				code = &code.Table[0]

			case CodeAppl:
				for i := range code.Table[:len(code.Table)-1] {
					switch code.Table[i].Kind {
					case CodeValue:
						stack = append(stack, code.Table[i].Value)
					case CodeVar:
						index := int32(len(data)) - code.X - 1
						stack = append(stack, data[index])
					default:
						stack = append(stack, &Thunk{Result: nil, Code: &code.Table[i], Data: data})
					}
				}
				code = &code.Table[len(code.Table)-1]

			case CodeSwitch:
				str := Reduce(globals, &Thunk{Result: nil, Code: &code.Table[0], Data: data}).(*Struct)
				data = append(data, str.Values...)
				code = &code.Table[str.Index+1]
			}
		}

		v.Result = result
		return result
	}
	panic("unreachable")
}
