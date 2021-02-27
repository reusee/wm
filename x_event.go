package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type SetupEventHandler func()

func (_ Def) SetupEventHandler(
	conn *xgb.Conn,
	cur *Scope,
) SetupEventHandler {
	return func() {

		go func() {

			for {
				ev, err := conn.WaitForEvent()
				if ev == nil && err != nil {
					ce(err)
				}

				if err != nil {
					pt("%s\n", err.Error())
				}

				if ev != nil {

					switch ev := ev.(type) {

					case xproto.ConfigureRequestEvent:
						var vals []uint32
						flags := ev.ValueMask
						if xproto.ConfigWindowX&flags > 0 {
							vals = append(vals, uint32(ev.X))
						}
						if xproto.ConfigWindowY&flags > 0 {
							vals = append(vals, uint32(ev.Y))
						}
						if xproto.ConfigWindowWidth&flags > 0 {
							vals = append(vals, uint32(ev.Width))
						}
						if xproto.ConfigWindowHeight&flags > 0 {
							vals = append(vals, uint32(ev.Height))
						}
						if xproto.ConfigWindowBorderWidth&flags > 0 {
							vals = append(vals, 0) // do not set border width
						}
						if xproto.ConfigWindowSibling&flags > 0 {
							vals = append(vals, uint32(ev.Sibling))
						}
						if xproto.ConfigWindowStackMode&flags > 0 {
							vals = append(vals, uint32(ev.StackMode))
						}
						xproto.ConfigureWindow(conn, ev.Window, flags, vals)

					case xproto.MapRequestEvent:
						xproto.MapWindow(conn, ev.Window)
						ce(xproto.ChangeWindowAttributesChecked(
							conn, ev.Window,
							xproto.CwEventMask,
							[]uint32{
								xproto.EventMaskPropertyChange,
							},
						).Check())
						cur.Call(func(
							manage ManageWindow,
							stack StackByLastFocus,
						) {
							manage(ev.Window)
							stack()
						})

					case xproto.UnmapNotifyEvent:
						cur.Call(func(
							unmanage UnmanageWindow,
							stack StackByLastFocus,
						) {
							unmanage(ev.Window)
							stack()
						})

					case xproto.CreateNotifyEvent:
					case xproto.PropertyNotifyEvent:
					case xproto.ExposeEvent:
					case xproto.ConfigureNotifyEvent:
					case xproto.ClientMessageEvent:
					case xproto.MapNotifyEvent:
					case xproto.MappingNotifyEvent:
					case xproto.DestroyNotifyEvent:

					default:
						pt("EVENT-> %v\n", ev)

					}
				}

			}

		}()

	}
}
