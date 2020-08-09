package api

type LuaVM interface {
	LuaState
	PC() int          //返回当前 PC （仅用于测试）
	AddPC(n int)      // 修改 PC （用于实现跳转指令）
	Fetch() uint32    // 取出当前指令；将 PC 指向下一条指令
	GetConst(idx int) // 将制定常量推入栈顶
	GetRK(rk int)     // 将指定常量或栈值推入栈顶
}
