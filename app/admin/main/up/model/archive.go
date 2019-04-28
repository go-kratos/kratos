package model

//Archive for db.
type Archive struct {
	ID    int64 `json:"id"`
	MID   int64 `json:"mid"`
	State int   `json:"state"`
}
