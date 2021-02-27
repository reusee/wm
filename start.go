package main

type Start func()

func (_ Def) Start(
	rootWin XRootWindow,
) Start {
	return func() {

		pt("starting\n")

		pt("root win: %v\n", rootWin)

		pt("started\n")
	}
}
