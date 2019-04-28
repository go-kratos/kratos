package model

import "go-common/library/time"

// const const value.
const (
	HotRegionTag  = int32(0)
	HotArchiveTag = int32(1)
)

// RankCount RankCount.
type RankCount struct {
	ID          int64     `json:"id"`
	Prid        int64     `json:"prid"`
	Rid         int64     `json:"rid"`
	Type        int32     `json:"type"`
	TopCount    int64     `json:"top"`
	ViewCount   int64     `json:"view"`
	FilterCount int64     `json:"filter"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
}

// RankFilter RankFilter.
type RankFilter struct {
	ID      int64     `json:"-"`
	Tid     int64     `json:"tid"`
	TName   string    `json:"tname"`
	TagType int64     `json:"tag_type"`
	Rank    int64     `json:"rank"`
	CTime   time.Time `json:"ctime"`
	MTime   time.Time `json:"mtime"`
}

// RankTop RankTop.
type RankTop struct {
	ID        int64     `json:"-"`
	Tid       int64     `json:"tid"`
	TName     string    `json:"tname"`
	TagType   int64     `json:"tag_type"`
	HighLight int64     `json:"highlight"`
	Rank      int64     `json:"rank"`
	Business  int32     `json:"business"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// RankResult RankResult.
type RankResult struct {
	ID        int64     `json:"-"`
	Tid       int64     `json:"tid"`
	TName     string    `json:"tname"`
	TagType   int64     `json:"tag_type"`
	Rank      int64     `json:"rank"`
	HighLight int64     `json:"highlight"`
	Business  int32     `json:"business"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// BasicTag BasicTag.
type BasicTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// HotRank HotRank.
type HotRank struct {
	Type   int32         `json:"type"`
	Prid   int64         `json:"prid"`
	Rid    int64         `json:"rid"`
	View   []*RankResult `json:"view"`
	Filter []*RankFilter `json:"filter"` //Table: rank_filter
	Top    []*RankTop    `json:"top"`    //Table: rank_top
}

// RankResultSort RankResultSort.
type RankResultSort []*RankResult

// Len Len.
func (t RankResultSort) Len() int {
	return len(t)
}

// Less Less.
func (t RankResultSort) Less(i, j int) bool {
	if t[i].Rank < t[j].Rank {
		return true
	} else if t[i].Rank == t[j].Rank {
		if t[i].Tid < t[j].Tid {
			return true
		}
	}
	return false
}

// Swap Swap.
func (t RankResultSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// RankTopSort RankTopSort.
type RankTopSort []*RankTop

// Len Len.
func (t RankTopSort) Len() int {
	return len(t)
}

// Less Less.
func (t RankTopSort) Less(i, j int) bool {
	if t[i].Rank < t[j].Rank {
		return true
	} else if t[i].Rank == t[j].Rank {
		if t[i].Tid < t[j].Tid {
			return true
		}
	}
	return false
}

// Swap Swap.
func (t RankTopSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// RankFilterSort RankFilterSort.
type RankFilterSort []*RankFilter

// Len Len.
func (t RankFilterSort) Len() int {
	return len(t)
}

// Less Less.
func (t RankFilterSort) Less(i, j int) bool {
	if t[i].Rank < t[j].Rank {
		return true
	} else if t[i].Rank == t[j].Rank {
		if t[i].Tid < t[j].Tid {
			return true
		}
	}
	return false
}

// Swap Swap.
func (t RankFilterSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
