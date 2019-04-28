package config

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/infra/config/model"
)

func TestConf(t *testing.T) {
	s := New2(nil)
	time.Sleep(1 * time.Second)

	// coin
	testPush(t, s)
	testSetToken(t, s)
	testHosts(t, s)
	testClearHost(t, s)
}

func testPush(t *testing.T, s *Service2) {
	arg := &model.ArgConf{
		App:      "zjx_test",
		BuildVer: "1_0_0_0",
		Ver:      113,
		Env:      "2",
	}
	if err := s.Push(context.TODO(), arg); err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
func testSetToken(t *testing.T, s *Service2) {
	arg := &model.ArgToken{
		App:   "zjx_test",
		Token: "123",
		Env:   "2",
	}
	if err := s.SetToken(context.TODO(), arg); err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func testHosts(t *testing.T, s *Service2) {
	if hosts, err := s.Hosts(context.TODO(), "testApp4890934756659"); err != nil {
		t.Log(err)
		t.FailNow()
	} else {
		t.Log(len(hosts))
	}
}

func testClearHost(t *testing.T, s *Service2) {
	if err := s.ClearHost(context.TODO(), "testApp4890934756659"); err != nil {
		t.Log(err)
		t.FailNow()
	}
}
