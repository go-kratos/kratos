package ecode

import (
	"testing"
)

func TestNew(t *testing.T) {
	defer func() {
		errStr := recover()
		if errStr != "ecode: 1 already exist" {
			t.Logf("New duplicate ecode should cause panic")
			t.FailNow()
		}
	}()
	var _ error = New(1)
	var _ error = New(2)
	var _ error = New(1)
}

func TestErrMessage(t *testing.T) {
	e1 := New(3)
	if e1.Error() != "3" {
		t.Logf("ecode message should be `3`")
		t.FailNow()
	}
	if e1.Message() != "3" {
		t.Logf("unregistered ecode message should be ecode number")
		t.FailNow()
	}
	Register(map[int]string{3: "testErr"})
	if e1.Message() != "testErr" {
		t.Logf("registered ecode message should be `testErr`")
		t.FailNow()
	}
}

func TestCause(t *testing.T) {
	e1 := New(4)
	var err error = e1
	e2 := Cause(err)
	if e2.Code() != 4 {
		t.Logf("parsed error code should be 4")
		t.FailNow()
	}
}

func TestInt(t *testing.T) {
	e1 := Int(1)
	if e1.Code() != 1 {
		t.Logf("int parsed error code should be 1")
		t.FailNow()
	}
	if e1.Error() != "1" || e1.Message() != "1" {
		t.Logf("int parsed error string should be `1`")
		t.FailNow()
	}
}

func TestString(t *testing.T) {
	eStr := String("123")
	if eStr.Code() != 123 {
		t.Logf("string parsed error code should be 123")
		t.FailNow()
	}
	if eStr.Error() != "123" || eStr.Message() != "123" {
		t.Logf("string parsed error string should be `123`")
		t.FailNow()
	}
	eStr = String("test")
	if eStr.Code() != -500 {
		t.Logf("invalid string parsed error code should be -500")
		t.FailNow()
	}
	if eStr.Error() != "-500" || eStr.Message() != "-500" {
		t.Logf("invalid string parsed error string should be `-500`")
		t.FailNow()
	}
}
