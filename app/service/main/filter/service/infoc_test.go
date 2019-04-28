package service

import (
	"testing"

	"go-common/app/service/main/filter/model/actriearea"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_HitReport(t *testing.T) {
	var mh []*actriearea.MatchHits
	mh = append(mh, &actriearea.MatchHits{Fid: 1})
	Convey("Test_CoverStart", t, func() {
		service.repostHitLog(ctx, "requestArea", "how are you", mh, "key")
	})
}
