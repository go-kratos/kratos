package service

import (
	"context"
	"time"

	spym "go-common/app/service/main/spy/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testSpyMid int64 = 111
)

// go test -test.v -test.run TestPutSpyScore
func TestPutSpyScore(t *testing.T) {
	Convey("Test_PutSpyScore put no reason msg", t, WithService(func(s *Service) {
		So(s.PutSpyScore(context.TODO(), &spym.ScoreChange{
			Mid:        testSpyMid,
			Score:      80,
			EventScore: 80,
			BaseScore:  80,
			TS:         time.Now().Unix(),
		}), ShouldBeNil)
	}))
	Convey("Test_PutSpyScore put coin-service reason msg", t, WithService(func(s *Service) {
		So(s.PutSpyScore(context.TODO(), &spym.ScoreChange{
			Mid:        testSpyMid,
			Score:      70,
			EventScore: 80,
			BaseScore:  80,
			TS:         time.Now().Unix(),
			Reason:     "coin-service",
			RiskLevel:  3,
		}), ShouldBeNil)
	}))
	Convey("Test_PutSpyScore put coin-service high risk reason msg", t, WithService(func(s *Service) {
		So(s.PutSpyScore(context.TODO(), &spym.ScoreChange{
			Mid:        testSpyMid,
			Score:      8,
			EventScore: 80,
			BaseScore:  80,
			TS:         time.Now().Unix(),
			Reason:     "coin-service",
			RiskLevel:  7,
		}), ShouldBeNil)
	}))
}
