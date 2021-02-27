package main

type Start func()

func (_ Def) Start(
	becomeWM BecomeWM,
) Start {
	return func() {

		pt("starting\n")

		becomeWM()

		pt("started\n")

	}
}
