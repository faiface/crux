package crux

import (
	"fmt"

	"github.com/faiface/crux/runtime"
)

type link struct {
	Name  string
	Index int32
}

func isFast(e Expr) bool {
	switch e := e.(type) {
	case *Char, *Int, *Float, *Operator, *Make, *Var, *Abst:
		return true
	case *Appl:
		if !isFast(e.Rator) {
			return false
		}
		for _, rand := range e.Rands {
			switch rand := rand.(type) {
			case *Var, *Strict:
				continue
			default:
				if hasLocals(rand) {
					return false
				}
			}
		}
		return true
	case *Switch:
		for _, cas := range e.Cases {
			if !isFast(cas) {
				return false
			}
		}
		return true
	default:
		panic("unreachable")
	}
}

func hasLocals(e Expr) bool {
	switch e := e.(type) {
	case *Char, *Int, *Float, *Operator, *Make:
		return false
	case *Var:
		return e.Index < 0
	case *Abst:
		return false
	case *Appl:
		if hasLocals(e.Rator) {
			return true
		}
		for _, rand := range e.Rands {
			if hasLocals(rand) {
				return true
			}
		}
		return false
	case *Strict:
		return hasLocals(e.Expr)
	case *Switch:
		if hasLocals(e.Expr) {
			return true
		}
		for _, cas := range e.Cases {
			if hasLocals(cas) {
				return true
			}
		}
		return false
	default:
		panic("unreachable")
	}
}

func Compile(globals map[string][]Expr) (
	globalIndices map[string][]int32,
	globalValues []runtime.Value,
	codes []runtime.Code,
) {
	// total hack, compile the first time just to get the number of codes
	// compile second time so that tables all refer the same codes slice
	_, _, codes = compile(0, globals)
	return compile(len(codes), globals)
}

func compile(alloc int, globals map[string][]Expr) (
	globalIndices map[string][]int32,
	globalValues []runtime.Value,
	codes []runtime.Code,
) {
	links := make(map[int]link)

	var process = func(i int) func(c runtime.Code, ln *link) runtime.Code {
		return func(c runtime.Code, ln *link) runtime.Code {
			if ln != nil {
				links[i] = *ln
			}
			return c
		}
	}

	var compile func(locals []string, e Expr) (runtime.Code, *link)
	compile = func(locals []string, e Expr) (runtime.Code, *link) {
		switch e := e.(type) {
		case *Char:
			return runtime.Code{
				Kind:  runtime.CodeValue,
				Value: &runtime.Char{Value: e.Value},
			}, nil

		case *Int:
			return runtime.Code{
				Kind:  runtime.CodeValue,
				Value: &runtime.Int{Value: e.Value},
			}, nil

		case *Float:
			return runtime.Code{
				Kind:  runtime.CodeValue,
				Value: &runtime.Float{Value: e.Value},
			}, nil

		case *Operator:
			return runtime.Code{
				Kind: runtime.CodeOperator,
				X:    e.Code,
			}, nil

		case *Make:
			return runtime.Code{
				Kind: runtime.CodeMake,
				X:    e.Index,
			}, nil

		case *Var:
			if e.Index >= 0 {
				return runtime.Code{Kind: runtime.CodeGlobal}, &link{e.Name, e.Index}
			}
			for i := len(locals) - 1; i >= 0; i-- {
				if e.Name == locals[i] {
					return runtime.Code{
						Kind: runtime.CodeVar,
						X:    int32(i),
					}, nil
				}
			}
			panic(fmt.Sprintf("%s not bound", e.Name))

		case *Abst:
			i := len(codes)
			codes = append(codes, runtime.Code{})
			codes[i] = process(i)(compile(e.Bound, e.Body))
			kind := runtime.CodeAbst
			if isFast(e.Body) {
				kind = runtime.CodeFastAbst
			}
			return runtime.Code{
				Kind:  kind,
				X:     int32(len(e.Bound)),
				Table: codes[i : i+1],
			}, nil

		case *Appl:
			i := len(codes)
			codes = append(codes, make([]runtime.Code, 1+len(e.Rands))...)
			codes[i] = process(i)(compile(locals, e.Rator))
			for j := 0; j < len(e.Rands); j++ {
				codes[i+1+j] = process(i + 1 + j)(compile(locals, e.Rands[j]))
			}
			return runtime.Code{
				Kind:  runtime.CodeAppl,
				Table: codes[i : i+1+len(e.Rands)],
			}, nil

		case *Strict:
			i := len(codes)
			codes = append(codes, runtime.Code{})
			codes[i] = process(i)(compile(locals, e.Expr))
			return runtime.Code{
				Kind:  runtime.CodeStrict,
				Table: codes[i : i+1],
			}, nil

		case *Switch:
			i := len(codes)
			codes = append(codes, make([]runtime.Code, 1+len(e.Cases))...)
			codes[i] = process(i)(compile(locals, e.Expr))
			for j := 0; j < len(e.Cases); j++ {
				codes[i+1+j] = process(i + 1 + j)(compile(locals, e.Cases[j]))
			}
			return runtime.Code{
				Kind:  runtime.CodeSwitch,
				Table: codes[i : i+1+len(e.Cases)],
			}, nil
		}
		panic("unreachable")
	}

	globalIndices = make(map[string][]int32)
	globalValues = nil
	codes = make([]runtime.Code, 0, alloc)

	for name := range globals {
		for index := range globals[name] {
			i := len(codes)
			codes = append(codes, runtime.Code{})
			codes[i] = process(i)(compile(nil, globals[name][index]))
			globalIndices[name] = append(globalIndices[name], int32(len(globalValues)))
			switch codes[i].Kind {
			case runtime.CodeValue:
				globalValues = append(globalValues, codes[i].Value)
			default:
				globalValues = append(globalValues, &runtime.Thunk{Code: &codes[i]})
			}
		}
	}

	for i, ln := range links {
		codes[i] = runtime.Code{
			Kind: runtime.CodeGlobal,
			X:    globalIndices[ln.Name][ln.Index],
		}
	}

	return globalIndices, globalValues, codes
}
