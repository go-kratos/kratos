package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_checkUserPwd(t *testing.T) {
	Convey("decode", t, func() {
		once.Do(startService)
		var (
			originPwd = "b87dab1b55edfbe7d8191e011ff81067"
			salt      = "8ZMdzvzF"
			rsaPwd    = "VwXtLeolUuYiNKhgjcwaWZcoJAqQyaIRuaOxBrntvqLy4of6oYvAG10EhfBZiITPhBZ7Wgsmubbl/vMRumDPzZ6WnG8KOLoqqYbre4G9AByeGg6LOxL7r+nW78ZWB/xQfe8u8Z7uGs3FQL1og09qLmbyqevvoP/Q4ipQPOvLDSE="
		)
		err := s.checkUserPwd(rsaPwd, originPwd, salt)
		ShouldBeNil(err)
	})
}
