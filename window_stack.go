package main

import (
	"sort"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type StackByLastFocus func()

func (def Def) StackByFocus(
	get GetWindowsArray,
	conn *xgb.Conn,
) StackByLastFocus {
	return func() {

		windows := get()
		if len(windows) < 2 {
			return
		}

		sort.SliceStable(windows, func(i, j int) bool {
			return windows[i].LastFocus.Before(windows[j].LastFocus)
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
