package client

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/up/model"
)

func TestRpcClient(t *testing.T) {
	s := New(nil)
	time.Sleep(1 * time.Second)
	testInfo(t, s)
	testSpecial(t, s)
	testUpStatBase(t, s)
	testUpSwitch(t, s)
}

func testInfo(t *testing.T, s *Service) {
	arg := model.ArgInfo{
		Mid:  2089809,
		From: 1,
	}
	t.Log(s.Info(context.TODO(), &arg))
}

func testSpecial(t *testing.T, s *Service) {
	arg := model.ArgSpecial{
		GroupID: 1,
	}
	t.Log(s.Special(context.TODO(), &arg))
}

func testUpStatBase(t *testing.T, s *Service) {
	arg := model.ArgMidWithDate{
		Mid: 12345,
	}
	t.Log(s.UpStatBase(context.TODO(), &arg))
}

func testUpSwitch(t *testing.T, s *Service) {
	arg := model.ArgUpSwitch{
		Mid:  1,
		From: 0,
	}
	t.Log(s.UpSwitch(context.TODO(), &arg))
}
