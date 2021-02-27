package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type WindowsMap = map[xproto.Window]*Window

func (_ Def) WindowsMap() WindowsMap {
	return make(WindowsMap)
}

func (_ Def) WindowsArray(
	m WindowsMap,
) (array []*Window) {
	for _, win := range m {
		array = append(array, win)
	}
	return
}

type ManageWindow func(xproto.Window)

type UnmanageWindow func(xproto.Window)

func (_ Def) ManageWindow(
	conn *xgb.Conn,
	winsMap WindowsMap,
	update Update,
) (
	manage ManageWindow,
	unmanage UnmanageWindow,
) {

	manage = func(id xproto.Window) {
		if _, ok := winsMap[id]; ok {
			return
		}

		win := &Window{
			XID: id,
		}
		winsMap[id] = win
		update(&winsMap)
	}

	unmanage = func(win xproto.Window) {
		delete(winsMap, win)
		update(&winsMap)
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
