package service

import (
	"testing"

	"go-common/app/job/main/passport-encrypt/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_IncMigration(t *testing.T) {
	once.Do(startService)
	c := conf.Conf
	Convey("test full migration ", t, func() {
		s.asobinlogconsumeproc()
		for i := 0; i < c.Group.AsoBinLog.Num; i++ {
			ch := make(chan *message, c.Group.AsoBinLog.Chan)
			s.merges[i] = ch
			s.asobinlogmergeproc(ch)
		}
		s.asobinlogcommitproc()
	})
}
