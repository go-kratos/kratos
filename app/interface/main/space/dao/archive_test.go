package dao

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_ArcSearchList(t *testing.T) {
	convey.Convey("test search arc list", t, func(ctx convey.C) {
		arg := &model.SearchArg{
			Mid:     2,
			Tid:     0,
			Order:   "",
			Keyword: "",
			Pn:      1,
			Ps:      20,
		}
		data, count, err := d.ArcSearchList(context.Background(), arg)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%d", count)
		str, _ := json.Marshal(data)
		convey.Printf("%s", string(str))
	})
}
