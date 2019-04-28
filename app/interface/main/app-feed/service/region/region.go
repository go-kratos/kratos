package region

import (
	"context"
	"time"

	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/tag"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emptyHotTags = []*tag.Hot{}
	_emptyTags    = []*tag.Tag{}
	_banRegion    = map[int16]struct{}{
		// bangumi
		33:  struct{}{},
		32:  struct{}{},
		153: struct{}{},
		51:  struct{}{},
		152: struct{}{},
		// music
		29:  struct{}{},
		54:  struct{}{},
		130: struct{}{},
		// tech
		37: struct{}{},
		96: struct{}{},
		// movie
		145: struct{}{},
		146: struct{}{},
		147: struct{}{},
		// TV series
		15: struct{}{},
		34: struct{}{},
		86: struct{}{},
		// entertainment
		71:  struct{}{},
		137: struct{}{},
		131: struct{}{},
	}
)

// HotTags get hot tags of region id.
func (s *Service) HotTags(c context.Context, mid int64, rid int16, ver string, plat int8, now time.Time) (hs []*tag.Hot, version string, err error) {
	if hs, err = s.tg.Hots(c, mid, rid, now); err != nil {
		log.Error("tg.HotTags(%d) error(%v)", rid, err)
		return
	}
	if model.IsOverseas(plat) {
		for _, hot := range hs {
			if _, ok := _banRegion[hot.Rid]; ok {
				hot.Tags = _emptyTags
			}
		}
	}
	if len(hs) == 0 {
		hs = _emptyHotTags
		return
	}
	version = s.md5(hs)
	if ver == version {
		err = ecode.NotModified
	}
	return
}

func (s *Service) SubTags(c context.Context, mid int64, pn, ps int) (t *tag.SubTag) {
	t = &tag.SubTag{}
	sub, err := s.tg.SubTags(c, mid, 0, pn, ps)
	if err != nil {
		log.Error("s.tg.SubTags(%d,%d,%d) error(%v)", mid, pn, ps, err)
		return
	}
	subTags := make([]*tag.Tag, 0, len(sub.Tags))
	for _, st := range sub.Tags {
		subTag := &tag.Tag{ID: st.ID, Name: st.Name, IsAtten: st.IsAtten}
		subTags = append(subTags, subTag)
	}
	if len(subTags) != 0 {
		t.SubTags = subTags
	} else {
		t.SubTags = _emptyTags
	}
	t.Count = sub.Total
	return
}

func (s *Service) AddTag(c context.Context, mid, tid int64, now time.Time) (err error) {
	if err = s.tg.Add(c, mid, tid, now); err != nil {
		log.Error("s.tg.Add(%d, %d) error(%v)", mid, tid, err)
	}
	return
}

func (s *Service) CancelTag(c context.Context, mid, tid int64, now time.Time) (err error) {
	if err = s.tg.Cancel(c, mid, tid, now); err != nil {
		log.Error("s.tg.Cancel(%d, %d) error(%v)", mid, tid, err)
	}
	return
}
