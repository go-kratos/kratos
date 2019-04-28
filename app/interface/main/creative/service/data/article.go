package data

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/creative/model/data"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
)

// ThirtyDayArticle for 30 days.
func (s *Service) ThirtyDayArticle(c context.Context, mid int64, ip string) (res []*artmdl.ThirtyDayArticle, err error) {
	if res, err = s.data.ThirtyDayArticleCache(c, mid); err != nil && len(res) != 0 {
		return
	}
	if res, err = s.art.ThirtyDayArticle(c, mid, ip); err != nil && res != nil {
		s.data.AddCache(func() {
			s.data.AddThirtyDayArticleCache(context.Background(), mid, res)
		})
	}
	return
}

//ArtThirtyDay for article.
func (s *Service) ArtThirtyDay(c context.Context, mid int64, ty byte) (res []*data.ArtTrend, err error) {
	if res, err = s.data.ArtThirtyDay(c, mid, ty); err != nil {
		log.Error("s.data.ArtThirtyDay mid(%d)|err(%v)", mid, err)
	}
	return
}

//ArtRank for article rank.
func (s *Service) ArtRank(c context.Context, mid int64, ty byte) (res *data.ArtRankList, err error) {
	daytime := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour)
	dt := daytime.Format("20060102")

	var ar *data.ArtRankMap
	if ar, err = s.data.ArtRank(c, mid, ty, dt); err != nil {
		log.Error("s.data.ArtRank mid(%d)|err(%v)", mid, err)
		return
	}

	if ar == nil {
		log.Info("s.data.ArtRank mid(%d) article rank(%+v)", mid, ar)
		return
	}

	sortK := make([]int, 0, len(ar.AIDs))
	aids := make([]int64, 0, len(ar.AIDs))
	for k, v := range ar.AIDs {
		aids = append(aids, v)
		sortK = append(sortK, k)
	}
	artMap, _ := s.art.ArticleMetas(c, aids, "")
	if len(artMap) == 0 {
		return
	}
	arts := make([]*data.ArtMeta, 0, len(artMap))

	sort.Ints(sortK)
	for _, k := range sortK {
		if aid, ok := ar.AIDs[k]; ok {
			if a, ok := artMap[aid]; ok {
				art := &data.ArtMeta{}
				art.AID = a.ID
				art.Title = a.Title
				art.PTime = a.PublishTime
				if v, ok := ar.Incrs[k]; ok {
					art.Incr = v
				}
				arts = append(arts, art)
			}
		}
	}
	res = &data.ArtRankList{Arts: arts}
	return
}

//ArtReadAnalysis for article.
func (s *Service) ArtReadAnalysis(c context.Context, mid int64) (res *data.ArtRead, err error) {
	if res, err = s.data.ReadAnalysis(c, mid); err != nil {
		log.Error("s.data.ReadAnalysis mid(%d)|err(%v)", mid, err)
	}
	return
}
