package dao

import (
	"encoding/hex"
	"flag"
	"os"
	"strconv"
	"strings"
	"testing"

	"go-common/app/interface/main/push-archive/conf"

	"github.com/smartystreets/goconvey/convey"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.push-archive")
		flag.Set("conf_token", "61c0d7d8527e8a4aad5b49826869e23c")
		flag.Set("tree_id", "7615")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/push-archive-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func Test_msgTemplateEncode(t *testing.T) {
	convey.Convey("msgTemplateDesc编码", t, func() {
		for _, g := range d.FanGroups {
			ascii := strconv.QuoteToASCII(g.MsgTemplateDesc)
			msgtemp := hex.EncodeToString([]byte(ascii))
			t.Logf("the group(%s) msgtemplate encoded(%v)\n", g.Name, msgtemp)

			ascii = strconv.QuoteToASCII(g.MsgTemplate)
			msgtemp2 := hex.EncodeToString([]byte(ascii))
			convey.So(msgtemp, convey.ShouldEqual, msgtemp2)
		}
	})
}

func Test_msgTemplateDecode(t *testing.T) {
	convey.Convey("msgTemplateDesc解码", t, func() {
		for _, g := range d.FanGroups {
			convey.So(g.MsgTemplate, convey.ShouldEqual, g.MsgTemplateDesc)
		}
	})
}

func Test_keyname(t *testing.T) {
	convey.Convey("fangroup keyname", t, func() {
		for gkey, g := range d.FanGroups {
			convey.So(gkey, convey.ShouldEqual, fanGroupKey(g.RelationType, g.Name))
		}
	})
}

func Test_conf(t *testing.T) {
	convey.Convey("配置结果", t, func() {
		for gkey, g := range d.FanGroups {
			convey.So(gkey, convey.ShouldEqual, fanGroupKey(g.RelationType, g.Name))
			convey.So(len(strings.Split(g.MsgTemplateDesc, "\r\n")), convey.ShouldEqual, 2)
			convey.So(g.MsgTemplate, convey.ShouldEqual, g.MsgTemplateDesc)
		}
		for i, g := range d.Proportions {
			proportion, _ := strconv.ParseFloat(d.c.ArcPush.Proportions[i].Proportion, 64)
			convey.So(g.MaxValue-g.MinValue+1, convey.ShouldEqual, proportion*100)
		}
		convey.So(len(d.GroupOrder), convey.ShouldEqual, len(d.c.ArcPush.Order))
	})
}
