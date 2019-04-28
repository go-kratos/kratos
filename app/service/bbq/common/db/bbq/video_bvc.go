package bbq

//VideoBvc ..
type VideoBvc struct {
	SVID            int64  `json:"svid"`
	Path            string `json:"path"`
	ResolutionRetio string `json:"resolution_retio"`
	CodeRate        int16  `json:"code_rate"`
	VideoCode       string `json:"video_code"`
	FileSize        int64  `json:"file_size"`
	Duration        int64  `json:"duration"`
}
