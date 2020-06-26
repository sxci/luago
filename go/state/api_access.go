package state

import (
	"fmt"

	"github.com/sxci/luago/go/api"
)

func (state *luaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRINNG:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (state *luaState) Type(idx int) api.LuaType {
	if state.stack.isValid(idx) {
		val := state.stack.get(idx)
		return typeOf(val)
	}
	return api.LUA_TNONE
}

func (state *luaState) IsNone(idx int) bool {
	return state.Type(idx) == api.LUA_TNONE
}

func (state *luaState) IsNil(idx int) bool {
	return state.Type(idx) == api.LUA_TNIL
}

func (state *luaState) IsNoneOrNil(idx int) bool {
	return state.Type(idx) <= api.LUA_TNIL
}

func (state *luaState) IsBoolean(idx int) bool {
	return state.Type(idx) == api.LUA_TBOOLEAN
}

func (state *luaState) IsString(idx int) bool {
	t := state.Type(idx)
	return t == api.LUA_TSTRINNG || t == api.LUA_TNUMBER
}

func (state *luaState) IsNumber(idx int) bool {
	_, ok := state.ToNumberX(idx)
	return ok
}

func (state *luaState) IsInteger(idx int) bool {
	val := state.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (state *luaState) ToBoolean(idx int) bool {
	val := state.stack.get(idx)
	return convertToBoolean(val)
}

func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func (state *luaState) ToNumber(idx int) float64 {
	n, _ := state.ToNumberX(idx)
	return n
}

func (state *luaState) ToNumberX(idx int) (float64, bool) {
	val := state.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (state *luaState) ToInteger(idx int) int64 {
	n, _ := state.ToIntegerX(idx)
	return n
}

func (state *luaState) ToIntegerX(idx int) (int64, bool) {
	val := state.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return int64(x), true
	case int64:
		return x, true
	default:
		return 0, false
	}
}

func (state *luaState) ToString(idx int) string {
	s, _ := state.ToStringX(idx)
	return s
}

func (state *luaState) ToStringX(idx int) (string, bool) {
	val := state.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		state.stack.set(idx, s) // ?
		return s, true
	default:
		return "", false
	}
}
