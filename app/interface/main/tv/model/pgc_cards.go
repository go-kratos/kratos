package model

// PgcResponse is the result structure from PGC API
type PgcResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  map[string]*SeasonCard `json:"result"`
}

// SeasonCard is the result structure from PGC API
type SeasonCard struct {
	NewEP *PgcNewEP `json:"new_ep"`
}

// PgcNewEP is the result from pgc return pgc new ep value
type PgcNewEP struct {
	ID        int    `json:"id"`
	IndexShow string `json:"index_show"`
	Cover     string `json:"cover"`
}
