package main

type SetGeometries func(
	windows []*Window,
)

func (_ Def) SetGeometries() SetGeometries {
	return func(
		windows []*Window,
	) {

		//TODO

	}
}
