package main

import (
	"sort"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type StackWindows func(
	less func(windows []*Window, i, j int) bool,
)

func (def Def) StackByFocus(
	get GetWindowsArray,
	conn *xgb.Conn,
) StackWindows {
	return func(
		less func(windows []*Window, i, j int) bool,
	) {

		windows := get()
		if len(windows) < 2 {
			return
		}

		sort.SliceStable(windows, func(i, j int) bool {
			return less(windows, i, j)
		})

		for i, win := range windows {
			if i > 0 {
				prev := windows[i-1] // prev
				ce(xproto.ConfigureWindowChecked(
					conn, win.XID,
					xproto.ConfigWindowSibling|
						xproto.ConfigWindowStackMode,
					[]uint32{
						uint32(prev.XID),
						xproto.StackModeAbove,
					},
				).Check())
			}
		}

	}
}

func byInteractTime(windows []*Window, i, j int) bool {
	a := windows[i]
	b := windows[j]
	if !a.LastKey.Equal(b.LastKey) {
		return a.LastKey.Before(b.LastKey)
	}
	return a.LastFocus.Before(b.LastFocus)
}
