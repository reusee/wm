package main

import (
	"sync"

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
	intern InternAtom,
) InternAtoms {
	return func() {

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

type InternAtom func(name string) xproto.Atom

func (_ Def) InternAtom(
	conn *xgb.Conn,
) InternAtom {
	return func(str string) xproto.Atom {
		r, err := xproto.InternAtom(conn, false, uint16(len(str)), str).Reply()
		ce(err, wi("intern %s", str))
		if r == nil {
			return 0
		}
		return r.Atom
	}
}

type AtomName func(xproto.Atom) string

func (_ Def) AtomName(
	conn *xgb.Conn,
) AtomName {
	var m sync.Map
	return func(atom xproto.Atom) string {
		v, ok := m.Load(atom)
		if ok {
			return v.(string)
		}
		r, err := xproto.GetAtomName(conn, atom).Reply()
		ce(err)
		name := r.Name
		m.Store(atom, name)
		return name
	}
}
