package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TxUpAccountBreach(t *testing.T) {
	Convey("TxUpAccountBreach", t, WithService(func(s *Service) {
		tx, _ := s.dao.BeginTran(context.Background())
		err := s.TxUpAccountBreach(context.Background(), tx, 11, -1, 1)
		So(err, ShouldNotBeNil)
	}))

	Convey("TxUpAccountBreach", t, WithService(func(s *Service) {
		tx, _ := s.dao.BeginTran(context.Background())
		err := s.TxUpAccountBreach(context.Background(), tx, 11, 1, -1)
		So(err, ShouldNotBeNil)
	}))
}
