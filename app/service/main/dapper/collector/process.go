package collector

import (
	"go-common/app/service/main/dapper/model"
)

// Processer span processer
type Processer interface {
	Process(span *model.Span) error
}

// ProcessFunc implement Processer
type ProcessFunc func(*model.Span) error

// Process implement Processer
func (p ProcessFunc) Process(span *model.Span) error { return p(span) }
