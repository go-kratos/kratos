package spy

import (
	"context"
	"testing"
	"time"

	model "go-common/app/service/main/spy/model"
)

func TestSpy(t *testing.T) {
	s := New(nil)
	time.Sleep(2 * time.Second)
	testUpdateEventScore(t, s)
	testUpdateBaseScore(t, s)
	testUserScore(t, s)
	testHandleEvent(t, s)
	testReBuildPortrait(t, s)
}

func testUpdateEventScore(t *testing.T, s *Service) {
	t.Log(s.UpdateEventScore(context.TODO(), &model.ArgReset{Mid: 23333, Operator: "admin test"}))
}

func testUpdateBaseScore(t *testing.T, s *Service) {
	t.Log(s.UpdateBaseScore(context.TODO(), &model.ArgReset{Mid: 23333, Operator: "admin test"}))
}

func testUserScore(t *testing.T, s *Service) {
	t.Log(s.UserScore(context.TODO(), &model.ArgUserScore{Mid: 23333, IP: "127.0.0.1"}))
}

func testHandleEvent(t *testing.T, s *Service) {
	t.Log(s.HandleEvent(context.TODO(), &model.ArgHandleEvent{
		IP:        "127.0.0.1",
		Service:   "spy_service",
		Event:     "bind_mail_only",
		ActiveMid: 23333,
		TargetMid: 23333,
		Effect:    "",
		RiskLevel: 1,
	}))
}

func testReBuildPortrait(t *testing.T, s *Service) {
	t.Log(s.ReBuildPortrait(context.TODO(), &model.ArgReBuild{Mid: 23333}))
}
