package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Tag(t *testing.T) {
	Convey("Test_Tag info  service ", t, WithService(func(s *Service) {
		var (
			c           = context.Background()
			tid   int64 = 2090
			name        = "美丽"
			names       = []string{"美丽", "超"}
			tids        = []int64{2090, 2181}
			mid   int64 = 14771787
		)
		s.Info(c, mid, tid)
		s.Infos(c, mid, tids)
		s.InfoByName(c, mid, name)
		s.InfosByNames(c, mid, names)
		s.RecommandTag(c)
		s.HotMap(c)
		s.Prids(c)
		s.TagGroup(c)
		s.Count(c, tid)
		s.Counts(c, tids)
	}))
}
