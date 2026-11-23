package log

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// SetDefault sets the default logger used by the package-level helpers and by
// [slog.Default].
func SetDefault(logger *slog.Logger) {
	slog.SetDefault(logger)
}

// Default returns the default logger.
func Default() *slog.Logger {
	return slog.Default()
}

// With returns a logger that includes the given attributes in each output
// operation. It mirrors [slog.Logger.With] on the default logger.
func With(args ...any) *slog.Logger {
	return slog.With(args...)
}

// WithGroup returns a logger that starts a group. It mirrors
// [slog.Logger.WithGroup] on the default logger.
func WithGroup(name string) *slog.Logger {
	return Default().WithGroup(name)
}

// Handler returns the default logger's handler. It mirrors
// [slog.Logger.Handler] on the default logger.
func Handler() slog.Handler {
	return Default().Handler()
}

// Enabled reports whether the default logger emits log records at the given
// context and level. It mirrors [slog.Logger.Enabled] on the default logger.
func Enabled(ctx context.Context, level Level) bool {
	return Default().Enabled(ctx, level)
}

// Debug logs at debug level. Signature mirrors [slog.Logger.Debug].
func Debug(msg string, args ...any) {
	log(context.Background(), LevelDebug, msg, args...)
}

// DebugContext logs at debug level with the provided context.
func DebugContext(ctx context.Context, msg string, args ...any) {
	log(ctx, LevelDebug, msg, args...)
}

// Info logs at info level. Signature mirrors [slog.Logger.Info].
func Info(msg string, args ...any) {
	log(context.Background(), LevelInfo, msg, args...)
}

// InfoContext logs at info level with the provided context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	log(ctx, LevelInfo, msg, args...)
}

// Warn logs at warn level. Signature mirrors [slog.Logger.Warn].
func Warn(msg string, args ...any) {
	log(context.Background(), LevelWarn, msg, args...)
}

// WarnContext logs at warn level with the provided context.
func WarnContext(ctx context.Context, msg string, args ...any) {
	log(ctx, LevelWarn, msg, args...)
}

// Error logs at error level. Signature mirrors [slog.Logger.Error].
func Error(msg string, args ...any) {
	log(context.Background(), LevelError, msg, args...)
}

// ErrorContext logs at error level with the provided context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	log(ctx, LevelError, msg, args...)
}

// Log emits a record at the given level. It mirrors [slog.Logger.Log] on the
// default logger.
func Log(ctx context.Context, level Level, msg string, args ...any) {
	log(ctx, level, msg, args...)
}

// LogAttrs emits a typed-attr record at the given level. It mirrors
// [slog.Logger.LogAttrs] on the default logger.
//
//nolint:revive // LogAttrs intentionally mirrors slog.Logger.LogAttrs.
func LogAttrs(ctx context.Context, level Level, msg string, attrs ...slog.Attr) {
	handler := slog.Default().Handler()
	if !handler.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// Skip [runtime.Callers, LogAttrs].
	runtime.Callers(2, pcs[:])
	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	record.AddAttrs(attrs...)
	_ = handler.Handle(ctx, record)
}

func log(ctx context.Context, level Level, msg string, args ...any) {
	handler := slog.Default().Handler()
	if !handler.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// Skip [runtime.Callers, log, exported helper].
	runtime.Callers(3, pcs[:])
	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	record.Add(args...)
	_ = handler.Handle(ctx, record)
}
