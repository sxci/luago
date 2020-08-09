package vm

import (
	"github.com/sxci/luago/go/api"
)

func move(i Instruction, vm api.LuaVM) { // copy
	a, b, _ := i.ABC()
	a++
	b++
	vm.Copy(b, a)
}

func jmp(i Instruction, vm api.LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("todo!")
	}
}
