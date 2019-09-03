package ecode

import (
	"testing"
)

func TestLocalizedError(t *testing.T) {
	e1 := New(3)
	if e1.Error() != "3" {
		t.Logf("ecode message should be `3`")
		t.FailNow()
	}
	if e1.Message() != "3" {
		t.Logf("unregistered ecode message should be ecode number")
		t.FailNow()
	}
	codes := map[int]map[string]string{
		3: {
			"default": "testErr",
			"en-US":   "test_Err",
			"zh-TW":   "我是測試，這是繁體字",
		},
	}
	Register(codes)
	codeEnUS := LocalizedError(e1, []string{LangEnUS, LangZhCN})
	if codeEnUS.Message() != "test_Err" {
		t.Logf("registered ecode message should be `testErr`")
		t.FailNow()
	}
	codeZhCN := LocalizedError(e1, []string{LangZhCN, LangZhHK})
	if codeZhCN.Message() != "testErr" {
		t.Logf("registered ecode message should be `testErr`")
		t.FailNow()
	}
	codeZhHK := LocalizedError(e1, []string{LangZhHK, LangEnUS})
	if codeZhHK.Message() != "我是測試，這是繁體字" {
		t.Logf("registered ecode message should be `testErr`")
		t.FailNow()
	}
	codeZhTW := LocalizedError(e1, []string{LangZhTW})
	if codeZhTW.Message() != "我是測試，這是繁體字" {
		t.Logf("registered ecode message should be `testErr`")
		t.FailNow()
	}
	codeJaJP := LocalizedError(e1, []string{LangJaJP})
	if codeJaJP.Message() != "testErr" {
		t.Logf("registered ecode message should be `testErr`")
		t.FailNow()
	}
}
