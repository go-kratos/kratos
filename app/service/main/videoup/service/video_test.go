package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ObtainCid(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("ObtainCid", t, WithService(func(s *Service) {
		_, err := svr.ObtainCid(c, "")
		So(err, ShouldBeNil)
	}))
}

func TestService_FindCidByFn(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("FindCidByFn", t, WithService(func(s *Service) {
		_, err := svr.FindCidByFn(c, "")
		So(err, ShouldBeNil)
	}))
}

func BenchmarkDo(b *testing.B) {
	var (
		c = context.TODO()
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = svr.ObtainCid(c, "22222222222")
		}
	})
}
