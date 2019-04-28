package price

import (
	"context"
	"os"
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	p   *Price
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	var err error
	p = New()
	var pcm = map[int64]map[int8][]*model.VipPriceConfig{
		1: {
			1: {{ID: 1}},
		},
	}
	if err = p.SetPriceConfig(ctx, pcm); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestPrice(t *testing.T) {
	Convey("TestPrice", t, func() {
		i, ok := p.GetPriceConfig(1)
		So(i, ShouldNotBeNil)
		So(ok, ShouldEqual, true)
		So(len(i[1]), ShouldEqual, 1)
	})
}
