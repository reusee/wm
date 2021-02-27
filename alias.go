package main

import (
	"fmt"

	"github.com/reusee/dscope"
)

type (
	Scope = dscope.Scope

	any = interface{}
)

var (
	pt = fmt.Printf
)
