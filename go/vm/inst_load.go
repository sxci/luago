package vm

import "github.com/sxci/luago/go/api"

func loadNil(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a++
	vm.PushNil()
	for i := 0; i < a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

func loadBool(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a++
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

func loadK(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a++
	vm.GetConst(bx)
	vm.Replace(a)
}

func loadKx(i Instruction, vm api.LuaVM) {
	a, _ := i.ABx()
	a++
	ax := Instruction(vm.Fetch()).Ax()
	vm.GetConst(ax)
	vm.Replace(a)
}
