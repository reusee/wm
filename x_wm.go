package main

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/reusee/e4"
)

type BecomeWM func()

func (_ Def) BecomeWM(
	conn *xgb.Conn,
	setup *xproto.SetupInfo,
) BecomeWM {
	return func() {
		for i, screen := range setup.Roots {
			err := xproto.ChangeWindowAttributesChecked(
				conn, screen.Root, xproto.CwEventMask,
				[]uint32{
					0 |
						//xproto.EventMaskButtonRelease |
						//xproto.EventMaskKeyPress |
						//xproto.EventMaskKeyRelease |
						//xproto.EventMaskPointerMotion |
						xproto.EventMaskButtonPress |
						xproto.EventMaskPropertyChange |
						xproto.EventMaskStructureNotify |
						xproto.EventMaskSubstructureNotify |
						xproto.EventMaskSubstructureRedirect |
						0,
				},
			).Check()
			var accessError xproto.AccessError
			if as(err, &accessError) {
				ce(
					fmt.Errorf("could not become the window manager"),
					e4.WithInfo("screen %d", i),
				)
			}
			ce(err)
		}
	}
}
