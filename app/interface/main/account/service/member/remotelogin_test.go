package member

import (
	"context"
	"testing"

	secmodel "go-common/app/service/main/secure/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Status(t *testing.T) {
	var (
		err error
		c   = context.TODO()
		res *secmodel.Msg
	)
	Convey("TestService_Status", func() {
		res, err = s.Status(c, 1, "xxxxx")
		if err != nil {
			t.Errorf("s.Status err(%v)", err)
		}
		t.Logf("s.Status result (%v)", res)
	})
}
