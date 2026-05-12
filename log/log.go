package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

// Log emits a record at the given level. It mirrors [slog.Log].
func Log(ctx context.Context, level Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}

// LogAttrs emits a typed-attr record at the given level. It mirrors
// [slog.LogAttrs].
//
//nolint:revive // LogAttrs intentionally mirrors slog.LogAttrs.
func LogAttrs(ctx context.Context, level Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, level, msg, attrs...)
}

// Debug logs at debug level. Signature mirrors [slog.Debug].
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// DebugContext logs at debug level with the provided context.
func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// Debugf logs at debug level using fmt-style formatting. Kratos extension;
// slog has no fmt variant.
func Debugf(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

// DebugfContext is the context-aware variant of [Debugf].
func DebugfContext(ctx context.Context, format string, args ...any) {
	slog.DebugContext(ctx, fmt.Sprintf(format, args...))
}

// Info logs at info level. Signature mirrors [slog.Info].
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// InfoContext logs at info level with the provided context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// Infof logs at info level using fmt-style formatting.
func Infof(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

// InfofContext is the context-aware variant of [Infof].
func InfofContext(ctx context.Context, format string, args ...any) {
	slog.InfoContext(ctx, fmt.Sprintf(format, args...))
}

// Warn logs at warn level. Signature mirrors [slog.Warn].
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// WarnContext logs at warn level with the provided context.
func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// Warnf logs at warn level using fmt-style formatting.
func Warnf(format string, args ...any) {
	slog.Warn(fmt.Sprintf(format, args...))
}

// WarnfContext is the context-aware variant of [Warnf].
func WarnfContext(ctx context.Context, format string, args ...any) {
	slog.WarnContext(ctx, fmt.Sprintf(format, args...))
}

// Error logs at error level. Signature mirrors [slog.Error].
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// ErrorContext logs at error level with the provided context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// Errorf logs at error level using fmt-style formatting.
func Errorf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
}

// ErrorfContext is the context-aware variant of [Errorf].
func ErrorfContext(ctx context.Context, format string, args ...any) {
	slog.ErrorContext(ctx, fmt.Sprintf(format, args...))
}

// Fatal logs at [LevelFatal] and then calls os.Exit(1). Kratos extension.
func Fatal(msg string, args ...any) {
	slog.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// FatalContext is the context-aware variant of [Fatal].
func FatalContext(ctx context.Context, msg string, args ...any) {
	slog.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatalf logs at [LevelFatal] with fmt-style formatting and then exits.
func Fatalf(format string, args ...any) {
	slog.Log(context.Background(), LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// FatalfContext is the context-aware variant of [Fatalf].
func FatalfContext(ctx context.Context, format string, args ...any) {
	slog.Log(ctx, LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}
