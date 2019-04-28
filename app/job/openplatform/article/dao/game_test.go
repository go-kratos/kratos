package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GameSync(t *testing.T) {
	Convey("work", t, WithDao(func(d *Dao) {
		err := d.GameSync(context.Background(), "add", 1)
		So(err, ShouldBeNil)
	}))
}

func Test_NewGameCache(t *testing.T) {
	Convey("work", t, WithDao(func(d *Dao) {
		res, err := d.GameList(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
