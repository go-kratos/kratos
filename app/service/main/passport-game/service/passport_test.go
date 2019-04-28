package service

import (
	"encoding/json"
	"testing"
)

func TestServiceApp(t *testing.T) {
	once.Do(startService)
	if app, ok := s.appMap[_gameAppKey]; !ok {
		t.Errorf("res is not correct, expect game app exists but nil")
		t.FailNow()
	} else {
		str, _ := json.Marshal(app)
		t.Logf("res: %s", str)
	}
}
