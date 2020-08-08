package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sxci/luago/go/api"
	"github.com/sxci/luago/go/binchunk"
	"github.com/sxci/luago/go/state"
	"github.com/sxci/luago/go/vm"
)

func main() {
	mainCalcState()
}

func mainCalcState() {
	ls := state.New()
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls)

	ls.Arith(api.LUA_OPADD)
	printStack(ls)
	ls.Arith(api.LUA_OPBNOT)
	printStack(ls)
	ls.Len(2)
	printStack(ls)
	ls.Concat(3)
	printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, api.LUA_OPEQ))
	printStack(ls)
}

func mainShowStack() {
	ls := state.New()
	ls.PushBoolean(true)
	printStack(ls)
	ls.PushInteger(10)
	printStack(ls)
	ls.PushNil()
	printStack(ls)
	ls.PushString("hello")
	printStack(ls)
	ls.PushValue(-4)
	printStack(ls)
	ls.Replace(3)
	printStack(ls)
	ls.SetTop(6)
	printStack(ls)
	ls.Remove(-3)
	printStack(ls)
	ls.SetTop(-5)
	printStack(ls)
}

func printStack(ls api.LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case api.LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case api.LUA_TSTRINNG:
			fmt.Printf("[%q]", ls.ToString(i))
		default:
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func mainPrintLuaChunk() {
	if len(os.Args) <= 1 {
		return
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	list(proto, true)
}

func list(f *binchunk.Prototype, detail bool) {
	printHeader(f)
	printCode(f)
	if detail {
		printDetail(f)
	}
	for _, p := range f.Protos {
		list(p, detail)
	}
}

func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}
	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf("\n%s <%s:%d,%d> (%d instructions at %p)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code), f)
	fmt.Printf("%d%s params, %d slots, %d upvalue, %d locals, %d constants, %d functions\n",
		f.NumParams, varargFlag, f.MaxStacksize, len(f.Upvalues), len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := vm.Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printOperands(i vm.Instruction) {
	str := ""
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		pb, pc := "", ""
		if i.BMode() != vm.OpArgN {
			if b > 0xFF {
				pb = strconv.Itoa(-1 - (b & 0xFF))
			} else {
				pb = strconv.Itoa(b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0xFF {
				pc = strconv.Itoa(-1 - (c & 0xFF))
			} else {
				pc = strconv.Itoa(c)
			}
		}
		str = fmt.Sprintf("%d %2s %2s", a, pb, pc)
	case vm.IABx:
		a, bx := i.ABx()
		pbx := bx
		if i.BMode() == vm.OpArgK {
			pbx = -1 - bx
		}
		str = fmt.Sprintf("%d %2d", a, pbx)
	case vm.IAsBx:
		a, sbx := i.AsBx()
		str = fmt.Sprintf("%d %2d", a, sbx)
	case vm.IAx:
		ax := i.Ax()
		str = fmt.Sprintf("%d", ax)
	}
	fmt.Printf("%-8s", str)
}

func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d) for %p:\n", len(f.Constants), f)
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}
	fmt.Printf("locals (%d) for %p:\n", len(f.LocVars), f)
	for i, l := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i+1, l.VarName, l.StartPC+1, l.EndPc+1)
	}
	fmt.Printf("upvalues (%d) for %p:\n", len(f.Upvalues), f)
	for i, u := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i+1, upvalName(f, i), u.Instack, u.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}
