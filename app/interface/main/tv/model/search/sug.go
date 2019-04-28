package search

// SugResponse is the structure of search api response
type SugResponse struct {
	Code       int         `json:"code"`
	Cost       *Cost       `json:"cost"`
	Result     *Result     `json:"result"`
	PageCaches *PageCaches `json:"page caches"`
	Sengine    *Sengine    `json:"sengine"`
	Stoken     string      `json:"stoken"`
}

// Cost def.
type Cost struct {
	About *About `json:"about"`
}

// About def.
type About struct {
	ParamsCheck string `json:"params_check"`
	Total       string `json:"total"`
	MainHandler string `json:"main_handler"`
}

// Result def.
type Result struct {
	Tag []*STag `json:"tag"`
}

// STag def.
type STag struct {
	Value string `json:"value"`
	Ref   int    `json:"ref"`
	Name  string `json:"name"`
	Spid  int    `json:"spid"`
	Type  string `json:"type"`
}

// PageCaches def.
type PageCaches struct {
	SaveCache string `json:"save cache"`
}

// Sengine def.
type Sengine struct {
	Usage int `json:"usage"`
}

// ReqSug def.
type ReqSug struct {
	MobiApp  string `form:"mobi_app"`
	Build    string `form:"build"`
	Platform string `form:"platform"`
	Term     string `form:"term" validate:"required"`
}
