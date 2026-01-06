package schema

import (
	"testing"
)

func TestNewTask(t *testing.T) {
	payload := map[string]interface{}{
		"url": "https://example.com/image.jpg",
	}
	
	task := NewTask("resize_image", payload)
	
	if task.TaskID == "" {
		t.Error("TaskID should not be empty")
	}
	
	if task.TraceID == "" {
		t.Error("TraceID should not be empty")
	}
	
	if task.Type != "resize_image" {
		t.Errorf("Expected type 'resize_image', got '%s'", task.Type)
	}
}

func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name    string
		task    *Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: &Task{
				TaskID:  "123e4567-e89b-12d3-a456-426614174000",
				Type:    "resize_image",
				Payload: map[string]interface{}{"url": "https://example.com/img.jpg"},
				TraceID: "123e4567-e89b-12d3-a456-426614174001",
			},
			wantErr: false,
		},
		{
			name: "invalid task type",
			task: &Task{
				TaskID:  "123e4567-e89b-12d3-a456-426614174000",
				Type:    "unknown_type",
				Payload: map[string]interface{}{"url": "https://example.com/img.jpg"},
				TraceID: "123e4567-e89b-12d3-a456-426614174001",
			},
			wantErr: true,
		},
		{
			name: "missing payload",
			task: &Task{
				TaskID:  "123e4567-e89b-12d3-a456-426614174000",
				Type:    "resize_image",
				TraceID: "123e4567-e89b-12d3-a456-426614174001",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONSerialization(t *testing.T) {
	original := NewTask("resize_image", map[string]interface{}{
		"url": "https://example.com/image.jpg",
	})

	// Serialize
	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Deserialize
	decoded, err := FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	if decoded.TaskID != original.TaskID {
		t.Errorf("TaskID mismatch: got %s, want %s", decoded.TaskID, original.TaskID)
	}
}

func TestValidatePayloadForType(t *testing.T) {
	tests := []struct {
		name     string
		taskType string
		payload  map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "valid resize_image",
			taskType: "resize_image",
			payload: map[string]interface{}{
				"url":    "https://example.com/image.jpg",
				"width":  800,
				"height": 600,
			},
			wantErr: false,
		},
		{
			name:     "invalid resize_image - missing url",
			taskType: "resize_image",
			payload: map[string]interface{}{
				"width": 800,
			},
			wantErr: true,
		},
		{
			name:     "invalid task type",
			taskType: "unknown",
			payload: map[string]interface{}{
				"url": "https://example.com/image.jpg",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePayloadForType(tt.taskType, tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePayloadForType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

