package model

// Otype playurl data type.
const (
	OtypeJSON = "json"
	OtypeXML  = "xml"
)

// PlayurlArg playurl arg.
type PlayurlArg struct {
	Cid           int64  `form:"cid" validate:"min=1"`
	Aid           int64  `form:"avid" validate:"min=1"`
	Qn            int    `form:"qn"`
	Type          string `form:"type"`
	MaxBackup     int    `form:"max_backup"`
	Npcybs        int    `form:"npcybs"`
	Platform      string `form:"platform"`
	Player        int    `form:"player"`
	Buvid         string `form:"buvid"`
	Resolution    string `form:"resolution"`
	Model         string `form:"model"`
	Build         int    `form:"build"`
	OType         string `form:"otype"`
	Fnver         int    `form:"fnver"`
	Fnval         int    `form:"fnval"`
	Session       string `form:"session"`
	HTML5         int    `form:"html5"`
	H5GoodQuality int    `form:"h5_good_quality"`
	HighQuality   int    `form:"high_quality"`
}

// PlayurlRes playurl res.
type PlayurlRes struct {
	From              string   `json:"from" xml:"from"`
	Result            string   `json:"result" xml:"result"`
	Message           string   `json:"message" xml:"message"`
	Quality           int      `json:"quality" xml:"quality"`
	Format            string   `json:"format" xml:"format"`
	Timelength        int64    `json:"timelength" xml:"timelength"`
	AcceptFormat      string   `json:"accept_format" xml:"accept_format"`
	AcceptDescription []string `json:"accept_description" xml:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality" xml:"accept_quality"`
	VideoCodeCid      int64    `json:"video_codecid" xml:"video_codecid"`
	SeekParam         string   `json:"seek_param" xml:"seek_param"`
	SeekType          string   `json:"seek_type" xml:"seek_type"`
	Abtid             int64    `json:"abtid,omitempty" xml:"abtid,omitempty"`
	Durl              []*struct {
		Order     int      `json:"order" xml:"order"`
		Length    int64    `json:"length" xml:"length"`
		Size      int64    `json:"size" xml:"size"`
		Ahead     string   `json:"ahead" xml:"ahead"`
		Vhead     string   `json:"vhead" xml:"vhead"`
		URL       string   `json:"url" xml:"url"`
		BackupURL []string `json:"backup_url" xml:"backup_url"`
	} `json:"durl,omitempty" xml:"durl,omitempty"`
	Dash *struct {
		Duration      int64       `json:"duration"`
		MinBufferTime float64     `json:"minBufferTime"`
		Video         []*DashItem `json:"video"`
		Audio         []*DashItem `json:"audio"`
	} `json:"dash,omitempty" xml:"dash,omitempty"`
}

// DashItem .
type DashItem struct {
	ID           int64    `json:"id" xml:"id"`
	BaseURL      string   `json:"baseUrl" xml:"baseUrl"`
	BackupURL    []string `json:"backupUrl" xml:"backupUrl"`
	Bandwidth    int64    `json:"bandwidth" xml:"bandwidth"`
	MimeType     string   `json:"mimeType" xml:"mimeType"`
	Codecs       string   `json:"codecs" xml:"codecs"`
	Width        int64    `json:"width" xml:"width"`
	Height       int64    `json:"height" xml:"height"`
	FrameRate    string   `json:"frameRate" xml:"frameRate"`
	Sar          string   `json:"sar" xml:"sar"`
	StartWithSAP int64    `json:"startWithSap" xml:"startWithSap"`
	SegmentBase  *struct {
		Initialization string `json:"Initialization" xml:"Initialization"`
		IndexRange     string `json:"indexRange" xml:"indexRange"`
	} `json:"SegmentBase" xml:"SegmentBase"`
	Codecid int64 `json:"codecid" xml:"codecid"`
}
