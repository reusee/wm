package main

type Start func()

func (_ Def) Start() Start {
	return func() {

		pt("start\n")

	}
}
