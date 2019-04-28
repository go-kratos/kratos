// Package generator defines an interface for code generators to implement.
//
// To use this package, you'll implement the "Package" and "Generator"
// interfaces; you'll call NewContext to load up the types you want to work
// with, and then you'll call one or more of the Execute methods. See the
// interface definitions for explanations. All output will have gofmt called on
// it automatically, so you do not need to worry about generating correct
// indentation.
//
// This package also exposes SnippetWriter. SnippetWriter reduces to a minimum
// the boilerplate involved in setting up a template from go's text/template
// package. Additionally, all naming systems in the Context will be added as
// functions to the parsed template, so that they can be called directly from
// your templates!
package generator // import "go-common/app/tool/gengo/generator"
