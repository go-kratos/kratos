package model

//Apply apply model
type Apply struct {
	Verify bool `json:"verify"`
	Forbid bool `json:"forbid"`
	Phone  int  `json:"phone"`
}

// Identify user verify info.
type Identify struct {
	Identify int `json:"identify"`
	Phone    int `json:"phone"`
}
