package service

import (
	"context"
	"go-common/app/service/main/search/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PgcMedia(t *testing.T) {

	var (
		err error
		c   = context.TODO()
		sp  *model.PgcMediaParams
	)

	Convey("PgcMedia", t, WithService(func(s *Service) {
		_, err = s.PgcMedia(c, sp)
		So(err, ShouldBeNil)
	}))
}
