package email

import (
	"go-common/app/job/main/archive/model/result"
	"go-common/app/service/main/archive/api"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PGCNotifyMail(t *testing.T) {
	Convey("PGCNotifyMail", t, func() {
		d.PGCNotifyMail(&api.Arc{}, &result.Archive{}, &result.Archive{})
	})
}
