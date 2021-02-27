package main

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type KeyCodeSymMap map[KeyCode]map[KeySym]struct{}

func (_ Def) KeyCodeSymMap(
	conn *xgb.Conn,
) KeyCodeSymMap {
	m := make(KeyCodeSymMap)
	const keyLow = 8
	const keyHigh = 255
	mapping, err := xproto.GetKeyboardMapping(conn, keyLow, keyHigh-keyLow+1).Reply()
	ce(err)
	if mapping == nil {
		ce(fmt.Errorf("no keyboard mapping"))
	}
	n := int(mapping.KeysymsPerKeycode)
	for i := keyLow; i <= keyHigh; i++ {
		syms := make(map[KeySym]struct{})
		for j := 0; j < int(mapping.KeysymsPerKeycode); j++ {
			sym := mapping.Keysyms[(i-keyLow)*n+j]
			if sym > 0 {
				syms[KeySym(sym)] = struct{}{}
			}
		}
		m[KeyCode(i)] = syms
	}
	return m
}

type SetupKeyboard func()

func (_ Def) SetupKeyboard(
	keyCodeToSymMap KeyCodeSymMap,
) SetupKeyboard {
	return func() {

		_ = keyCodeToSymMap

		//TODO grab keys

	}
}
