package log

import (
	"context"
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	// Logger JSON estruturado
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	defaultLogger = slog.New(handler)
}

// SetLogger permite configurar um logger customizado
func SetLogger(logger *slog.Logger) {
	defaultLogger = logger
}

// GetLogger retorna o logger padrão
func GetLogger() *slog.Logger {
	return defaultLogger
}

// Info loga mensagem de info
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Error loga mensagem de erro
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// Warn loga mensagem de warning
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Debug loga mensagem de debug
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// WithContext retorna um logger com contexto
func WithContext(ctx context.Context) *slog.Logger {
	return defaultLogger.With()
}

// WithFields retorna um logger com campos adicionais
func WithFields(fields map[string]interface{}) *slog.Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return defaultLogger.With(args...)
}

// TaskLogger cria um logger com campos de contexto de tarefa
func TaskLogger(taskID, traceID string) *slog.Logger {
	return defaultLogger.With(
		"task_id", taskID,
		"trace_id", traceID,
	)
}

// MetricsLogger cria um logger específico para métricas
func MetricsLogger() *slog.Logger {
	return defaultLogger.With("component", "metrics")
}

// WorkerLogger cria um logger específico para workers
func WorkerLogger(workerID int) *slog.Logger {
	return defaultLogger.With(
		"component", "worker",
		"worker_id", workerID,
	)
}
