package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type DefaultCursor xproto.Cursor

func (_ Def) DefaultCursor(
	conn *xgb.Conn,
) DefaultCursor {
	cursorFont, err := xproto.NewFontId(conn)
	ce(err)
	ce(xproto.OpenFontChecked(conn, cursorFont, uint16(len("cursor")), "cursor").Check())
	cursor, err := xproto.NewCursorId(conn)
	ce(err)
	const leftPtr = 68
	ce(xproto.CreateGlyphCursorChecked(conn,
		cursor, cursorFont, cursorFont,
		leftPtr, leftPtr+1,
		0xffff, 0xffff, 0xffff,
		0, 0, 0,
	).Check())
	ce(xproto.CloseFontChecked(conn, cursorFont).Check())
	return DefaultCursor(cursor)
}
