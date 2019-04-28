package model

import (
	"encoding/json"
)

// ES .
type ES struct {
	Addr string
}

// Page .
type Page struct {
	Pn    int   `json:"num"`
	Ps    int   `json:"size"`
	Total int64 `json:"total"`
}

// SearchResult search result (deprecated).
type SearchResult struct {
	Order  string            `json:"order"`
	Sort   string            `json:"sort"`
	Result []json.RawMessage `json:"result"`
	Debug  string            `json:"debug"`
	Page   *Page             `json:"page"`
}

// BasicSearchParams (deprecated).
type BasicSearchParams struct {
	AppID      string   `form:"appid" params:"appid"`
	Pattern    string   `form:"pattern" params:"pattern" default:"equal"` //关键字匹配模式，完成匹配：equal，模糊查询：like
	KW         string   `form:"kw" params:"kw"`
	KwFields   []string `form:"kw_fields,split" params:"kw_fields"`
	KWs        []string `form:"kws,split" params:"kws"` //关键词组，用于AND OR连接
	Order      []string `form:"order,split" params:"order"`
	Sort       []string `form:"sort,split" params:"sort" default:"desc"`
	Pn         int      `form:"pn" params:"pn;Range(1,5000)" default:"1"`
	Ps         int      `form:"ps" params:"ps;Range(1,1000)" default:"10"`
	Highlight  bool     `form:"highlight" params:"highlight" default:"false"`
	ScoreFirst bool     `form:"score_first" params:"score_first" default:"true"`
	Debug      bool     `form:"debug" params:"debug"`
	Source     []string
}

// BasicMNGSearchParams .
type BasicMNGSearchParams struct {
	Order string `form:"order" params:"order"`
	Sort  string `form:"sort" params:"sort" default:"desc"`
	Pn    int    `form:"pn" params:"pn;Range(1,5000)" default:"1"`
	Ps    int    `form:"ps" params:"ps;Range(1,1000)" default:"10"`
}

// BasicUpdateParams (deprecated).
type BasicUpdateParams struct {
	AppID string
}

// UpdateParams update params (deprecated).
type UpdateParams map[string]interface{}
