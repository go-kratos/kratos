package model

// Location .
type Location struct {
	ID    int32       `json:"id"`
	PID   int32       `json:"pid"`
	Name  string      `json:"name"`
	Child []*Location `json:"child,omitempty"`
}
