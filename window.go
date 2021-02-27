package main

import (
	"sort"
	"time"

	"github.com/jezek/xgb/xproto"
)

type Window struct {
	XID       xproto.Window
	LastFocus time.Time
}

type WindowsSorter struct {
	windows []*Window
	less    func(i, j int) bool
	swap    func(i, j int)
}

var _ sort.Interface = WindowsSorter{}

func (s WindowsSorter) Len() int {
	return len(s.windows)
}

func (s WindowsSorter) Less(i, j int) bool {
	return s.less(i, j)
}

func (s WindowsSorter) Swap(i, j int) {
	s.swap(i, j)
	s.windows[i], s.windows[j] = s.windows[j], s.windows[i]
}
