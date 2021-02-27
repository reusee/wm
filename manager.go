package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type ManageWindow func(xproto.Window)

type UnmanageWindow func(xproto.Window)

func (_ Def) ManageWindow(
	conn *xgb.Conn,
) (
	manage ManageWindow,
	unmanage UnmanageWindow,
) {

	manage = func(win xproto.Window) {
		pt("manage %v\n", win)
	}

	unmanage = func(win xproto.Window) {
		pt("unmanage %v\n", win)
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
