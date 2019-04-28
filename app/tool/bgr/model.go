package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type lint struct {
	s  *script
	fn func(curDir string, f *ast.File, node ast.Node) bool
}

type script struct {
	dir string
	ts  []string // type slice
	v   string
	l   string
	d   string
}

func (s script) String() string {
	return fmt.Sprintf("script path: %s, type: %s, value: %s, level: %s", s.dir, strings.Join(s.ts, "."), s.v, s.l)
}
