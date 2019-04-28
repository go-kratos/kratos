package discovery

// Status node
type Status struct {
	Status int         `json:"status"`
	Others interface{} `json:"others"`
}

// Addr addr
type Addr struct {
	Addr   string `json:"addr"`
	Status int    `json:"status"`
}
