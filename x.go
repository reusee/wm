package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
)

func (_ Def) X() (
	conn *xgb.Conn,
	setupInfo *xproto.SetupInfo,
) {
	var err error

	conn, err = xgb.NewConn()
	ce(err)

	//TODO ewmh
	//TODO randr
	//TODO icccm

	ce(xinerama.Init(conn))

	setupInfo = xproto.Setup(conn)

	return
}
