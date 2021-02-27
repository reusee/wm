package main

import (
	"sort"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type StackWindows func() []*Window

func (def Def) StackByFocus(
	get GetWindowsArray,
	conn *xgb.Conn,
) StackWindows {
	return func() []*Window {

		windows := get()
		if len(windows) < 2 {
			return windows
		}

		sort.SliceStable(windows, func(i, j int) bool {
			a := windows[i]
			b := windows[j]
			if !a.LastRaise.Equal(b.LastRaise) {
				return a.LastRaise.Before(b.LastRaise)
			}
			return a.LastFocus.Before(b.LastFocus)
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

		return windows
	}
}
