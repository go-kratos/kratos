package model

import artmdl "go-common/app/interface/openplatform/article/model"

var (
	// ArticleSortType article list sort types.
	ArticleSortType = map[string]int{
		"publish_time": artmdl.FieldDefault,
		"view":         artmdl.FieldView,
		"fav":          artmdl.FieldFav,
	}
	// PrivacyFields privacy allowed field.
	PrivacyFields = []string{
		"bangumi",
		"tags",
		"fav_video",
		"coins_video",
		"groups",
		"played_game",
		"channel",
		"user_info",
		"likes_video",
	}
	//ArcCheckType search arc check type.
	ArcCheckType = map[string]int{
		"channel": 1,
	}
)

// Page page return data struct.
type Page struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

// SearchArg arc search param.
type SearchArg struct {
	Mid       int64  `form:"mid" validate:"gt=0"`
	Tid       int64  `form:"tid"`
	Order     string `form:"order"`
	Keyword   string `form:"keyword"`
	Pn        int    `form:"pn" validate:"gt=0"`
	Ps        int    `form:"ps" validate:"gt=0,lte=100"`
	CheckType string `form:"check_type"`
	CheckID   int64  `form:"check_id"`
}

// WebIndex .
type WebIndex struct {
	Account *AccInfo `json:"account"`
	Setting *Setting `json:"setting"`
	Archive *WebArc  `json:"archive"`
}

// WebArc .
type WebArc struct {
	Page     WebPage    `json:"page"`
	Archives []*ArcItem `json:"archives"`
}

// WebPage .
type WebPage struct {
	Pn    int32 `json:"pn"`
	Ps    int32 `json:"ps"`
	Count int64 `json:"count"`
}
