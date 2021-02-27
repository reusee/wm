package main

import "os"

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
		var err error
		defer he(&err, func(prev error) error {
			pt("%s\n", prev.Error())
			os.Exit(-1)
			return prev
		})
		start()
		select {}
	})

}
