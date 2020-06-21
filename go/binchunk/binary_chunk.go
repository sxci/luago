package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUA_VERSION      = 0x53 // 5 * 16 + 3
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\x0D\x0A\x1A\x0A" // "\x19\x93\r\n\x1A\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUAINTEGER_SIZE  = 8
	LUANUMBER_SIZE   = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header                  //头部
	sizeUpvalues byte       //主函数 upvalue 变量
	mainFunc     *Prototype //主函数原型
}

type header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	csizetSize      byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64
	luacNum         float64
}

type Prototype struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStacksize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPc   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader() // header
	reader.readByte()    // upvalue size in main
	return reader.readProto("")
}
