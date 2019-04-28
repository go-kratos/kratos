package search

import (
	"go-common/app/admin/main/workflow/model"
)

// ChallSearchCond is the condition model to send challenge search request
type ChallSearchCond struct {
	// Using int64 directly
	Cids      []int64
	Gids      []int64
	Mids      []int64
	Tids      []int64
	TagRounds []int64
	States    []int64

	Keyword   string
	CTimeFrom string
	CTimeTo   string

	PN    int64
	PS    int64
	Order string
	Sort  string
}

// FormatState .
func (cc *ChallSearchCond) FormatState() {
	for _, st := range cc.States {
		if st == model.QueueStateBefore {
			cc.States = append(cc.States, model.QueueState)
		}
	}
}

// ArcSearchResult is the model to parse search archive appeal result
type ArcSearchResult struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	TTL     int32  `json:"ttl"`

	Data struct {
		Page   *model.Page             `json:"page"`
		Result []GroupSearchCommonData `json:"result"`
	} `json:"data"`
}

// ChallSearchResult is the model to parse search challenge result
type ChallSearchResult struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	TTL     int32  `json:"ttl"`

	Data struct {
		Order string `json:"order"`
		Sort  string `json:"sort"`

		Page struct {
			Num   int64 `json:"num"`
			Size  int64 `json:"size"`
			Total int64 `json:"total"`
		} `json:"page"`

		Result []struct {
			ID    int64  `json:"id"`
			Gid   int64  `json:"gid"`
			Mid   int64  `json:"mid"`
			Tid   int64  `json:"tid"`
			CTime string `json:"ctime"`
		} `json:"result"`
	} `json:"data"`
}

// ChallListPage is the model for challenge list result
type ChallListPage struct {
	Items      []*model.Chall `json:"items"`
	TotalCount int32          `json:"total_count"`
	PN         int32          `json:"pn"`
	PS         int32          `json:"ps"`
}

// ChallListPageCommon model for challenge/list2 result
type ChallListPageCommon struct {
	Items []*model.Chall `json:"items"`
	Page  *model.Page    `json:"page"`
}

// ChallCount is the model for challenge count result
type ChallCount struct {
	TotalCount    int64          `json:"total_count"`
	BusinessCount map[int8]int64 `json:"business_count"`
}
