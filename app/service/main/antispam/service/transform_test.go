package service

import "testing"

func TestToDaoArea(t *testing.T) {
	ToDaoArea("reply")
}

func TestToModelArea(t *testing.T) {
	ToModelArea(1)
}

func TestToDaoState(t *testing.T) {
	ToDaoState("default")
}

func TestToModelState(t *testing.T) {
	ToModelState(1)
}
