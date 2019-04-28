package model

// is
const (
	TypeForUpdateVideo   = int(0)
	TypeForDelVideo      = int(1)
	TypeForUpdateArchive = int(2)
)

// RetryItem struct
type RetryItem struct {
	Tp     int      `json:"type"`
	AID    int64    `json:"aid"`
	CID    int64    `json:"cid"`
	Old    *Archive `json:"new_archive"`
	Nw     *Archive `json:"old_archive"`
	Action string   `json:"action"`
}
