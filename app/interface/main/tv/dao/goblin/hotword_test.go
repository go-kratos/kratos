package goblin

import (
	"testing"

	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Hotword(t *testing.T) {
	Convey("Hotword Test", t, WithDao(func(d *Dao) {
		ctx := context.TODO()
		conn := d.mc.Get(ctx)
		s := []*model.Hotword{
			{
				Keyword: "Test1",
			},
			{
				Keyword: "Test2",
			},
		}
		defer conn.Close()
		err := conn.Set(&memcache.Item{Key: _hotwordKey, Object: s, Flags: memcache.FlagJSON, Expiration: 1200})
		So(err, ShouldBeNil)
		hotwordList, err := d.Hotword(ctx)
		So(err, ShouldBeNil)
		So(len(hotwordList), ShouldBeGreaterThan, 0)
	}))
}
