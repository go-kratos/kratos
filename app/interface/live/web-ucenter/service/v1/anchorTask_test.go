package v1

import (
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	api "go-common/app/interface/live/web-ucenter/api/http/v1"
	"go-common/app/interface/live/web-ucenter/conf"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"testing"
)

var (
	AnchorTask *AnchorTaskService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	AnchorTask = NewAnchorTaskService(conf.Conf)
}

// go test  -test.v -test.run TestServiceAllowanceList
func TestMyReward(t *testing.T) {
	Convey("TestMyReward", t, func() {

		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := AnchorTask.MyReward(ctx, &api.AnchorTaskMyRewardReq{
			Page: 1,
		})
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}
