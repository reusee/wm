package main

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type SetupEventHandler func()

func (_ Def) SetupEventHandler(
	conn *xgb.Conn,
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

					case xproto.CreateNotifyEvent:
						ce(xproto.ChangeWindowAttributesChecked(
							conn, ev.Window,
							xproto.CwEventMask,
							[]uint32{
								xproto.EventMaskPropertyChange,
							},
						).Check())
						//TODO

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
						//TODO

					case xproto.MapRequestEvent:
						xproto.MapWindow(conn, ev.Window)
						//TODO

					case xproto.PropertyNotifyEvent:
						//TODO

					case xproto.ExposeEvent:
					case xproto.ConfigureNotifyEvent:
					case xproto.ClientMessageEvent:
					case xproto.MapNotifyEvent:

					default:
						pt("EVENT-> %v\n", ev)

					}
				}

			}

		}()

	}
}
