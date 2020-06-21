package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sxci/luago/go/binchunk"
)

func main() {
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
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
	}
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
