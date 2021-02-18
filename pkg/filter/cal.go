package filter

import (
	log "github.com/mhchlib/logger"
	lua "github.com/yuin/gopher-lua"
)

func CalMepFilter(code string, metatdata map[string]string) bool {
	return false
}

func CalSimpleFilter(code string, metatdata map[string]string) bool {
	return false
}

func CalLuaFilter(code string, metatdata map[string]string) bool {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(code); err != nil {
		log.Error("lua filter code parser err", err.Error(), "code", code)
		return false
	}
	metatable := L.NewTable()
	for k, v := range metatdata {
		metatable.RawSetString(k, lua.LString(v))
	}
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("Filter"),
		NRet:    1,
		Protect: true,
	}, metatable); err != nil {
		log.Error("lua filter code parser err", err.Error(), "code", code)
		return false
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)
	if ret.String() == "true" {
		return true
	}
	return false
}
