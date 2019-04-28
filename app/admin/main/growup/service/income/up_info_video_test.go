package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpInfoVideo(t *testing.T) {
	Convey("update up_info_video score", t, WithService(func(s *Service) {
		var (
			mid   = int64(1101)
			score = 100
		)
		tx, err := s.dao.BeginTran(context.Background())
		So(err, ShouldBeNil)
		err = s.TxUpdateUpInfoScore(tx, mid, score)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	}))
}
