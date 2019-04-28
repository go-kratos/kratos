package history

// HisParam fro history
type HisParam struct {
	MobiApp  string `form:"mobi_app"`
	Device   string `form:"device"`
	Build    int64  `form:"build"`
	Platform string `form:"platform"`
	Pn       int    `form:"pn"`
	Ps       int    `form:"ps"`
	Mid      int64  `form:"mid"`
	Max      int64  `form:"max"`
	MaxTP    int8   `form:"max_tp"`
	Business string `form:"business"`
}

// LiveParam statue param
type LiveParam struct {
	RoomIDs string `form:"room_ids"`
}

// DelParam del param
type DelParam struct {
	Mid   int64    `form:"mid"`
	Boids []string `form:"boids,split" validate:"min=1"`
}

// ClearParam clear param
type ClearParam struct {
	Mid      int64  `form:"mid"`
	Business string `form:"business"`
}

// ListRes for history
type ListRes struct {
	Title   string   `json:"title"`
	Covers  []string `json:"covers,omitempty"`
	Cover   string   `json:"cover,omitempty"`
	URI     string   `json:"uri"`
	History struct {
		Oid      int64  `json:"oid"`
		Tp       int8   `json:"tp"`
		Cid      int64  `json:"cid,omitempty"`
		Page     int32  `json:"page,omitempty"`
		Part     string `json:"part,omitempty"`
		Business string `json:"business"`
	} `json:"history"`
	Videos     int64  `json:"videos,omitempty"`
	Name       string `json:"name,omitempty"`
	Mid        int64  `json:"mid,omitempty"`
	Goto       string `json:"goto"`
	Badge      string `json:"badge,omitempty"`
	ViewAt     int64  `json:"view_at"`
	Progress   int64  `json:"progress,omitempty"`
	Duration   int64  `json:"duration,omitempty"`
	ShowTitle  string `json:"show_title,omitempty"`
	TagName    string `json:"tag_name,omitempty"`
	LiveStatus int    `json:"live_status,omitempty"`
	Current    string `json:"current,omitempty"`
	Total      string `json:"total,omitempty"`
	NewDesc    string `json:"new_desc,omitempty"`
	IsFinish   int8   `json:"is_finish,omitempty"`
}

// PGCRes for history
type PGCRes struct {
	EpID      int64  `json:"ep_id"`
	Cover     string `json:"cover"`
	URI       string `json:"uri"`
	Title     string `json:"title"`
	ShowTitle string `json:"show_title"`
	Season    struct {
		Title string `json:"title"`
	} `json:"season"`
}

// ListCursor for history
type ListCursor struct {
	Tab    []*BusTab  `json:"tab"`
	List   []*ListRes `json:"list"`
	Cursor *Cursor    `json:"cursor"`
}

// BusTab business tab
type BusTab struct {
	Business string `json:"business"`
	Name     string `json:"name"`
}

// Cursor for history
type Cursor struct {
	Max   int64 `json:"max"`
	MaxTP int8  `json:"max_tp"`
	Ps    int   `json:"ps"`
}
