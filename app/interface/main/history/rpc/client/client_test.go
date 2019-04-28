package client

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/history/model"
)

func TestHistory(t *testing.T) {
	s := New(nil)
	time.Sleep(1 * time.Second)
	testProgress(t, s)
	testAdd(t, s)
}

// testProgress test progress rpc.
func testProgress(t *testing.T, s *Service) {
	if res, err := s.Progress(context.TODO(), &model.ArgPro{Mid: 88888966, Aids: []int64{5463286}}); err != nil {
		t.Errorf("Service: Progress err: %v", err)
	} else {
		t.Logf("Service: zone res: %+v", res)
	}
}

// testAdd test add rpc .
func testAdd(t *testing.T, s *Service) {
	h := &model.History{Mid: 88888966, Aid: 5463286, TP: -1, Unix: 1494217922}
	if err := s.Add(context.TODO(), &model.ArgHistory{Mid: 88888966, History: h}); err != nil {
		t.Errorf("Service: Add err: %v", err)
	}
}
