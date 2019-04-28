package service

import (
	"bytes"
	"container/heap"
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/up-rating/model"
)

// ratingTop get top ups
func (s *Service) ratingTop(c context.Context, date time.Time, source chan []*model.Rating) (topRating map[int]map[int64]*RatingHeap, err error) {
	topRating = make(map[int]map[int64]*RatingHeap) // map[ctype][tagID]
	topRating[CreativeType] = make(map[int64]*RatingHeap)
	topRating[InfluenceType] = make(map[int64]*RatingHeap)
	for rating := range source {
		for _, r := range rating {
			if _, ok := topRating[CreativeType][r.TagID]; !ok {
				topRating[CreativeType][r.TagID] = &RatingHeap{}
			}
			pushTopRating(topRating[CreativeType][r.TagID], CreativeType, r)
			if _, ok := topRating[InfluenceType][r.TagID]; !ok {
				topRating[InfluenceType][r.TagID] = &RatingHeap{}
			}
			pushTopRating(topRating[InfluenceType][r.TagID], InfluenceType, r)
		}
	}
	return
}

func pushTopRating(h *RatingHeap, ctype int, r *model.Rating) {
	tr := &model.TopRating{
		MID:   r.MID,
		CType: ctype,
		TagID: r.TagID,
	}
	switch ctype {
	case CreativeType:
		tr.Score = r.CreativityScore
	case InfluenceType:
		tr.Score = r.InfluenceScore
	}
	heap.Push(h, tr)
	if h.Len() > 10 {
		heap.Pop(h)
	}
}

func (s *Service) insertTopRating(c context.Context, date time.Time, topRating map[int]map[int64]*RatingHeap, baseInfo map[int64]*model.BaseInfo) (rows int64, err error) {
	return s.dao.InsertTopRating(c, assemberTopRating(date, topRating, baseInfo))
}

func assemberTopRating(date time.Time, topRating map[int]map[int64]*RatingHeap, baseInfo map[int64]*model.BaseInfo) (vals string) {
	var buf bytes.Buffer
	for _, tagTop := range topRating {
		for _, h := range tagTop {
			for h.Len() > 0 {
				tr := heap.Pop(h).(*model.TopRating)
				info := baseInfo[tr.MID]
				if info == nil {
					info = &model.BaseInfo{}
				}
				buf.WriteString("(")
				buf.WriteString(strconv.FormatInt(tr.MID, 10))
				buf.WriteByte(',')
				buf.WriteString(strconv.Itoa(tr.CType))
				buf.WriteByte(',')
				buf.WriteString(strconv.FormatInt(tr.TagID, 10))
				buf.WriteByte(',')
				buf.WriteString(strconv.FormatInt(tr.Score, 10))
				buf.WriteByte(',')
				buf.WriteString(strconv.FormatInt(info.TotalFans, 10))
				buf.WriteByte(',')
				buf.WriteString(strconv.FormatInt(info.TotalPlay, 10))
				buf.WriteByte(',')
				buf.WriteString("'" + date.Format(_layout) + "'")
				buf.WriteString(")")
				buf.WriteByte(',')
			}
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}

// RatingHeap rating heap for topK
type RatingHeap []*model.TopRating

// Len len
func (r RatingHeap) Len() int { return len(r) }

// Less less
func (r RatingHeap) Less(i, j int) bool { return r[i].Score < r[j].Score }

// Swap swap
func (r RatingHeap) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Push push to heap
func (r *RatingHeap) Push(x interface{}) {
	*r = append(*r, x.(*model.TopRating))
}

// Pop pop from heap
func (r *RatingHeap) Pop() interface{} {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}
