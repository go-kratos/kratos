package search

//RespForClient def.
type RespForClient struct {
	*ResultResponse
	SearchType string          `json:"search_type"`
	PageInfo   *pageinfo       `json:"pageinfo"`
	ResultAll  *AllForClient   `json:"resultall"`
	PGC        []*CommonResult `json:"pgc"`
	UGC        []*CommonResult `json:"ugc"`
}

//AllForClient def.
type AllForClient struct {
	Pgc []*CommonResult `json:"tvpgc"`
	Ugc []*CommonResult `json:"tvugc"`
}
