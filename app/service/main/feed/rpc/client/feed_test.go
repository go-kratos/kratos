package client

import (
	"context"
	"testing"
	"time"

	model "go-common/app/service/main/feed/model"
)

func TestFeed(t *testing.T) {
	s := New(nil)
	time.Sleep(1 * time.Second)
	testAppFeed(t, s)
	testWebFeed(t, s)
	testAddArc(t, s)
	testDelArc(t, s)
	testPurgeFeedCache(t, s)
	testFold(t, s)
}

func testAppFeed(t *testing.T, s *Service) {
	if res, err := s.AppFeed(context.TODO(), &model.ArgFeed{Mid: 27515256, Pn: 1, Ps: 20}); err != nil {
		t.Errorf("Service: AppFeed err: %v", err)
	} else {
		t.Logf("Service: AppFeed: %v", res)
	}
}

func testWebFeed(t *testing.T, s *Service) {
	if res, err := s.WebFeed(context.TODO(), &model.ArgFeed{Mid: 27515256, Pn: 1, Ps: 20}); err != nil {
		t.Errorf("Service: WebFeed err: %v", err)
	} else {
		t.Logf("Service: WebFeed: %v", res)
	}
}

func testArchiveFeed(t *testing.T, s *Service) {
	if res, err := s.ArchiveFeed(context.TODO(), &model.ArgFeed{Mid: 27515256, Pn: 1, Ps: 20}); err != nil {
		t.Errorf("Service: ArchiveFeed err: %v", err)
	} else {
		t.Logf("Service: ArchiveFeed: %v", res)
	}
}

func testBangumiFeed(t *testing.T, s *Service) {
	if res, err := s.BangumiFeed(context.TODO(), &model.ArgFeed{Mid: 27515256, Pn: 1, Ps: 20}); err != nil {
		t.Errorf("Service: BangumiFeed err: %v", err)
	} else {
		t.Logf("Service: BangumiFeed: %v", res)
	}
}

func testAddArc(t *testing.T, s *Service) {
	if err := s.AddArc(context.TODO(), &model.ArgArc{Aid: 1}); err != nil {
		t.Errorf("Service: AddArc err: %v", err)
	}
}

func testDelArc(t *testing.T, s *Service) {
	if err := s.DelArc(context.TODO(), &model.ArgAidMid{Aid: 1}); err != nil {
		t.Errorf("Service: DelArc err: %v", err)
	}
}

func testPurgeFeedCache(t *testing.T, s *Service) {
	if err := s.PurgeFeedCache(context.TODO(), &model.ArgMid{Mid: 27515256}); err != nil {
		t.Errorf("Service: PurgeFeedCache err: %v", err)
	}
}

func testFold(t *testing.T, s *Service) {
	if res, err := s.Fold(context.TODO(), &model.ArgFold{Aid: 1, Mid: 27515256}); err != nil {
		t.Errorf("Service: Fold err: %v", err)
	} else {
		t.Logf("Service: Fold: %v", res)
	}
}

func testAppUnreadCount(t *testing.T, s *Service) {
	if res, err := s.AppUnreadCount(context.TODO(), &model.ArgUnreadCount{Mid: 27515256, WithoutBangumi: false}); err != nil {
		t.Errorf("Service: UnreadCount err: %v", err)
	} else {
		t.Logf("Service: UnreadCount: %v", res)
	}
}

func testWebUnreadCount(t *testing.T, s *Service) {
	if res, err := s.WebUnreadCount(context.TODO(), &model.ArgMid{Mid: 27515256}); err != nil {
		t.Errorf("Service: UnreadCount err: %v", err)
	} else {
		t.Logf("Service: UnreadCount: %v", res)
	}
}

func testChangeArcUpper(t *testing.T, s *Service) {
	if err := s.ChangeArcUpper(context.TODO(), &model.ArgChangeUpper{Aid: 1, OldMid: 1, NewMid: 2, RealIP: "127.0.0.1"}); err != nil {
		t.Errorf("Service: ChangeArcUpper err: %v", err)
	} else {
		t.Logf("Service: ChangeArcUpper ok")
	}
}
