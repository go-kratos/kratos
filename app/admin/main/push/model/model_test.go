package model

import "testing"

func Test_Funcs(t *testing.T) {
	s := RandomString(32)
	t.Logf("random string: %s", s)
}
