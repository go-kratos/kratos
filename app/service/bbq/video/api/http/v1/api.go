package v1

// ViewsAddRequest .
type ViewsAddRequest struct {
	Svid  int64 `json:"svid" form:"svid" validate:"required"`
	Views int   `json:"views" form:"views"`
}

// ViewsAddResponse .
type ViewsAddResponse struct {
	Affected int64 `json:"affected"`
}
