package model

// VerUpdate Params
type VerUpdate struct {
	MobiApp string
	Build   int
	Channel string
	Seed    int
	Sdkint  int
	Model   string
	OldID   string
}

// HTTPData response
type HTTPData struct {
	Ver     string `json:"ver"`
	Build   int    `json:"build"`
	Info    string `json:"info"`
	Size    string `json:"size"`
	URL     string `json:"url"`
	Hash    string `json:"hash"`
	Policy  int    `json:"policy"`
	IsForce int    `json:"is_force"`
	IsPush  int    `json:"is_push"`
	IsGray  int    `json:"is_gray"`
	Mtime   int    `json:"mtime"`
	Patch   *Patch `json:"patch"`
}

// Patch fix
type Patch struct {
	NewID string `json:"new_id"`
	OldID string `json:"old_id"`
	URL   string `json:"url"`
	Md5   string `json:"md5"`
	Size  int    `json:"size"`
}
