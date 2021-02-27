package main

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
)

func (_ Def) X() (
	xConn *xgb.Conn,
	xSetupInfo *xproto.SetupInfo,
	xRootWindow xproto.Window,
) {
	var err error

	xConn, err = xgb.NewConn()
	ce(err)

	ce(xinerama.Init(xConn))

	xSetupInfo = xproto.Setup(xConn)
	if n := len(xSetupInfo.Roots); n != 1 {
		ce(fmt.Errorf("too many roots: %d", n))
	}

	xRootWindow = xSetupInfo.Roots[0].Root

	return
}

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
