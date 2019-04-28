package email

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestEmailNotifyEmailTemplate(t *testing.T) {
	var (
		params = map[string]string{
			"aid":        "444",
			"title":      "随意一个标题",
			"upName":     "随意一个up主昵称",
			"condition":  "优质up主",
			"change":     "",
			"username":   "cxxx",
			"department": "r&d",
			"typeId":     "4",
			"fromVideo":  "false",
			"uid":        "441",
		}
	)
	convey.Convey("NotifyEmailTemplate", t, func(ctx convey.C) {
		tpl := d.NotifyEmailTemplate(params)
		ctx.Convey("Then tpl should not be nil.", func(ctx convey.C) {
			ctx.So(tpl, convey.ShouldNotBeNil)
		})
	})
}

func TestEmailPrivateEmailTemplate(t *testing.T) {
	var (
		params = map[string]string{
			"aid":        "444",
			"title":      "随意一个标题",
			"upName":     "随意一个up主昵称",
			"condition":  "优质up主",
			"change":     "",
			"username":   "cxxx",
			"department": "r&d",
			"typeId":     "4",
			"fromVideo":  "false",
			"uid":        "441",
		}
	)
	convey.Convey("PrivateEmailTemplate", t, func(ctx convey.C) {
		tpl := d.PrivateEmailTemplate(params)
		ctx.Convey("Then tpl should not be nil.", func(ctx convey.C) {
			ctx.So(tpl, convey.ShouldNotBeNil)
		})
	})
}
