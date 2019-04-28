package main

type fileInfo struct {
	path, rel, name, ext string

	// packageName is the Go package name of a .go file, without the
	// "_test" suffix if it was present. It is empty for non-Go files.
	packageName string

	// importPath is the canonical import path for this file's package.
	// This may be read from a package comment (in Go) or a go_package
	// option (in proto). This field is empty for files that don't specify
	// an import path.
	importPath string

	// category is the type of file, based on extension.
	category extCategory

	// isTest is true if the file stem (the part before the extension)
	// ends with "_test.go". This is never true for non-Go files.
	isTest bool

	// isXTest is true for test Go files whose declared package name ends
	// with "_test".
	isXTest bool

	// imports is a list of packages imported by a file. It does not include
	// "C" or anything from the standard library.
	imports []string

	// isCgo is true for .go files that import "C".
	isCgo bool

	// goos and goarch contain the OS and architecture suffixes in the filename,
	// if they were present.
	goos, goarch string

	// tags is a list of build tag lines. Each entry is the trimmed text of
	// a line after a "+build" prefix.
	tags []tagLine

	// copts and clinkopts contain flags that are part of CFLAGS, CPPFLAGS,
	// CXXFLAGS, and LDFLAGS directives in cgo comments.
	copts, clinkopts []taggedOpts

	// hasServices indicates whether a .proto file has service definitions.
	hasServices bool
}

// tagLine represents the space-separated disjunction of build tag groups
// in a line comment.
type tagLine []tagGroup

// tagGroup represents a comma-separated conjuction of build tags.
type tagGroup []string

// extCategory indicates how a file should be treated, based on extension.
type extCategory int

// taggedOpts a list of compile or link options which should only be applied
// if the given set of build tags are satisfied. These options have already
// been tokenized using the same algorithm that "go build" uses, then joined
// with OptSeparator.
type taggedOpts struct {
	tags tagLine
	opts string
}

func fileNameInfo(dir, rel, name string) fileInfo {
	return fileInfo{}
}
