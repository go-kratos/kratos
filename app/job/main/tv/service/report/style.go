package report

import (
	"context"
	"strings"
	"time"

	mdlpgc "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

func (s *Service) showStyle() {
	var (
		err      error
		res      []*mdlpgc.StyleRes
		styleStr []*mdlpgc.ParamStyle
		styleRes = make(map[int][]*mdlpgc.ParamStyle)
		ctx      = context.Background()
	)
	for {
		if res, err = s.dao.FindStyle(ctx); err != nil {
			log.Error("s.dao.FindStyle error(%v)", err)
			time.Sleep(time.Second * 5)
			continue
		}
		if len(res) != 0 {
			for _, v := range res {
				styleStr = make([]*mdlpgc.ParamStyle, 0)
				if m, ok := s.labelRes[v.Category]; ok {
					a := strings.Split(v.Style, ",")
					for _, v1 := range a {
						r := &mdlpgc.ParamStyle{}
						if m1, ok1 := m[v1]; ok1 {
							r.Name = v1
							r.StyleID = m1
							styleStr = append(styleStr, r)
						}
					}
					if len(styleStr) != 0 {
						styleRes[v.ID] = styleStr
					}
				}

			}
		}
		if len(styleRes) > 0 {
			s.cache.Do(ctx, func(ctx context.Context) {
				// set style data to mc
				s.dao.SetStyleCache(ctx, styleRes)
			})
		}
		time.Sleep(time.Duration(s.c.Style.StyleSpan))
	}
}

func (s *Service) showLabel() {
	var (
		err error
		res map[int]map[string]int
		ctx = context.Background()
	)
	for {
		if res, err = s.dao.FindLabelID(ctx); err != nil {
			log.Error("s.dao.FindLabelID error(%v)", err)
			time.Sleep(time.Second * 5)
			continue
		}
		if len(res) != 0 {
			s.labelRes = res
			s.cache.Do(ctx, func(ctx context.Context) {
				// set label data to mc
				s.dao.SetLabelCache(ctx, s.labelRes)
			})
		}
		time.Sleep(time.Duration(s.c.Style.LabelSpan))
	}
}

func (s *Service) readLabelCache() {
	var (
		err error
		m   map[int]map[string]int
	)
	if m, err = s.dao.GetLabelCache(context.Background()); err != nil {
		log.Error("s.dao.GetLabelCache error(%v)", err)
		panic(err)
	}
	s.labelRes = m
}
