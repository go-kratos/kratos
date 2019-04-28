package model

// PairKey def
type PairKey struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// TeamInfoResp def
type TeamInfoResp struct {
	Department []*PairKey `json:"department"`
	Business   []*PairKey `json:"business"`
}

// Developer def
type Developer struct {
	Department string `json:"department"`
	Total      int    `json:"total"`
	Android    int    `json:"android"`
	Ios        int    `json:"ios"`
	Web        int    `json:"web"`
	Service    int    `json:"service"`
	Other      int    `json:"other"`
}
