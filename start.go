package main

type Start func()

func (_ Def) Start(
	becomeWM BecomeWM,
	internAtoms InternAtoms,
	setupDesktop SetupDesktop,
) Start {
	return func() {

		pt("starting\n")

		becomeWM()
		internAtoms()
		setupDesktop()

		pt("started\n")

	}
}
