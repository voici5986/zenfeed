// Copyright (C) 2025 wangyusong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	slogdedup "github.com/veqryn/slog-dedup"

	"github.com/glidea/zenfeed/pkg/model"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func SetLevel(level Level) error {
	if level == "" {
		level = LevelInfo
	}

	var logLevel slog.Level
	switch level {
	case LevelDebug:
		logLevel = slog.LevelDebug
	case LevelInfo:
		logLevel = slog.LevelInfo
	case LevelWarn:
		logLevel = slog.LevelWarn
	case LevelError:
		logLevel = slog.LevelError
	default:
		return errors.Errorf("invalid log level, valid values are: %v", []Level{LevelDebug, LevelInfo, LevelWarn, LevelError})
	}

	newLogger := slog.New(slogdedup.NewOverwriteHandler(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
		nil,
	))

	mu.Lock()
	defaultLogger = newLogger
	mu.Unlock()

	return nil
}

// With returns a new context with additional labels added to the logger.
func With(ctx context.Context, keyvals ...any) context.Context {
	logger := from(ctx)

	return with(ctx, logger.With(keyvals...))
}

// Debug logs a debug message with stack trace.
func Debug(ctx context.Context, msg string, args ...any) {
	logWithStack(ctx, slog.LevelDebug, msg, args...)
}

// Info logs an informational message with stack trace.
func Info(ctx context.Context, msg string, args ...any) {
	logWithStack(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs a warning message with stack trace.
func Warn(ctx context.Context, err error, args ...any) {
	logWithStack(ctx, slog.LevelWarn, err.Error(), args...)
}

// Error logs an error message with call stack trace.
func Error(ctx context.Context, err error, args ...any) {
	logWithStack(ctx, slog.LevelError, err.Error(), args...)
}

// Fatal logs a fatal message with call stack trace.
// It will call os.Exit(1) after logging.
func Fatal(ctx context.Context, err error, args ...any) {
	logWithStack(ctx, slog.LevelError, err.Error(), args...)
	os.Exit(1)
}

type ctxKey uint8

var (
	loggerCtxKey  = ctxKey(0)
	defaultLogger = slog.New(slogdedup.NewOverwriteHandler(slog.NewTextHandler(os.Stdout, nil), nil))
	mu            sync.RWMutex
	// withStackLevel controls which log level and above will include stack traces.
	withStackLevel atomic.Int32
)

func init() {
	// Default to include stack traces for Warn and above.
	SetWithStackLevel(slog.LevelWarn)
}

// SetWithStackLevel sets the minimum log level that will include stack traces.
// It should not be called in init().
func SetWithStackLevel(level slog.Level) {
	withStackLevel.Store(int32(level))
}

// with returns a new context with the given logger.
func with(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// from retrieves the logger from context.
// Returns default logger if context has no logger.
func from(ctx context.Context) *slog.Logger {
	mu.RLock()
	defer mu.RUnlock()
	if ctx == nil {
		return defaultLogger
	}

	if logger, ok := ctx.Value(loggerCtxKey).(*slog.Logger); ok {
		return logger
	}

	return defaultLogger
}

const (
	stackSkip   = 2 // Skip ERROR../logWithStack.
	stackDepth  = 5 // Maximum number of stack frames to capture.
	avgFrameLen = 64
)

func logWithStack(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger := from(ctx)
	if !logger.Enabled(ctx, level) {
		// avoid to get stack trace if logging is disabled for this level
		return
	}

	// Only include stack trace if level is >= withStackLevel
	newArgs := make([]any, 0, len(args)+2)
	newArgs = append(newArgs, args...)
	if level >= slog.Level(withStackLevel.Load()) {
		newArgs = append(newArgs, "stack", getStack(stackSkip, stackDepth))
	}

	logger.Log(ctx, level, msg, newArgs...)
}

// getStack returns a formatted call stack trace.
func getStack(skip, depth int) string {
	pc := make([]uintptr, depth)
	n := runtime.Callers(skip+2, pc) // skip itself and runtime.Callers
	if n == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(n * avgFrameLen)

	frames := runtime.CallersFrames(pc[:n])
	first := true
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if !first {
			b.WriteString(" <- ")
		}
		first = false

		fn := strings.TrimPrefix(frame.Function, model.Module) // no module prefix for zenfeed self.
		b.WriteString(fn)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(frame.Line))
	}

	return b.String()
}
