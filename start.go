package main

type Start func()

func (_ Def) Start(
	becomeWM BecomeWM,
	internAtoms InternAtoms,
) Start {
	return func() {

		pt("starting\n")

		becomeWM()
		internAtoms()

		pt("started\n")

	}
}
