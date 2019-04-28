package archive

import "go-common/app/admin/main/videoup/model/utils"

//Watermark 水印
type Watermark struct {
	ID       int64            `json:"id"`
	Info     string           `json:"info"`
	MD5      string           `json:"m5"`
	MID      string           `json:"mid"`
	Position string           `json:"position"`
	Type     string           `json:"type"`
	Uname    string           `json:"uname"`
	URL      string           `json:"url"`
	State    string           `json:"state"`
	MTime    utils.FormatTime `json:"mtime"`
}
