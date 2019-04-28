package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Region(t *testing.T) {
	once.Do(startService)
	Convey("regions", t, func() {
		m := s.Regions(context.TODO())
		So(len(m), ShouldBeGreaterThan, 0)
		res, _ := json.Marshal(m)
		t.Logf("res: %s", res)
	})
}
