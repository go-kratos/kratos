package monitor

import (
	"context"
	"testing"

	"go-common/app/job/main/archive/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Send(t *testing.T) {
	Convey("Send", t, func() {
		err := d.Send(context.TODO(), conf.Conf.WeChantUsers, "报警短信test", conf.Conf.WeChatToken, conf.Conf.WeChatSecret)
		So(err, ShouldBeNil)
	})
}
