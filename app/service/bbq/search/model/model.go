package model

//Query .
type Query struct {
	Calc   *Calc                  `json:"calc"`
	Where  *Where                 `json:"where"`
	Filter map[string]interface{} `json:"filter"`
	From   int                    `json:"from"`
	Size   int                    `json:"size"`
}

//Calc .
type Calc struct {
	Open       int64   `json:"open"`
	PlayRatio  float64 `json:"play_ratio"`
	FavRatio   float64 `json:"fav_ratio"`
	LikeRatio  float64 `json:"like_ratio"`
	ShareRatio float64 `json:"share_ratio"`
	CoinRatio  float64 `json:"coin_ratio"`
	ReplyRatio float64 `json:"reply_ratio"`
}

//Where .
type Where struct {
	In    map[string][]interface{} `json:"in"`
	NotIn map[string][]interface{} `json:"not_in"`
	Lte   map[string]int64         `json:"lte"`
	Gte   map[string]int64         `json:"gte"`
}

// EsParam es请求参数
type EsParam struct {
	From  int                               `json:"from"`
	Size  int                               `json:"size"`
	Query map[string]map[string]interface{} `json:"query"`
	Sort  []map[string]*Script              `json:"sort"`
}

// Script .
type Script struct {
	Order  string                 `json:"order"`
	Script map[string]interface{} `json:"script"`
	Type   string                 `json:"type"`
}
