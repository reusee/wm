package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type StartMouseMovingWindow func(
	ev xproto.Window,
	x int16,
	y int16,
)

type StopMouseMovingWindow func()

type MouseMoveWindow func(int16, int16)

func (_ Def) MouseMovingWindow(
	conn *xgb.Conn,
	setup *xproto.SetupInfo,
) (
	start StartMouseMovingWindow,
	stop StopMouseMovingWindow,
	move MouseMoveWindow,
) {

	var moving xproto.Window
	var lastX, lastY int16

	start = func(win xproto.Window, x int16, y int16) {
		pt("start moving %d\n", win)
		moving = win
		lastX = x
		lastY = y

		// grab
		_, err := xproto.GrabPointer(
			conn, false, win,
			xproto.EventMaskPointerMotion|
				xproto.EventMaskButtonPress|
				xproto.EventMaskButtonRelease,
			xproto.GrabModeAsync,
			xproto.GrabModeAsync,
			xproto.WindowNone,
			xproto.CursorNone,
			xproto.TimeCurrentTime,
		).Reply()
		ce(err)

	}

	stop = func() {
		pt("stop moving %d\n", moving)
		moving = 0
		ce(xproto.UngrabPointerChecked(
			conn, xproto.TimeCurrentTime,
		).Check())
	}

	move = func(x, y int16) {
		if moving == 0 {
			return
		}
		deltaX := x - lastX
		deltaY := y - lastY
		if deltaX == 0 && deltaY == 0 {
			return
		}
		pt("delta %d %d\n", deltaX, deltaY)
		geom, err := xproto.GetGeometry(conn, xproto.Drawable(moving)).Reply()
		ce(err)
		pt("geometry %+v\n", geom)
		offset, err := xproto.TranslateCoordinates(
			conn, moving, setup.DefaultScreen(conn).Root,
			geom.X, geom.Y,
		).Reply()
		ce(err)
		pt("offset %+v\n", offset)
		winX := offset.DstX + deltaX
		winY := offset.DstY + deltaY
		pt("move %d to %d %d\n", moving, winX, winY)
		ce(xproto.ConfigureWindowChecked(conn, moving,
			xproto.ConfigWindowX|xproto.ConfigWindowY,
			[]uint32{
				uint32(winX),
				uint32(winY),
			},
		).Check())
		lastX = x
		lastY = y
	}

	return
}
