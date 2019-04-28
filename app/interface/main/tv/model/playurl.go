package model

// PlayURLReq is used for getting ugc play url param from app
type PlayURLReq struct {
	Platform  string `form:"platform" validate:"required"`
	Device    string `form:"device"`
	Expire    string `form:"expire"`
	Cid       string `form:"cid" validate:"required"`
	Avid      int64  `form:"avid" validate:"required"`
	Build     string `form:"build"`
	Qn        string `form:"qn"`
	Mid       string `form:"mid"`
	Npcybs    string `form:"npcybs"`
	Buvid     string `form:"buvid"`
	TrackPath string `form:"track_path"`
	AccessKey string `form:"access_key"`
}

//PlayURLResp is used for return ugc play url result
type PlayURLResp struct {
	Code              int      `json:"code"`
	Result            string   `json:"result"`
	Message           string   `json:"message"`
	From              string   `json:"from"`
	Quality           int      `json:"quality"`
	Format            string   `json:"format"`
	Timelength        int      `json:"timelength"`
	AcceptFormat      string   `json:"accept_format"`
	AcceptDescription []string `json:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality"`
	AcceptWatermark   []bool   `json:"accept_watermark"`
	VideoCodecid      int      `json:"video_codecid"`
	VideoProject      bool     `json:"video_project"`
	SeekParam         string   `json:"seek_param"`
	SeekType          string   `json:"seek_type"`
	Durl              []struct {
		Order  int    `json:"order"`
		Length int    `json:"length"`
		Size   int    `json:"size"`
		Ahead  string `json:"ahead"`
		Vhead  string `json:"vhead"`
		URL    string `json:"url"`
	} `json:"durl"`
}
