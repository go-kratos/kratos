package model

//Archive for db.
type Archive struct {
	AID    int64 `json:"aid"`
	MID    int64 `json:"mid"`
	State  int   `json:"state"`
	UpFrom int8  `json:"up_from"`
}

// ArcVideo str
type ArcVideo struct {
	Archive *Archive
}
