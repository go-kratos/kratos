package main

// go build -ldflags "-X main.Version=x.y.yz"
var (
	// Version is the version of the compiled software.
	Version string
	// Branch is current branch name the code is built off
	Branch string
	// Revision is the short commit hash of source tree
	Revision string
	// BuildDate is the date when the binary was built.
	BuildDate string
)
