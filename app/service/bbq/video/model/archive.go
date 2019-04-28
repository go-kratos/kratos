package model

// ArchiveNotify .
type ArchiveNotify struct {
	Action string   `json:"action"`
	Table  string   `json:"table"`
	New    *Archive `json:"new"`
	Old    *Archive `json:"old"`
}

// Archive .
type Archive struct {
	ID          int    `json:"id"`
	AID         int64  `json:"aid"`
	CID         int64  `json:"cid"`
	MID         int64  `json:"mid"`
	TypeID      int32  `json:"typeid"`
	Videos      int    `json:"videos"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Duration    int    `json:"duration"`
	Attribute   int    `json:"attribute"`
	Copyright   int    `json:"copyright"`
	Access      int    `json:"access"`
	PubTime     string `json:"pubtime"`
	CTime       string `json:"ctime"`
	MTime       string `json:"mtime"`
	State       int    `json:"state"`
	MissionID   int    `json:"mission_id"`
	OrderID     int    `json:"order_id"`
	RedirectURL string `json:"redirect_url"`
	Forward     int    `json:"forward"`
	TID         int32  `json:"tid"`
	SubTID      int32  `json:"sub_tid"`
}

// ArchiveTypeResponse .
type ArchiveTypeResponse struct {
	Code    int                     `json:"code"`
	Data    map[string]*ArchiveType `json:"data"`
	Message string                  `json:"message"`
	TTL     int                     `json:"ttl"`
}

// ArchiveType .
type ArchiveType struct {
	ID   int    `json:"id"`
	PID  int    `json:"pid"`
	Name string `json:"name"`
}
