package project

import "testing"

func TestNewProjInfo(t *testing.T) {
	proj, err := NewProjInfo("../api/stub.txt", "go-common", "go-common")
	if err != nil {
		t.Fatalf("init project fail %+v", err)
	} else {
		if proj.Name != "pkg" {
			t.Fatalf("name fail, expect project, got %s", proj.Name)
		}
		if proj.Department != "bmproto" {
			t.Fatalf("department fail, expect pkg, got %s", proj.Department)
		}
		if proj.Typ != "tool" {
			t.Fatalf("typ fail, expect tool, got %s", proj.Typ)
		}
		if proj.PathRefToProj != ".." {
			t.Fatalf("pathRefToProj fail, expect .., got %s", proj.PathRefToProj)
		}
		if proj.ImportPath != "kratos/tool/bmproto/pkg" {
			t.Fatalf("pathRefToProj fail, expect kratos/tool/bmproto/pkg, got %s", proj.ImportPath)
		}
		if proj.HasInternalPkg != false {
			t.Fatalf("pathRefToProj fail, expect false got %v", proj.HasInternalPkg)
		}
	}
}
