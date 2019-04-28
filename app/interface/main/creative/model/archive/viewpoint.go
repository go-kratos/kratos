package archive

// ViewPointRow video highlight viewpoint
type ViewPointRow struct {
	ID     int64        `json:"id"`
	AID    int64        `json:"aid"`
	CID    int64        `json:"cid"`
	Points []*ViewPoint `json:"points"`
	State  int32        `json:"state"`
	CTime  string       `json:"ctime"`
	MTime  string       `json:"mtime"`
}

// ViewPoint viewpoint struct
type ViewPoint struct {
	Type    int8   `json:"type"`
	From    int    `json:"from"`
	To      int    `json:"to"`
	Content string `json:"content"`
	State   int8   `json:"state"`
}
