package log

import "strconv"

// VerboseKey is logger verbose key.
const VerboseKey = "verbose"

// Verbose .
type Verbose int

func (v Verbose) String() string {
	return strconv.Itoa(int(v))
}

// Enabled .
func (v Verbose) Enabled(n Verbose) bool {
	return v > n
}
