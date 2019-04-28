package model

// ArgDMHistory .
type ArgDMHistory struct {
	Oid   int64
	Date  string
	Pn    int
	Ps    int
	Order string
	Sort  string
}

// DMHistory .
type DMHistory struct {
	ID int64 `json:"id"`
}
