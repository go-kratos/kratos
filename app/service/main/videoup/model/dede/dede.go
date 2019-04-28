package dede

type PadInfo struct {
	Aid      int64            `json:"aid"`
	Mid      int64            `json:"mid"`
	Fnm      map[string]int64 `json:"fnm"`
	Paded    bool             `json:"paded"`
	OK       chan bool        `json:"-"`
	IsUpload bool             `json:"is_upload"`
	CodeMode bool             `json:"code_mode"`
	IsUGC    bool             `json:"is_ugc"`
}
