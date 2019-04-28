package agent

import (
	"testing"
)

func TestRunnerStart(t *testing.T) {
	stdout, stderr, err := RunnerStart()
	if err != nil {
		t.Error(stdout, stderr)
	}
}

func TestRunnerRegister(t *testing.T) {
	stdout, stderr, err := RunnerRegister("http://gitlab.bilibili.co/", "pxZPKWk1JQHLHNzYyj5p", "mac-test-1")
	if err != nil {
		t.Error(stdout, stderr)
	}
}

func TestRunnerUnRegister(t *testing.T) {
	stdout, stderr, err := RunnerUnRegister("mac-test-1")
	if err != nil {
		t.Error(stdout, stderr)
	}
}

func TestRunnerUnRegisterAll(t *testing.T) {
	stdout, stderr, err := RunnerUnRegisterAll()
	if err != nil {
		t.Error(stdout, stderr)
	}
}
