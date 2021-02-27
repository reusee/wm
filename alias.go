package main

import (
	"errors"
	"fmt"

	"github.com/reusee/dscope"
	"github.com/reusee/e4"
)

type (
	Scope = dscope.Scope

	any = interface{}
)

var (
	pt     = fmt.Printf
	ce, he = e4.Check, e4.Handle
	is     = errors.Is
	as     = errors.As
)
