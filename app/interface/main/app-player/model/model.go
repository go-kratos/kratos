package model

// Playurl is http://git.bilibili.co/video/playurl_doc/blob/master/PlayurlV2%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.md
type Playurl struct {
	From              string   `json:"from"`
	Result            string   `json:"result"`
	Quality           int64    `json:"quality"`
	Format            string   `json:"format"`
	Timelength        int64    `json:"timelength"`
	AcceptFormat      string   `json:"accept_format"`
	AcceptDescription []string `json:"accept_description,omitempty"`
	AcceptQuality     []int64  `json:"accept_quality"`
	VideoCodecid      int      `json:"video_codecid"`
	Fnver             int      `json:"fnver"`
	Fnval             int      `json:"fnval"`
	VideoProject      bool     `json:"video_project"`
	SeekParam         string   `json:"seek_param"`
	SeekType          string   `json:"seek_type"`
	Abtid             int      `json:"abtid,omitempty"`
	Durl              []*Durl  `json:"durl,omitempty"`
	Dash              *Dash    `json:"dash,omitempty"`
}

// FormatDash dash 驼峰 -> 下划线
func (p *Playurl) FormatDash() {
	if p.Dash != nil {
		var as, vs []*DashItem
		for _, a := range p.Dash.Audio {
			as = append(as, &DashItem{
				ID:        a.ID,
				BaseURL:   a.BaseURLRes,
				BackupURL: a.BackupURLRes,
				Bandwidth: a.Bandwidth,
				Codecid:   a.Codecid,
			})
		}
		for _, v := range p.Dash.Video {
			vs = append(vs, &DashItem{
				ID:        v.ID,
				BaseURL:   v.BaseURLRes,
				BackupURL: v.BackupURLRes,
				Bandwidth: v.Bandwidth,
				Codecid:   v.Codecid,
			})
		}
		p.Dash.Audio = as
		p.Dash.Video = vs
	}
}

// Durl is
type Durl struct {
	Order     int      `json:"order"`
	Length    int64    `json:"length"`
	Size      int64    `json:"size"`
	AHead     string   `json:"ahead,omitempty"`
	VHead     string   `json:"vhead,omitempty"`
	URL       string   `json:"url"`
	BackupURL []string `json:"backup_url,omitempty"`
}

//Param is
type Param struct {
	AID       int64  `form:"aid"`
	CID       int64  `form:"cid"`
	Qn        int64  `form:"qn"`
	Npcybs    int    `form:"npcybs"`
	Otype     string `form:"otype"`
	MobiApp   string `form:"mobi_app"`
	Fnver     int    `form:"fnver"`
	Fnval     int    `form:"fnval"`
	Session   string `form:"session"`
	Build     int    `form:"build"`
	Device    string `form:"device"`
	ForceHost int    `form:"force_host"`
}

// Dash is
type Dash struct {
	Video []*DashItem `json:"video"`
	Audio []*DashItem `json:"audio"`
}

// DashItem is
type DashItem struct {
	ID           int64    `json:"id"`
	BaseURL      string   `json:"base_url"`
	BackupURL    []string `json:"backup_url,omitempty"`
	BaseURLRes   string   `json:"baseUrl,omitempty"`
	BackupURLRes []string `json:"backupUrl,omitempty"`
	Bandwidth    int64    `json:"bandwidth"`
	Codecid      int64    `json:"codecid"`
}
