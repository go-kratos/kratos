package model

// CompareRes the result of comparing aso account between local and cloud.
type CompareRes struct {
	Flags          uint8             `json:"flag"`
	FlagsDesc      string            `json:"flags_desc"`
	Seq            int64             `json:"seq"`
	Local          *OriginAsoAccount `json:"local"`
	LocalEncrypted *AsoAccount       `json:"local_encrypted"`
	Cloud          *AsoAccount       `json:"cloud"`
}

// DiffParseResp diff parse resp.
type DiffParseResp struct {
	Total            int                   `json:"total"`
	SeqAndPercents   []*SeqCountAndPercent `json:"seq_and_percents"`
	CountAndPercents []*CountAndPercent    `json:"count_and_percents"`
	CompareResList   []*CompareRes         `json:"compare_res_list"`
}

// CountAndPercent count and percent.
type CountAndPercent struct {
	DiffType string `json:"diff_type"`
	Count    int    `json:"count"`
	Percent  string `json:"percent"`
}

// SeqCountAndPercent process goroutine seq count and percent.
type SeqCountAndPercent struct {
	Seq     int64  `json:"seq"`
	Count   int    `json:"count"`
	Percent string `json:"percent"`
}
