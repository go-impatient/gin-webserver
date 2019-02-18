package log

import (
	"github.com/sevennt/wzap"
)

// Debug logs debug level messages with default logger.
func Debug(msg string, args ...interface{}) {
	wzap.Debug(msg, args...)
}

// Debugf logs debug level messages with default logger in printf-style.
func Debugf(msg string, args ...interface{}) {
	wzap.Debugf(msg, args...)
}

// Info logs Info level messages with default logger in structured-style.
func Info(msg string, args ...interface{}) {
	wzap.Info(msg, args...)
}

// Infof logs Info level messages with default logger in printf-style.
func Infof(msg string, args ...interface{}) {
	wzap.Infof(msg, args...)
}

// Warn logs Warn level messages with default logger in structured-style.
func Warn(msg string, args ...interface{}) {
	wzap.Warn(msg, args...)
}

// Warnf logs Warn level messages with default logger in printf-style.
func Warnf(msg string, args ...interface{}) {
	wzap.Warnf(msg, args...)
}

// Error logs Error level messages with default logger in structured-style.
func Error(msg string, args ...interface{}) {
	wzap.Error(msg, args...)
}

// Errorf logs Error level messages with default logger in printf-style.
func Errorf(msg string, args ...interface{}) {
	wzap.Errorf(msg, args...)
}

// Panic logs Panic level messages with default logger in structured-style.
func Panic(msg string, args ...interface{}) {
	wzap.Panic(msg, args...)
}

// Panicf logs Panicf level messages with default logger in printf-style.
func Panicf(msg string, args ...interface{}) {
	wzap.Panicf(msg, args...)
}

// Fatal logs Fatal level messages with default logger in structured-style.
func Fatal(msg string, args ...interface{}) {
	wzap.Fatal(msg, args...)
}

// Fatalf logs Fatalf level messages with default logger in printf-style.
func Fatalf(msg string, args ...interface{}) {
	wzap.Fatalf(msg, args...)
}
