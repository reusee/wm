package main

type Start func()

func (_ Def) Start(
	becomeWM BecomeWM,
	internAtoms InternAtoms,
	setupAllDesktop SetupAllDesktop,
	setupKeyboard SetupKeyboard,
) Start {
	return func() {

		pt("starting\n")

		becomeWM()
		internAtoms()
		setupAllDesktop()
		setupKeyboard()

		pt("started\n")

	}
}
