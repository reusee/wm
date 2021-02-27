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
	throw  = e4.Throw
	is     = errors.Is
	as     = errors.As
)
