package model

const (
	// OwnerInstall is_activated=1.
	OwnerInstall = 1
	// OwnerUninstall is_activated=0.
	OwnerUninstall = 0
)

// UpInfo .
type UpInfo struct {
	Code int            `json:"code"`
	Data []*UpInfoGroup `json:"results"`
}

// UpInfoGroup .
type UpInfoGroup struct {
	ID   int64   `json:"id"`
	Desc string  `json:"desc"`
	Mids []int64 `json:"mids"`
}
