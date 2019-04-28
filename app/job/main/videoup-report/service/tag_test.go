package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestService_upbind(t *testing.T) {
	Convey("upbind", t, func() {
		a, _ := s.arc.ArchiveByAid(context.Background(), 10110255)
		err := s.upBindTag(context.Background(), a.Mid, a.ID, a.Tag, a.TypeID)
		t.Logf("err(%+v)", err)
	})
}
