package main

import (
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
) (
	manage ManageWindow,
	unmanage UnmanageWindow,
) {

	manage = func(id xproto.Window) {
		if _, ok := winsMap[id]; ok {
			return
		}

		win := &Window{
			XID:       id,
			LastFocus: time.Now(),
		}

		r, err := xproto.GetProperty(
			conn, false, id,
			AtomWM_TRANSIENT_FOR,
			xproto.GetPropertyTypeAny, 0, 60,
		).Reply()
		ce(err)
		if len(r.Value) > 0 {
			win.TransientFor = xproto.Window(xgb.Get32(r.Value))
		}

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
