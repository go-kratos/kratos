package model

// SearchMemberResult is.
type SearchMemberResult struct {
	Order  string `json:"order"`
	Sort   string `json:"sort"`
	Result []struct {
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
	} `json:"result"`
	Page Page `json:"page"`
}

// Mids is.
func (r *SearchMemberResult) Mids() []int64 {
	mids := make([]int64, 0, len(r.Result))
	for _, r := range r.Result {
		mids = append(mids, r.Mid)
	}
	return mids
}

// Pagination is.
func (r *SearchMemberResult) Pagination() *CommonPagination {
	return &CommonPagination{
		Page: r.Page,
	}
}

// SearchUserPropertyReviewResult is.
type SearchUserPropertyReviewResult struct {
	Order  string `json:"order"`
	Sort   string `json:"sort"`
	Result []struct {
		ID int64 `json:"id"`
	} `json:"result"`
	Page Page `json:"page"`
}

// IDs is.
func (r *SearchUserPropertyReviewResult) IDs() []int64 {
	ids := make([]int64, 0, len(r.Result))
	for _, r := range r.Result {
		ids = append(ids, r.ID)
	}
	return ids
}

// Total is.
func (r *SearchUserPropertyReviewResult) Total() int {
	return r.Page.Total
}

// SearchLogResult is.
type SearchLogResult struct {
	Order  string     `json:"order"`
	Sort   string     `json:"sort"`
	Result []AuditLog `json:"result"`
	Page   Page       `json:"page"`
}

// AuditLog is.
type AuditLog struct {
	UID    int64  `json:"uid"`
	Uname  string `json:"uname"`
	OID    int64  `json:"oid"`
	Type   int8   `json:"type"`
	Action string `json:"action"`
	Str0   string `json:"str_0"`
	Str1   string `json:"str_1"`
	Str2   string `json:"str_2"`
	Int0   int    `json:"int_0"`
	Int1   int    `json:"int_1"`
	Int2   int    `json:"int_2"`
	Ctime  string `json:"ctime"`
	Extra  string `json:"extra_data"`
}
