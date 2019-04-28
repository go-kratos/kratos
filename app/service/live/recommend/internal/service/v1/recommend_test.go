package v1

import (
	"context"
	"flag"
	"testing"

	"go-common/app/service/live/recommend/api/grpc/v1"
	"go-common/app/service/live/recommend/internal/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *RecommendService
)

func init() {
	flag.Set("conf", "../../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = NewRecommendService(conf.Conf)
}

// go test  -test.v -test.run TestRecommend_RandomRecsByUser
func TestRecommend_RandomRecsByUser(t *testing.T) {
	Convey("TestRecommend_RandomRecsByUser", t, func() {
		res, err := s.RandomRecsByUser(context.TODO(), &v1.GetRandomRecReq{Uid: 4158272, Count: 3})
		t.Logf("%v msg", res)
		if err != nil {
			t.Logf("err=%+v", err)
		}
		So(err, ShouldBeNil)
	})
}
