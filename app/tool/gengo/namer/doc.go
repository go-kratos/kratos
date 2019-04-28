// Package namer has support for making different type naming systems.
//
// This is because sometimes you want to refer to the literal type, sometimes
// you want to make a name for the thing you're generating, and you want to
// make the name based on the type. For example, if you have `type foo string`,
// you want to be able to generate something like `func FooPrinter(f *foo) {
// Print(string(*f)) }`; that is, you want to refer to a public name, a literal
// name, and the underlying literal name.
//
// This package supports the idea of a "Namer" and a set of "NameSystems" to
// support these use cases.
//
// Additionally, a "RawNamer" can optionally keep track of what needs to be
// imported.
package namer // import "go-common/app/tool/gengo/namer"
