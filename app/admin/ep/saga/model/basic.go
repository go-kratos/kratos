package model

// QueryTypeItem ...
type QueryTypeItem struct {
	Name  interface{} `json:"name" validate:"required"`
	Value interface{} `json:"value"`
}
