package main

import "os/exec"

type Start func()

func (_ Def) Start(
	becomeWM BecomeWM,
	internAtoms InternAtoms,
	setupAllDesktop SetupAllDesktop,
	setupKeyboard SetupKeyboard,
	manageExisting ManageExistingWindows,
	setupEventHandler SetupEventHandler,
) Start {
	return func() {

		pt("starting\n")

		becomeWM()
		internAtoms()
		setupAllDesktop()
		setupKeyboard()
		manageExisting()
		setupEventHandler()

		ce(exec.Command("terminal").Start())
		ce(exec.Command("gedit").Start())
		ce(exec.Command("gnome-calculator").Start())

		pt("started\n")

	}
}
