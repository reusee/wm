package main

import (
	"time"

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
						// map
						xproto.MapWindow(conn, ev.Window)

						// set event mark
						ce(xproto.ChangeWindowAttributesChecked(
							conn, ev.Window,
							xproto.CwEventMask,
							[]uint32{
								xproto.EventMaskPropertyChange |
									xproto.EventMaskEnterWindow,
							},
						).Check())

						// manage
						cur.Call(func(
							manage ManageWindow,
							stack StackWindows,
						) {
							manage(ev.Window)
							stack(byInteractTime)
						})

					case xproto.UnmapNotifyEvent:
						cur.Call(func(
							unmanage UnmanageWindow,
							stack StackWindows,
						) {
							unmanage(ev.Window)
							stack(byInteractTime)
						})

					case xproto.EnterNotifyEvent:
						cur.Call(func(
							wins WindowsMap,
						) {
							// update LastFocus
							if w, ok := wins[ev.Event]; ok {
								w.LastFocus = time.Now()
							}
							// focus pointer root
							ce(xproto.SetInputFocusChecked(
								conn, 0, xproto.InputFocusPointerRoot, 0,
							).Check())
						})

					case xproto.ButtonPressEvent:
						cur.Call(func(
							wins WindowsMap,
							stack StackWindows,
							conn *xgb.Conn,
						) {
							// update LastKey
							if w, ok := wins[ev.Event]; ok {
								w.LastKey = time.Now()
							}
							// stack
							stack(byInteractTime)
							// allow events
							ce(xproto.AllowEventsChecked(conn, xproto.AllowReplayPointer, ev.Time).Check())
							ce(xproto.AllowEventsChecked(conn, xproto.AllowReplayKeyboard, ev.Time).Check())
						})

					case xproto.ButtonReleaseEvent:
						cur.Call(func(
							conn *xgb.Conn,
						) {
							// allow events
							ce(xproto.AllowEventsChecked(conn, xproto.AllowReplayPointer, ev.Time).Check())
							ce(xproto.AllowEventsChecked(conn, xproto.AllowReplayKeyboard, ev.Time).Check())
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
