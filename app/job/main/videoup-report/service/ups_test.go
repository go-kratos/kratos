package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/manager"
	"testing"
)

func TestHdlManagerUpsBinlog(t *testing.T) {
	Convey("hdlManagerUpsBinlog", t, func() {
		str := "{\"ctime\":\"2018-10-31 15:29:18\",\"id\":502,\"mid\":27515256,\"mtime\":\"2018-10-31 15:29:18\",\"note\":\"1111\",\"type\":18,\"uid\":277}"
		bs := []byte(str)
		m := &manager.BinMsg{
			Action: "insert",
			Table:  "ups",
			New:    bs,
		}
		s.hdlManagerUpsBinlog(m)
	})
}
