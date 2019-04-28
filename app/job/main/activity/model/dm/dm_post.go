package dm

import "encoding/json"

const (
	// ActionPost dm action
	ActionPost = "post"
	// BroadcastCMD dm broadcast command
	BroadcastCMD = "DM"
	// BroadcastCMDACT BroadcastCMDACT
	BroadcastCMDACT = "ACT"
)

// XML dm xml info
type XML struct {
	PlayTime float32 `json:"playtime"`
	Mode     int32   `json:"mode"`
	FontSize int32   `json:"fontsize"`
	Color    int32   `json:"color"`
	Times    int64   `json:"times"`
	PoolID   int32   `json:"poolid"`
	Hash     string  `json:"hash"`
	DMID     int64   `json:"dmid"`
	Msg      string  `json:"msg"`
	Random   string  `json:"rnd"`
}

// Broadcast dm broadcast
type Broadcast struct {
	RoomID int64           `json:"roomid"`
	CMD    string          `json:"cmd"`
	Info   json.RawMessage `json:"info"`
}

// ActDM ActDM
type ActDM struct {
	Act    int64  `json:"act"`
	Aid    int64  `json:"aid"`
	Next   int64  `json:"next"`
	No     int64  `json:"no"`
	Yes    int64  `json:"yes"`
	Stage  int64  `json:"stage"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Tname  string `json:"tname"`
}
