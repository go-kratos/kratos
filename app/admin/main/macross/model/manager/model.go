package manager

import "go-common/library/time"

// Model for manager.
type Model struct {
	ModelID       int64     `json:"model_id"`
	System        string    `json:"-"`
	ModelName     string    `json:"model_name"`
	ModelFlag     string    `json:"model_flag"`
	HasDependence bool      `json:"has_dependence"`
	GitURL        string    `json:"git_url"`
	Count         int64     `json:"count"`
	CTime         time.Time `json:"-"`
	MTime         time.Time `json:"-"`
}

// Models model sorted.
type Models []*Model

func (a Models) Len() int           { return len(a) }
func (a Models) Less(i, j int) bool { return int64(a[i].ModelID) < int64(a[j].ModelID) }
func (a Models) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
