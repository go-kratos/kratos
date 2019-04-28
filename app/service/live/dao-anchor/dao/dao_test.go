package dao

import (
	"context"
	"encoding/json"
	"flag"

	//"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/service/live/dao-anchor/conf"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	// TODO: other environments?
	flag.Set("conf", "../cmd/test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestCanConsume(t *testing.T) {
	flag.Set("conf", "../cmd/test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	Convey("", t, func(c C) {
		ctx := context.TODO()
		d := New(conf.Conf)

		msg := &databus.Message{
			Topic: "test-topic",
			Value: json.RawMessage(`{"msg_id":"test-msg-id", "other_key":"value"}`),
		}
		d.clearConsumed(ctx, msg)
		can := d.CanConsume(ctx, msg)
		So(can, ShouldBeTrue)

		can = d.CanConsume(ctx, msg)
		So(can, ShouldBeFalse)

		d.clearConsumed(ctx, msg)
		can = d.CanConsume(ctx, msg)
		So(can, ShouldBeTrue)
	})
}
