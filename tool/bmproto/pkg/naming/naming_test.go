package naming

import "testing"

func TestGetGoImportPathForPb(t *testing.T) {
	p, err := GetGoImportPathForPb("naming.go", "go-common", "go-common")
	if err != nil {
		t.Fatalf("err is not nil %v", err)
	} else {
		if p != "github.com/bilibili/kratos/tool/bmproto/pkg/naming" {
			t.Fatalf("path is not correct" + p)
		}
	}
}
