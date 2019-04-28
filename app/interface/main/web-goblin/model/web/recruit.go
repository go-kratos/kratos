package web

// Params .
type Params struct {
	Route string `form:"route" validate:"required"`
	Mode  string `form:"mode" validate:"required"`
	JobID string `form:"job_id"`
}
