package secure

import (
	"context"
	"testing"
	"time"

	model "go-common/app/service/main/secure/model"
)

var s *Service

func TestSecure(t *testing.T) {
	s = New(nil)
	time.Sleep(1000 * time.Second)
	testStatus(t)
	testExpectionLoc(t)
}

// TestStatus test status rpc.
func testStatus(t *testing.T) {
	if res, err := s.Status(context.TODO(), &model.ArgSecure{Mid: 1, UUID: "2"}); err != nil {
		t.Errorf("Service: Status err: %v", err)
	} else {
		t.Logf("Service: Status res: %+v", res)
	}
}

func testExpectionLoc(t *testing.T) {
	if res, err := s.ExpectionLoc(context.TODO(), &model.ArgSecure{Mid: 1, UUID: "2"}); err != nil {
		t.Errorf("Service : Expection err:%v", err)
	} else {
		t.Logf("Service: Status res: %+v", res)
	}
}
