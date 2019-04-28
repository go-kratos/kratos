package v1

//HotWordRequest .
type HotWordRequest struct {
}

//HotWordResponse .
type HotWordResponse struct {
	List []string `json:"list,omitempty"`
}

// VideoSearchList 搜索视频结构
type VideoSearchList struct {
	VideoResponse
	TitleHighlight string   `json:"title_highlight"`
	HitColumns     []string `json:"hit_columns,omitempty"`
	CursorValue    string   `json:"cursor_value"`
	Offset         int64    // 这里不返回给客户端，只是为了代码方便处理
}

// VideoSearchRes 搜索用视频结果
type VideoSearchRes struct {
	List    []*VideoSearchList `json:"list,omitempty"`
	NumPage int64              `json:"numPages"`
	Page    int64              `json:"page"`
	HasMore bool               `json:"has_more"`
}

// UserSearchRes 搜索用户结果
type UserSearchRes struct {
	List    []*UserSearchList `json:"list,omitempty"`
	NumPage int64             `json:"numPages"`
	Page    int64             `json:"page"`
	HasMore bool              `json:"has_more"`
}

// UserSearchList 搜索用户结构
type UserSearchList struct {
	UserInfo
	UserStatic     *UserStatic `json:"user_statistics"`
	UnameHighlight string      `json:"uname_highlight"`
	HitColumns     []string    `json:"hit_columns"`
	CursorValue    string      `json:"cursor_value"`
	Offset         int64       // 这里不返回给客户端，只是为了代码方便处理
}

// UserStatic 用户统计信息
type UserStatic struct {
	Fan         int64 `json:"fan"`
	Follow      int64 `json:"follow"`
	Like        int64 `json:"like"`
	Liked       int64 `json:"liked"`
	FollowState int8  `json:"follow_state"`
}

// BaseSearchReq 基础搜索请求
type BaseSearchReq struct {
	Key       string `form:"keyword" validate:"required"`
	Page      int64  `form:"page"`
	PageSize  int64  `form:"pagesize"`
	Highlight int8   `form:"highlight"`
	Qn        int64  `form:"qn"`

	// TODO:v2接口，当page=0时生效，由于不久会拆接口，因此这里就复用老接口
	CursorPrev string `form:"cursor_prev"`
	CursorNext string `form:"cursor_next"`
}

// SugTag sug tag结构
type SugTag struct {
	Value string `json:"value"`
	Name  string `json:"name" `
	Type  string `json:"type"`
	Ref   int64  `json:"ref"`
}

// SugReq sug请求
type SugReq struct {
	KeyWord   string `form:"keyword" validate:"required"`
	PageSize  int64  `form:"pagesize"`
	Highlight int8   `form:"highlight"`
}
