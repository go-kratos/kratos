package dao

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_ModPage(t *testing.T) {
	Convey("TestDao_ModPage", t, WithDao(func(d *Dao) {
		// home page
		homeMods, err := d.ModPage(ctx, 0)
		So(err, ShouldBeNil)
		So(len(homeMods), ShouldBeGreaterThan, 0)
		// jp page
		jpMods, err2 := d.ModPage(ctx, 1)
		So(err2, ShouldBeNil)
		So(len(jpMods), ShouldBeGreaterThan, 0)
	}))
}

func TestDao_PassedSns(t *testing.T) {
	Convey("TestDao_PassedSns", t, WithDao(func(d *Dao) {
		ids, err := d.PassedSns(ctx)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThan, 0)
		fmt.Println(ids)
	}))
}
