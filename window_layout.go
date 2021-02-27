package main

type Relayout func()

func (_ Def) Relayout(
	stack StackWindows,
	setGeometries SetGeometries,
) Relayout {
	return func() {
		windows := stack()
		setGeometries(windows)
	}
}
