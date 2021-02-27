package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type SetupAllDesktop func()

func (_ Def) SetupAllDesktop(
	wins DesktopWindows,
) SetupAllDesktop {
	return func() {
		_ = wins
	}
}

type DesktopWindows map[xproto.Window]xproto.Window

func (_ Def) DesktopWindows(
	setupInfo *xproto.SetupInfo,
	setup SetupDesktop,
) DesktopWindows {
	m := make(DesktopWindows)
	for _, screen := range setupInfo.Roots {
		m[screen.Root] = setup(screen)
	}
	return m
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

type SetupDesktop func(xproto.ScreenInfo) xproto.Window

func (_ Def) SetupDesktop(
	conn *xgb.Conn,
	cursor XCursor,
) SetupDesktop {

	return func(
		screen xproto.ScreenInfo,
	) (
		win xproto.Window,
	) {
		var err error

		win, err = xproto.NewWindowId(conn)
		ce(err)

		gctx, err := xproto.NewGcontextId(conn)
		ce(err)

		// TODO xsettings

		ce(xproto.CreateWindowChecked(
			conn, screen.RootDepth, win, screen.Root,
			0, 0,
			uint16(screen.WidthInPixels),
			uint16(screen.HeightInPixels),
			0,
			xproto.WindowClassInputOutput,
			screen.RootVisual,
			xproto.CwOverrideRedirect|xproto.CwEventMask,
			[]uint32{
				1,
				xproto.EventMaskExposure,
			},
		).Check())

		ce(xproto.ConfigureWindowChecked(
			conn, win,
			xproto.ConfigWindowStackMode,
			[]uint32{
				xproto.StackModeBelow,
			},
		).Check())

		ce(xproto.ChangeWindowAttributesChecked(
			conn, win,
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
			conn, win,
		).Check())

		return
	}
}
