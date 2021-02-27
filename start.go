package main

type Start func()

func (_ Def) Start(
	becomeWM BecomeWM,
	internAtoms InternAtoms,
	setupAllDesktop SetupAllDesktop,
) Start {
	return func() {

		pt("starting\n")

		becomeWM()
		internAtoms()
		setupAllDesktop()

		pt("started\n")

	}
}
