package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_FullMigration(t *testing.T) {
	once.Do(startService)
	Convey("test full migration ", t, func() {
		s.fullMigration(0, 10, 5, 50, "")
	})
}
