package health

type Status int

const (
	StatusDown Status = iota
	StatusUp
)

type CheckerStatus struct {
	Status
	Detail interface{}
	Err    error
}
