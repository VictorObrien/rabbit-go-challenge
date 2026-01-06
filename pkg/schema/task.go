package schema

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Task representa uma tarefa a ser processada
type Task struct {
	TaskID   string                 `json:"task_id" validate:"required,uuid"`
	Type     string                 `json:"type" validate:"required,oneof=resize_image process_video generate_thumbnail"`
	Payload  map[string]interface{} `json:"payload" validate:"required"`
	TraceID  string                 `json:"trace_id" validate:"required,uuid"`
	CreateAt time.Time              `json:"created_at"`
}

// TaskPayload representa payloads específicos por tipo
type ResizeImagePayload struct {
	URL    string `json:"url" validate:"required,url"`
	Width  int    `json:"width,omitempty" validate:"omitempty,min=1,max=4096"`
	Height int    `json:"height,omitempty" validate:"omitempty,min=1,max=4096"`
}

type ProcessVideoPayload struct {
	URL    string `json:"url" validate:"required,url"`
	Format string `json:"format,omitempty" validate:"omitempty,oneof=mp4 webm"`
}

// NewTask cria uma nova tarefa com valores padrão
func NewTask(taskType string, payload map[string]interface{}) *Task {
	return &Task{
		TaskID:   uuid.New().String(),
		Type:     taskType,
		Payload:  payload,
		TraceID:  uuid.New().String(),
		CreateAt: time.Now().UTC(),
	}
}

// Validate valida a estrutura da tarefa
func (t *Task) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

// ToJSON serializa a tarefa para JSON
func (t *Task) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// FromJSON deserializa JSON para Task
func FromJSON(data []byte) (*Task, error) {
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// ValidatePayloadForType valida o payload específico baseado no tipo
func ValidatePayloadForType(taskType string, payload map[string]interface{}) error {
	validate := validator.New()
	
	// Converter map para bytes e depois para struct específica
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	switch taskType {
	case "resize_image":
		var p ResizeImagePayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			return err
		}
		return validate.Struct(p)
	case "process_video":
		var p ProcessVideoPayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			return err
		}
		return validate.Struct(p)
	case "generate_thumbnail":
		// Thumbnail usa mesmo payload que resize_image
		var p ResizeImagePayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			return err
		}
		return validate.Struct(p)
	default:
		return ErrInvalidTaskType
	}
}

// Erros customizados
var (
	ErrInvalidTaskType = &ValidationError{Message: "invalid task type"}
	ErrInvalidPayload  = &ValidationError{Message: "invalid payload"}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
