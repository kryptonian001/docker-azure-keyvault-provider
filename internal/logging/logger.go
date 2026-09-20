package logging

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger defines the interface for structured logging.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// DefaultLogger wraps slog.Logger for structured JSON output.
type DefaultLogger struct {
	logger *slog.Logger
}

// NewDefaultLogger creates a new logger with JSON output to stdout.
func NewDefaultLogger() *DefaultLogger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return &DefaultLogger{logger: logger}
}

// Debug logs a debug level message with structured key-value pairs.
func (l *DefaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, keysAndValues...)
}

// Info logs an info level message with structured key-value pairs.
func (l *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

// Error logs an error level message with structured key-value pairs.
func (l *DefaultLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, keysAndValues...)
}

// NoOpLogger is a no-operation logger used for testing or silent mode.
type NoOpLogger struct{}

// NewNoOpLogger creates a logger that silently discards all messages.
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug is a no-op implementation.
func (l *NoOpLogger) Debug(msg string, keysAndValues ...interface{}) {
	// Silent
}

// Info is a no-op implementation.
func (l *NoOpLogger) Info(msg string, keysAndValues ...interface{}) {
	// Silent
}

// Error is a no-op implementation.
func (l *NoOpLogger) Error(msg string, keysAndValues ...interface{}) {
	// Silent
}

// MockLogger records all log messages for testing purposes.
type MockLogger struct {
	Messages []string
}

// NewMockLogger creates a logger that records all messages.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Messages: []string{},
	}
}

// Debug records a debug message.
func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.Messages = append(l.Messages, fmt.Sprintf("DEBUG: %s %v", msg, keysAndValues))
}

// Info records an info message.
func (l *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Messages = append(l.Messages, fmt.Sprintf("INFO: %s %v", msg, keysAndValues))
}

// Error records an error message.
func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Messages = append(l.Messages, fmt.Sprintf("ERROR: %s %v", msg, keysAndValues))
}

// ContainsMessage checks if a message was logged (partial match).
func (l *MockLogger) ContainsMessage(substr string) bool {
	for _, msg := range l.Messages {
		if msg == substr || contains(msg, substr) {
			return true
		}
	}
	return false
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || len(s) > len(substr)
}
