package filter

import (
	"encoding/json"
	mep "github.com/ChenHaoHu/ExpressionParser/ep"
	log "github.com/mhchlib/logger"
	lua "github.com/yuin/gopher-lua"
)

func CalMepFilter(code string, metatdata map[string]string) bool {
	epEngine, err := mep.NewEpEngine(code)
	if err != nil {
		log.Error("mep filter code parse err", err.Error(), "code", code)
	}
	check := epEngine.Check(metatdata)
	return check
}

func CalSimpleFilter(code string, metatdata map[string]string) bool {
	metajson, _ := json.Marshal(metatdata)
	if string(metajson) == code {
		return true
	}
	return false
}

func CalLuaFilter(code string, metatdata map[string]string) bool {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(code); err != nil {
		log.Error("lua filter code parse err", err.Error(), "code", code)
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
		log.Error("lua filter code parse err", err.Error(), "code", code)
		return false
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)
	if ret.String() == "true" {
		return true
	}
	return false
}
