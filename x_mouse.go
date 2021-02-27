package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type SetupMouse func()

func (_ Def) SetupMouse(
	conn *xgb.Conn,
	setup *xproto.SetupInfo,
) SetupMouse {
	return func() {

		for _, screen := range setup.Roots {
			// async grab all buttons
			ce(xproto.GrabButtonChecked(
				conn,
				true,
				screen.Root,
				xproto.EventMaskButtonPress|
					xproto.EventMaskButtonRelease,
				xproto.GrabModeSync,
				xproto.GrabModeSync,
				xproto.WindowNone,
				xproto.CursorNone,
				xproto.ButtonIndexAny,
				xproto.ModMaskAny,
			).Check())
		}

	}
}
