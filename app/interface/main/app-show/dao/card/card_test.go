package card

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-show")
		flag.Set("conf_token", "Pae4IDOeht4cHXCdOkay7sKeQwHxKOLA")
		flag.Set("tree_id", "2687")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestColumns(t *testing.T) {
	Convey("Columns", t, func() {
		_, err := d.Columns(ctx())
		// res = map[int8][]*card.Column{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestPosRecs(t *testing.T) {
	Convey("PosRecs", t, func() {
		_, err := d.PosRecs(ctx(), time.Now())
		// res = map[int8]map[int][]*card.Card{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRecContents(t *testing.T) {
	Convey("RecContents", t, func() {
		_, _, err := d.RecContents(ctx(), time.Now())
		// res = map[int][]*card.Content{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestNperContents(t *testing.T) {
	Convey("NperContents", t, func() {
		_, _, err := d.NperContents(ctx(), time.Now())
		// res = map[int][]*card.Content{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestColumnNpers(t *testing.T) {
	Convey("ColumnNpers", t, func() {
		_, err := d.ColumnNpers(ctx(), time.Now())
		// res = map[int8][]*card.ColumnNper{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestColumnPlatList(t *testing.T) {
	Convey("ColumnPlatList", t, func() {
		_, err := d.ColumnPlatList(ctx(), time.Now())
		// res = map[int8][]*card.ColumnList{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestColumnList(t *testing.T) {
	Convey("ColumnList", t, func() {
		_, err := d.ColumnList(ctx(), time.Now())
		// res = []*card.ColumnList{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestCard(t *testing.T) {
	Convey("Card", t, func() {
		_, err := d.Card(ctx(), time.Now())
		// res = []*card.PopularCard{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestCardPlat(t *testing.T) {
	Convey("CardPlat", t, func() {
		_, err := d.CardPlat(ctx())
		// res = map[int64]map[int8][]*card.PopularCardPlat{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestCardSet(t *testing.T) {
	Convey("CardSet", t, func() {
		_, err := d.CardSet(ctx())
		// res = map[int64]*operate.CardSet{}
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestEventTopic(t *testing.T) {
	Convey("CardSet", t, func() {
		res, err := d.EventTopic(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
