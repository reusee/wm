package main

func main() {

	scope := NewScope()
	scope = scope.Sub(func() *Scope {
		return &scope
	})

	update := Update(func(decls ...any) Scope {
		scope = scope.Sub(decls...)
		return scope
	})
	scope = scope.Sub(&update)

	scope.Call(func(
		start Start,
	) {
		start()
	})

}
