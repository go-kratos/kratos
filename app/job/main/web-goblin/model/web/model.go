package web

// ArcMsg archive .
type ArcMsg struct {
	Action string      `json:"action"`
	Table  string      `json:"table"`
	New    *ArchiveSub `json:"new"`
	Old    *ArchiveSub `json:"old"`
}

// ArchiveSub archive .
type ArchiveSub struct {
	Aid     int64  `json:"aid"`
	Mid     int64  `json:"mid"`
	PubTime string `json:"pubtime"`
	CTime   string `json:"ctime"`
	MTime   string `json:"mtime"`
	State   int    `json:"state"`
}
