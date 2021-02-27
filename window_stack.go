package main

import (
	"sort"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type AdjustWindow func(i, j int)

func (def Def) AdjustWindow(
	windows []*Window,
	conn *xgb.Conn,
) AdjustWindow {
	return func(i, j int) {
		if j == 0 {
			// below next
			sibling := windows[1].XID
			ce(xproto.ConfigureWindowChecked(
				conn, windows[i].XID,
				xproto.ConfigWindowSibling|
					xproto.ConfigWindowStackMode,
				[]uint32{
					uint32(sibling),
					xproto.StackModeBelow,
				},
			).Check())

		} else {
			// above previous
			sibling := windows[j-1].XID
			ce(xproto.ConfigureWindowChecked(
				conn, windows[i].XID,
				xproto.ConfigWindowSibling|
					xproto.ConfigWindowStackMode,
				[]uint32{
					uint32(sibling),
					xproto.StackModeAbove,
				},
			).Check())
		}
	}
}

type StackByLastFocus func()

func (def Def) StackByFocus(
	windows []*Window,
	update Update,
	adjust AdjustWindow,
) StackByLastFocus {
	return func() {

		if len(windows) < 2 {
			return
		}

		updated := false
		sort.Sort(WindowsSorter{
			windows: windows,
			less: func(i, j int) bool {
				return windows[i].LastFocus.Before(windows[j].LastFocus)
			},
			swap: func(i, j int) {
				updated = true
				adjust(i, j)
				adjust(j, i)
			},
		})
		if updated {
			update(def.WindowsArray)
		}

	}
}
