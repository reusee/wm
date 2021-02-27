package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var (
	Atom_NET_ACTIVE_WINDOW,
	Atom_NET_WM_NAME,
	AtomWINDOW,
	AtomWM_CLASS,
	AtomWM_DELETE_WINDOW,
	AtomWM_PROTOCOLS,
	AtomWM_TAKE_FOCUS,
	AtomWM_TRANSIENT_FOR xproto.Atom
)

type InternAtoms func()

func (_ Def) InternAtoms(
	xConn *xgb.Conn,
) InternAtoms {
	return func() {

		intern := func(str string) xproto.Atom {
			r, err := xproto.InternAtom(xConn, false, uint16(len(str)), str).Reply()
			ce(err, wi("intern %s", str))
			if r == nil {
				return 0
			}
			return r.Atom
		}

		Atom_NET_ACTIVE_WINDOW = intern("_NET_ACTIVE_WINDOW")
		Atom_NET_WM_NAME = intern("_NET_WM_NAME")
		AtomWINDOW = intern("WINDOW")
		AtomWM_CLASS = intern("WM_CLASS")
		AtomWM_DELETE_WINDOW = intern("WM_DELETE_WINDOW")
		AtomWM_PROTOCOLS = intern("WM_PROTOCOLS")
		AtomWM_TAKE_FOCUS = intern("WM_TAKE_FOCUS")
		AtomWM_TRANSIENT_FOR = intern("WM_TRANSIENT_FOR")

	}
}
