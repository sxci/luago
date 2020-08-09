package binchunk

import (
	"encoding/binary"
	"math"
	"unsafe"
)

type reader struct {
	data []byte
}

func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

func (r *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

func (r *reader) readString() string {
	size := uint(r.readByte())

	if size == 0 { // Nil 字符串
		return ""
	}
	if size == 0xFF { // 长字符串
		size = uint(r.readUint64())
	}
	bytes := r.readBytes(size - 1)
	return string(bytes)
}

func bytesToString(bytes *[]byte) *string {
	return (*string)(unsafe.Pointer(&bytes))
}

func (r *reader) readBytes(n uint) []byte {
	bytes := r.data[:n]
	r.data = r.data[n:]
	return bytes
}

func (r *reader) checkHeader() {
	if string(r.readBytes(4)) != LUA_SIGNATURE {
		errReport("not a precompiled chunk!")
	} else if r.readByte() != LUA_VERSION {
		errReport("version mismatch!")
	} else if r.readByte() != LUAC_FORMAT {
		errReport("format mismatch!")
	} else if string(r.readBytes(6)) != LUAC_DATA {
		errReport("corrupted!")
	} else if r.readByte() != CINT_SIZE {
		errReport("int size mismatch!")
	} else if r.readByte() != CSIZET_SIZE {
		errReport("size_t size mismatch!")
	} else if r.readByte() != INSTRUCTION_SIZE {
		errReport("instruction size mismatch!")
	} else if r.readByte() != LUAINTEGER_SIZE {
		errReport("lua_Integer size mismatch!")
	} else if r.readByte() != LUANUMBER_SIZE {
		errReport("lua_Number size mismatch!")
	} else if r.readLuaInteger() != LUAC_INT {
		errReport("endianness mismatch!")
	} else if r.readLuaNumber() != LUAC_NUM {
		errReport("float format mismatch!")
	}
}

func (r *reader) readProto(parentSource string) *Prototype {
	source := r.readString()
	if source == "" {
		source = parentSource
	}
	// 能够保证从上到下的执行顺序 ?? 可以的。
	return &Prototype{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStacksize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *reader) readConstant() interface{} {
	switch r.readByte() { // tag
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() != 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STR:
		return r.readString()
	case TAG_LONG_STR:
		return r.readString()
	default:
		panic("corrupted!")
	}
}

func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		// go keep the reading order ?
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}
func (r *reader) readProtos(source string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(source)
	}
	return protos
}

func (r *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}
	return lineInfo
}

func (r *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	// go keep the reading order ?
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPc:   r.readUint32(),
		}
	}
	return locVars
}

func (r *reader) readUpvalueNames() []string {
	upvalueNames := make([]string, r.readUint32())
	for i := range upvalueNames {
		upvalueNames[i] = r.readString()
	}
	return upvalueNames
}

func errReport(msg interface{}) {
	panic(msg)
}
