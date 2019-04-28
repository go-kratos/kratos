package pgc

// PlayurlResp is the response struct from Playurl API
type PlayurlResp struct {
	Code          int     `json:"code"`
	Message       string  `json:"message"`
	From          string  `json:"from"`
	Result        string  `json:"result"`
	Quality       int     `json:"quality"`
	Format        string  `json:"format"`
	Timelength    int     `json:"timelength"`
	AcceptFormat  string  `json:"accept_format"`
	AcceptQuality []int   `json:"accept_quality"`
	SeekParam     string  `json:"seek_param"`
	SeekType      string  `json:"seek_type"`
	Durl          []*Durl `json:"durl"`
}

// Durl def.
type Durl struct {
	Order  int    `json:"order"`
	Length int    `json:"length"`
	Size   int    `json:"size"`
	URL    string `json:"url"`
}
