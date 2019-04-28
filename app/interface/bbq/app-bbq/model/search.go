package model

import "encoding/json"

// SearchBaseReq 搜索基本请求
type SearchBaseReq struct {
	KeyWord   string `json:"keyword"`
	Type      string `json:"search_type"`
	Page      int64  `json:"page"`
	PageSize  int64  `json:"pagesize"`
	Highlight int64  `json:"highlight"`
	Platform  string `json:"platform"`
	MobiApp   string `json:"mobi_app"`
	Build     string `json:"build"`
	Device    string `json:"device"`
}

// SearchBaseRet 搜索基本返回
type SearchBaseRet struct {
	Code     int64  `json:"code"`
	NumPages int64  `json:"numPages"`
	PageSize int64  `json:"pagesize"`
	Seid     string `json:"seid"`
	Msg      string `json:"msg"`
	Page     int64  `json:"page"`
}

// VideoSearchRet 视频搜索结果
type VideoSearchRet struct {
	SearchBaseRet
	Result []*VideoSearchResult `json:"result,omitempty"`
}

// VideoSearchResult 视频搜索result
type VideoSearchResult struct {
	ID         int32    `json:"id"`
	Title      string   `json:"title"`
	HitColumns []string `json:"hit_columns,omitempty"`
}

// UserSearchResult 用户搜索结果
type UserSearchResult struct {
	ID         int64    `json:"id"`
	Uname      string   `json:"uname"`
	HitColumns []string `json:"hit_columns"`
}

// RawSearchRes .
type RawSearchRes struct {
	Code    int             `json:"code"`
	SeID    string          `json:"seid"`
	Msg     string          `json:"msg"`
	Page    int64           `json:"page"`
	PageNum int64           `json:"NumPages"`
	Res     json.RawMessage `json:"Result"`
}

// SugBaseReq Sug基本请求
type SugBaseReq struct {
	Term        string `json:"term"`
	SuggestType string `json:"suggest_type"`
	MainVer     string `json:"main_ver"`
	SugNum      int64  `json:"sug_num"`
	Highlight   int64  `json:"highlight"`
	Platform    string `json:"platform"`
	MobiApp     string `json:"mobi_app"`
	Build       string `json:"build"`
	Device      string `json:"device"`
}

// RawSugTag SugTag结构
type RawSugTag struct {
	Value string `json:"value"`
	Ref   int64  `json:"ref"`
	Name  string `json:"name"`
	Spid  int64  `json:"spid"`
	Type  string `json:"type"`
}
