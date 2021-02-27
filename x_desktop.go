package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type SetupDesktop func()

func (_ Def) SetupDeskTop(
	desktopWin DesktopXWindow,
) SetupDesktop {
	return func() {

	}
}

type XCursor xproto.Cursor

func (_ Def) XCursor(
	conn *xgb.Conn,
) XCursor {
	cursorFont, err := xproto.NewFontId(conn)
	ce(err)
	ce(xproto.OpenFontChecked(conn, cursorFont, uint16(len("cursor")), "cursor").Check())
	cursor, err := xproto.NewCursorId(conn)
	ce(err)
	const leftPtr = 80
	ce(xproto.CreateGlyphCursorChecked(conn,
		cursor, cursorFont, cursorFont,
		leftPtr, leftPtr+1,
		0xffff, 0xffff, 0xffff,
		0, 0, 0,
	).Check())
	ce(xproto.CloseFontChecked(conn, cursorFont).Check())
	return XCursor(cursor)
}

type DesktopXWindow xproto.Window

type DesktopXGContext xproto.Gcontext

func (_ Def) DesktopVars(
	conn *xgb.Conn,
	screen xproto.ScreenInfo,
	width DesktopWidth,
	height DesktopHeight,
	cursor XCursor,
) (
	win DesktopXWindow,
	ctx DesktopXGContext,
) {

	id, err := xproto.NewWindowId(conn)
	ce(err)
	win = DesktopXWindow(id)

	gctx, err := xproto.NewGcontextId(conn)
	ce(err)
	ctx = DesktopXGContext(gctx)

	// TODO xsettings

	ce(xproto.CreateWindowChecked(
		conn, screen.RootDepth, id, screen.Root,
		0, 0, uint16(width), uint16(height), 0,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		xproto.CwOverrideRedirect|xproto.CwEventMask,
		[]uint32{
			1,
			xproto.EventMaskExposure,
		},
	).Check())

	ce(xproto.ConfigureWindowChecked(
		conn, id,
		xproto.ConfigWindowStackMode,
		[]uint32{
			xproto.StackModeBelow,
		},
	).Check())

	ce(xproto.ChangeWindowAttributesChecked(
		conn, id,
		xproto.CwBackPixel|xproto.CwCursor,
		[]uint32{
			screen.WhitePixel,
			uint32(cursor),
		},
	).Check())

	font, err := xproto.NewFontId(conn)
	ce(err)
	ce(xproto.OpenFontChecked(
		conn, font,
		uint16(len("6x13")), "6x13",
	).Check())
	defer xproto.CloseFont(conn, font)

	ce(xproto.CreateGCChecked(
		conn, gctx,
		xproto.Drawable(screen.Root),
		xproto.GcFont,
		[]uint32{
			uint32(font),
		}).Check())

	ce(xproto.MapWindowChecked(
		conn, id,
	).Check())

	return
}

type DesktopWidth int

type DesktopHeight int

func (_ Def) DesktopSize(
	info *xproto.SetupInfo,
) (
	w DesktopWidth,
	h DesktopHeight,
) {
	w = DesktopWidth(info.Roots[0].WidthInPixels)
	h = DesktopHeight(info.Roots[0].HeightInPixels)
	return
}
