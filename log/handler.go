package log

import (
	"context"
	"log/slog"
)

// discardHandler drops all records. It mirrors discardHandler{} from go1.24
// while keeping the module compatible with go1.22.
type discardHandler struct{}

func (discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (h discardHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h discardHandler) WithGroup(string) slog.Handler           { return h }
