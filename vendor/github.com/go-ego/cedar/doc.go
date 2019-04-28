// Package cedar implements double-array trie.
//
// It is a golang port of cedar (http://www.tkl.iis.u-tokyo.ac.jp/~ynaga/cedar)
// which is written in C++ by Naoki Yoshinaga.
// Currently cedar-go implements the `reduced` verion of cedar.
// This package is not thread safe if there is one goroutine doing
// insertions or deletions.
//
// Note
//
// key must be `[]byte` without zero items,
// while value must be integer in the range [0, 2<<63-2] or
// [0, 2<<31-2] depends on the platform.
package cedar
