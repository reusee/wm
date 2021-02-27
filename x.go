package main

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
	"github.com/reusee/e4"
)

type XConn struct {
	*xgb.Conn
}

type XSetupInfo struct {
	*xproto.SetupInfo
}

type XRootWindow struct {
	xproto.Window
}

func (_ Def) X() (
	xConn XConn,
	xSetupInfo XSetupInfo,
	xRootWindow XRootWindow,
) {

	conn, err := xgb.NewConn()
	ce(err)
	xConn.Conn = conn

	ce(xinerama.Init(conn))

	setupInfo := xproto.Setup(conn)
	if n := len(setupInfo.Roots); n != 1 {
		e4.Throw(fmt.Errorf("too many roots: %d", n))
	}
	xSetupInfo.SetupInfo = setupInfo

	xRootWindow.Window = setupInfo.Roots[0].Root

	return
}
