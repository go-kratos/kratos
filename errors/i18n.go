package errors

import (
	"context"
	"sync"
)

// I18nMessage An interface to internationalize error messages is defined
type I18nMessage interface {
	// Localize Localization of error causes based on context and data
	Localize(ctx context.Context, reason string, data any) string
}

// The global i18n manager
var globalI18n I18nMessage
var globalI18nOnce sync.Once

// RegisterI18nManager Register the global i18n manager
func RegisterI18nManager(i18n I18nMessage) {
	globalI18nOnce.Do(func() {
		globalI18n = i18n
	})
}
