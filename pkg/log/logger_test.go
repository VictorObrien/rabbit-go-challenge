package log

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	// Capturar output
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	testLogger := slog.New(handler)
	SetLogger(testLogger)

	// Testar log
	Info("test message", "key", "value")

	// Verificar output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["msg"] != "test message" {
		t.Errorf("Expected msg 'test message', got '%v'", logEntry["msg"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("Expected key 'value', got '%v'", logEntry["key"])
	}
}

func TestTaskLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	testLogger := slog.New(handler)
	SetLogger(testLogger)

	taskLogger := TaskLogger("task-123", "trace-456")
	taskLogger.Info("processing task")

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["task_id"] != "task-123" {
		t.Errorf("Expected task_id 'task-123', got '%v'", logEntry["task_id"])
	}

	if logEntry["trace_id"] != "trace-456" {
		t.Errorf("Expected trace_id 'trace-456', got '%v'", logEntry["trace_id"])
	}
}

