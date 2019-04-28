package report

// table name .
const (
	ArchiveClick   = "archive_click"
	ActiveDuration = "active_duration"
	PlayDuration   = "play_duration"
	VisitEvent     = "visit_event"
)

// DpCheckJobResult .
type DpCheckJobResult struct {
	Code      int      `json:"code"`
	Msg       string   `json:"msg"`
	StatusID  int      `json:"statusId"`
	StatusMsg string   `json:"statusMsg"`
	Files     []string `json:"hdfsPath"`
}

// ArcClickParam .
func ArcClickParam(b [][]byte) (res map[string]interface{}) {
	key := []string{"r_type", "avid", "cid", "part", "mid", "stime", "did", "ip", "ua", "buvid", "cookie_sid", "refer", "type", "sub_type", "sid", "epid", "platform", "device", "request_uri", "time_iso", "ip", "version", "buvid", "fts", "proid", "chid", "pid", "brand", "deviceid", "model", "osver", "ctime", "mid", "ver", "net", "oid", "eid", "start_time", "end_time", "duration", "openudid", "idfa", "mac", "is_coldstart", "session_id", "buvid_ext", "stime", "build", "buvid", "mobi_app", "platform", "session", "mid", "aid", "cid", "sid", "epid", "type", "sub_type", "quality", "total_time", "paused_time", "played_time", "video_duration", "play_type", "network_type", "last_play_progress_time", "max_play_progress_time", "device", "epid_status", "play_status", "user_status", "actual_played_time", "auto_play", "detail_play_time", "list_play_time", "request_uri", "time_iso", "ip", "version", "buvid", "fts", "proid", "chid", "pid", "brand", "deviceid", "model", "osver", "ctime", "mid", "ver", "net", "oid", "page_name", "page_arg", "ua", "h5_chid"}
	value := make([]string, len(b))
	res = make(map[string]interface{}, 4)
	mac := make(map[string]string, 18)
	mad := make(map[string]string, 28)
	mpd := make(map[string]string, 30)
	mve := make(map[string]string, 22)
	for k, v := range b {
		value[k] = string(v)
	}
	for k := range key {
		if k < len(value) {
			if k < 18 {
				mac[key[k]] = value[k]
			} else if k >= 18 && k < 46 {
				mad[key[k]] = value[k]
			} else if k >= 46 && k < 76 {
				mpd[key[k]] = value[k]
			} else if k >= 76 && k < 98 {
				mve[key[k]] = value[k]
			}
		}
	}
	res[ArchiveClick] = mac
	res[ActiveDuration] = mad
	res[PlayDuration] = mpd
	res[VisitEvent] = mve
	return
}
