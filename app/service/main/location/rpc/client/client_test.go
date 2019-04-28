package client

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/location/model"
)

const (
	_aid = 16428
	_gid = 317
	_mid = 0
	_ip  = "139.214.144.59"
	_cip = "127.0.0.1"
)

func TestLocation(t *testing.T) {
	s := New(nil)
	time.Sleep(1 * time.Second)
	testArchive(t, s)
	testArchive2(t, s)
	testGroup(t, s)
	testAuthPIDs(t, s)
	testInfo(t, s)
	testInfos(t, s)
	testInfoComplete(t, s)
	testInfosComplete(t, s)
}

func testArchive(t *testing.T, s *Service) {
	if res, err := s.Archive(context.TODO(), &model.Archive{Aid: _aid, Mid: _mid, IP: _ip, CIP: _cip}); err != nil {
		t.Errorf("Service: archive err: %v", err)
	} else {
		t.Logf("Service: archive res: %v", res)
	}
}

func testArchive2(t *testing.T, s *Service) {
	if res, err := s.Archive2(context.TODO(), &model.Archive{Aid: _aid, Mid: _mid, IP: _ip, CIP: _cip}); err != nil {
		t.Errorf("Service: archive2 err: %v", err)
	} else {
		t.Logf("Service: archive2 res: %v", res)
	}
}

func testGroup(t *testing.T, s *Service) {
	if res, err := s.Group(context.TODO(), &model.Group{Gid: _gid, Mid: _mid, IP: _ip, CIP: _cip}); err != nil {
		t.Errorf("Service: group err: %v", err)
	} else {
		t.Logf("Service: group res: %v", res)
	}
}

func testAuthPIDs(t *testing.T, s *Service) {
	if res, err := s.AuthPIDs(context.TODO(), &model.ArgPids{IP: _ip, Pids: "2,1163,86,87", CIP: _cip}); err != nil {
		t.Errorf("Service: AuthPIDs err: %v", err)
	} else {
		t.Logf("Service: AuthPIDs res: %v", res)
	}
}

func testInfo(t *testing.T, s *Service) {
	if res, err := s.Info(context.TODO(), &model.ArgIP{IP: _ip}); err != nil {
		t.Errorf("Service: info err: %v", err)
	} else {
		t.Logf("Service: info res: %v", res)
	}
}

func testInfos(t *testing.T, s *Service) {
	if res, err := s.Infos(context.TODO(), []string{"61.216.166.156", "211.139.80.6"}); err != nil {
		t.Errorf("Service: infos err: %v", err)
	} else {
		t.Logf("Service: infos res: %v", res)
	}
}

func testInfoComplete(t *testing.T, s *Service) {
	if res, err := s.InfoComplete(context.TODO(), &model.ArgIP{IP: _ip}); err != nil {
		t.Errorf("Service: infoComplete err: %v", err)
	} else {
		t.Logf("Service: infoComplete res: %v", res)
	}
}

func testInfosComplete(t *testing.T, s *Service) {
	if res, err := s.InfosComplete(context.TODO(), []string{"61.216.166.156", "211.139.80.6"}); err != nil {
		t.Errorf("Service: infosComplete err: %v", err)
	} else {
		t.Logf("Service: infosComplete res: %v", res)
	}
}
