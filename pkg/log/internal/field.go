package core

// Field is for encoder
type Field interface {
	AddTo(enc ObjectEncoder)
}
