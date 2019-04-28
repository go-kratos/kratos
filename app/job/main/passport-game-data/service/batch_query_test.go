package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_batchQueryCloudByMid(t *testing.T) {
	once.Do(startService)

	Convey("batch query cloud by mid", t, func() {
		mids := []int64{88888970, 88894849}
		res, miss := s.batchQueryCloudByMid(context.TODO(), mids, 1)
		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)

		So(len(miss), ShouldEqual, 0)
		So(len(res), ShouldEqual, 2)
	})
}

func TestService_batchQueryCloudByMidNonMiss(t *testing.T) {
	once.Do(startService)

	Convey("batch query cloud by mid", t, func() {
		mids := []int64{88888970, 88894849}
		res := s.batchQueryCloudNonMiss(context.TODO(), mids, 1001, 1)
		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)

		So(len(res), ShouldEqual, 2)
	})
}
