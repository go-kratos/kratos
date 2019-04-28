package tomllint

import (
	"bytes"
	"strings"
	"testing"
)

func TestSyntaxError(t *testing.T) {
	lint := &tomllint{}
	r := bytes.NewBufferString(synataxerrordata)
	lintErr := lint.Lint(r)
	if lintErr == nil {
		t.Fatalf("expect lintErr != nil")
	}
	if lintErr[0].Line == -1 {
		t.Errorf("expect get line number")
	}
}

func TestTomlLintOK(t *testing.T) {
	lint := &tomllint{}
	r := bytes.NewBufferString(normaldata)
	lintErr := lint.Lint(r)
	if lintErr != nil {
		t.Errorf("error %v", lintErr)
	}
}

func TestNoCommon(t *testing.T) {
	lint := &tomllint{}
	r := bytes.NewBufferString(nocommondata)
	lintErr := lint.Lint(r)
	if lintErr == nil {
		t.Fatalf("expect lintErr != nil")
	}
	message := lintErr.Error()
	if !strings.Contains(message, "Common") {
		t.Errorf("expect error contains common")
	}
}

func TestNoIdentify(t *testing.T) {
	lint := &tomllint{}
	r := bytes.NewBufferString(noidentify)
	lintErr := lint.Lint(r)
	if lintErr == nil {
		t.Fatalf("expect lintErr != nil")
	}
	message := lintErr.Error()
	if !strings.Contains(message, "Identify") {
		t.Errorf("expect error Identify common")
	}
}

func TestNoApp(t *testing.T) {
	lint := &tomllint{}
	r := bytes.NewBufferString(noapp)
	lintErr := lint.Lint(r)
	if lintErr == nil {
		t.Fatalf("expect lintErr != nil")
	}
	message := lintErr.Error()
	if !strings.Contains(message, "App") {
		t.Errorf("expect error App common")
	}
}
