package mail

import (
	"flag"
	"fmt"
	"testing"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	var err error
	flag.Set("conf", "../../cmd/saga-admin-test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
}

// go test  -test.v -test.run TestMail
func TestMail(t *testing.T) {
	Convey("Test mail", t, func() {
		m := &model.Mail{
			ToAddress: []*model.MailAddress{{Name: "baihai", Address: "yubaihai@bilibili.com"},
				{Name: "muyan", Address: "muyang@bilibili.com"}},
			Subject: fmt.Sprintf("【Sage 提醒】%s项目发生Merge Request事件", "kk"),
		}
		mergeOut := " Merge made by the 'recursive' strategy.\n" +
			"tools/saga/CHANGELOG.md          |  4 ++++\n" +
			"business/interface/app-show/service/rank/rank.go  | 28 +++++++++++------------\n" +
			"business/interface/app-show/service/show/cache.go |  6 ++---\n" +
			"3 files changed, 21 insertions(+), 17 deletions(-)"
		err := SendMail(m, &model.MailData{
			UserName:     "baihai",
			SourceBranch: "featre_answer",
			TargetBranch: "master",
			Title:        "修改变量A",
			Description:  "内容就是",
			URL:          "http://www.baidu.com",
			Info:         mergeOut,
		})
		So(err, ShouldBeNil)
	})
}
