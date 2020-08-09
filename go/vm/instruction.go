package vm

import "github.com/sxci/luago/go/api"

type Instruction uint32

const (
	MAXARG_BX  = 1<<18 - 1      // 2^18 - 1 = 262143
	MAXARG_sBx = MAXARG_BX >> 1 // 262143 / 2 = 131071
)

func (i Instruction) Opcode() int {
	return int(i & 0x3F) // 6 bits
}

// [  B:9  ][  C:9  ][ A:8  ][OP:6]
func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)   // 8 bits
	c = int(i >> 14 & 0x1FF) // 9 bits
	b = int(i >> 23 & 0x1FF) //9 bits
	return
}

// [      Bx:18     ][ A:8  ][OP:6]
func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14) // 18 bits
	return
}

// [     sBx:18     ][ A:8  ][OP:6]
func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

// [           Ax:26        ][OP:6]
func (i Instruction) Ax() int {
	return int(i >> 6)
}

func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

func (i Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

func (i Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}

func (i Instruction) Execute(vm api.LuaVM) {
	action := opcodes[i.Opcode()].action
	if action != nil {
		action(i, vm)
	} else {
		panic(i.OpName())
	}
}
