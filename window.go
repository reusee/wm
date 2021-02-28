package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Window struct {
	XID          xproto.Window
	LastFocus    time.Time
	LastRaise    time.Time
	Layer        int
	Tags         map[Tag]bool
	TransientFor xproto.Window
	Name         string
}

type WindowsMap = map[xproto.Window]*Window

func (_ Def) WindowsMap() WindowsMap {
	return make(WindowsMap)
}

type GetWindowsArray func() []*Window

func (_ Def) WindowsArray(
	m WindowsMap,
) GetWindowsArray {
	return func() []*Window {
		var array []*Window
		for _, win := range m {
			array = append(array, win)
		}
		return array
	}
}

type ManageWindow func(xproto.Window)

type UnmanageWindow func(xproto.Window)

func (_ Def) ManageWindow(
	conn *xgb.Conn,
	winsMap WindowsMap,
	cursor DefaultCursor,
	updateProperty UpdateWindowProperty,
) (
	manage ManageWindow,
	unmanage UnmanageWindow,
) {

	manage = func(id xproto.Window) {
		if _, ok := winsMap[id]; ok {
			return
		}

		// set event mark
		ce(xproto.ChangeWindowAttributesChecked(
			conn, id,
			xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskPropertyChange |
					xproto.EventMaskEnterWindow,
			},
		).Check())

		// set cursor
		ce(xproto.ChangeWindowAttributesChecked(
			conn, id,
			xproto.CwCursor,
			[]uint32{
				uint32(cursor),
			},
		).Check())

		win := &Window{
			XID:       id,
			LastFocus: time.Now(),
		}

		// properties
		updateProperty(id, AtomWM_TRANSIENT_FOR, &win.TransientFor)
		updateProperty(id, Atom_NET_WM_NAME, &win.Name)

		winsMap[id] = win
	}

	unmanage = func(win xproto.Window) {
		delete(winsMap, win)
	}

	return
}

type ManageExistingWindows func()

func (_ Def) ManageExistingWindows(
	setupInfo *xproto.SetupInfo,
	conn *xgb.Conn,
	desktopWins DesktopWindows,
	manage ManageWindow,
) ManageExistingWindows {
	return func() {

		for _, screen := range setupInfo.Roots {
			tree, err := xproto.QueryTree(conn, screen.Root).Reply()
			ce(err)
			if tree != nil {
				for _, win := range tree.Children {
					if win == desktopWins[screen.Root] {
						continue
					}
					attrs, err := xproto.GetWindowAttributes(conn, win).Reply()
					if attrs == nil || err != nil {
						continue
					}
					if attrs.OverrideRedirect || attrs.MapState == xproto.MapStateUnmapped {
						continue
					}
					manage(win)
				}
			}
		}

	}
}

type UpdateWindowProperty func(
	win xproto.Window,
	atom xproto.Atom,
	target any,
)

func (_ Def) UpdateWindowProperty(
	conn *xgb.Conn,
) UpdateWindowProperty {
	return func(
		win xproto.Window,
		atom xproto.Atom,
		target any,
	) {

		r, err := xproto.GetProperty(
			conn, false, win,
			atom, xproto.GetPropertyTypeAny,
			0, 60,
		).Reply()
		var windowError xproto.WindowError
		if as(err, &windowError) {
			return
		}
		ce(err)

		if len(r.Value) > 0 {
			targetValue := reflect.ValueOf(target)
			switch targetValue.Elem().Kind() {
			case reflect.Uint32:
				targetValue.Elem().SetUint(
					uint64(xgb.Get32(r.Value)),
				)
			case reflect.String:
				start := 0
				var strs []string
				for i, c := range r.Value {
					if c == 0 {
						strs = append(strs, string(r.Value[start:i]))
						start = i + 1
					}
				}
				if start < int(r.ValueLen) {
					strs = append(strs, string(r.Value[start:]))
				}
				//TODO cut or join?
				if len(strs) > 0 {
					targetValue.Elem().SetString(strs[0])
				}
			default:
				panic(fmt.Errorf("bad target type: %T", target))
			}
		}

	}
}
