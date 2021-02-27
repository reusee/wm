// +build ignore

package main

import (
	"bytes"
	"fmt"
	"github.com/reusee/e4"
	"go/format"
	"io/ioutil"
	"regexp"
)

var (
	ce = e4.Check
)

func main() {
	buf := new(bytes.Buffer)
	w := func(format string, args ...interface{}) {
		fmt.Fprintf(buf, format, args...)
	}
	w(`package main

const (
`)
	content, err := ioutil.ReadFile("/usr/include/X11/keysymdef.h")
	ce(err)
	defPattern := regexp.MustCompile(`(?m:^#define (?:XK_)([^\s]+)\s+(\S+))`)
	matches := defPattern.FindAllSubmatch(content, -1)
	for _, group := range matches {
		w("\tKey_%s KeySym = %s\n", group[1], group[2])
	}
	w(")\n")
	src, err := format.Source(buf.Bytes())
	ce(err)
	err = ioutil.WriteFile("x_key_sym.go", src, 0644)
	ce(err)
}
