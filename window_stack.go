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
	winsMap WindowsMap,
) StackWindows {
	return func() []*Window {

		windows := get()
		if len(windows) < 2 {
			return windows
		}

		sort.SliceStable(windows, func(i, j int) bool {
			a := windows[i]
			aID := a.XID
			if a.TransientFor > 0 {
				if w, ok := winsMap[a.TransientFor]; ok {
					a = w
				}
			}
			b := windows[j]
			bID := b.XID
			if b.TransientFor > 0 {
				if w, ok := winsMap[b.TransientFor]; ok {
					b = w
				}
			}
			if a.XID == b.XID {
				// same transient group
				return aID < bID
			}
			if a.Layer != b.Layer {
				return a.Layer < b.Layer
			}
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
