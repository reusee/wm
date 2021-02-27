package main

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type BecomeWM func()

func (_ Def) BecomeWM(
	conn *xgb.Conn,
	rootWin xproto.Window,
) BecomeWM {
	return func() {
		err := xproto.ChangeWindowAttributesChecked(
			conn, rootWin, xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskButtonPress |
					xproto.EventMaskButtonRelease |
					xproto.EventMaskPointerMotion |
					xproto.EventMaskStructureNotify |
					xproto.EventMaskSubstructureRedirect,
			},
		).Check()
		var accessError xproto.AccessError
		if as(err, &accessError) {
			ce(fmt.Errorf("could not become the window manager"))
		}
		ce(err)
	}
}
