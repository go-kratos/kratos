package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/interface/live/app-blink/api/http/v1"
	"go-common/library/net/metadata"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/interface/live/app-blink/conf"
)

var topic *TopicService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	topic = NewTopicService(conf.Conf)
}

// group=qa01 DEPLOY_ENV=uat go test -run TestGetTopicList
func TestGetTopicList(t *testing.T) {
	Convey("TestGetTopicList", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(16299525),
		})
		res, err := topic.GetTopicList(ctx, &v1.GetTopicListReq{
			Platform: "ios",
		})
		fmt.Println(1111, res, err, 22222)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

func TestCheckTopic(t *testing.T) {
	Convey("TestCheckTopic", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(16299525),
		})
		res, err := topic.CheckTopic(ctx, &v1.CheckTopicReq{
			Platform: "ios",
			Topic:    "我是习",
		})

		fmt.Println(1111, res, err, 22222)
		t.Logf("%v,%s", res, err)
	})
}
