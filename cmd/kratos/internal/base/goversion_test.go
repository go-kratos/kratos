package base

import "testing"

func TestGetGoVersion(t *testing.T) {
	version, err := GetGoVersion()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(version)
}
