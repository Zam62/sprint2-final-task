package agent

import (
	"encoding/json"
	"go-final-task/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgent_ProcessTask(t *testing.T) {
	tests := []struct {
		name      string
		task      models.Task
		wantErr   bool
		wantValue float64
	}{
		{
			name: "сложение",
			task: models.Task{
				ID:            "1",
				Arg1:          1,
				Arg2:          4,
				Operation:     "+",
				OperationTime: 10,
			},
			wantErr:   false,
			wantValue: 5,
		},
		{
			name: "вычитание",
			task: models.Task{
				ID:            "2",
				Arg1:          9,
				Arg2:          7,
				Operation:     "-",
				OperationTime: 10,
			},
			wantErr:   false,
			wantValue: 2,
		},
		{
			name: "умножение",
			task: models.Task{
				ID:            "3",
				Arg1:          3,
				Arg2:          2,
				Operation:     "*",
				OperationTime: 10,
			},
			wantErr:   false,
			wantValue: 6,
		},
		{
			name: "деление",
			task: models.Task{
				ID:            "4",
				Arg1:          8,
				Arg2:          4,
				Operation:     "/",
				OperationTime: 10,
			},
			wantErr:   false,
			wantValue: 2,
		},
		{
			name: "деление на ноль",
			task: models.Task{
				ID:            "5",
				Arg1:          7,
				Arg2:          0,
				Operation:     "/",
				OperationTime: 10,
			},
			wantErr: true,
		},
		{
			name: "неизвестная операция",
			task: models.Task{
				ID:            "6",
				Arg1:          5,
				Arg2:          8,
				Operation:     "%",
				OperationTime: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/internal/task" && r.Method == http.MethodPost {
					var result models.TaskResult
					json.NewDecoder(r.Body).Decode(&result)

					if !tt.wantErr {
						if result.Result != tt.wantValue {
							t.Errorf("Expected result %f, got %f", tt.wantValue, result.Result)
						}
					}

					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"status": "success"})
				}
			}))
			defer server.Close()

			agent := NewAgent(server.URL)
			err := agent.processTask(tt.task)

			if (err != nil) != tt.wantErr {
				t.Errorf("processTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_GetTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" && r.Method == http.MethodGet {
			Task := models.Task{
				ID:            "1",
				Arg1:          2,
				Arg2:          3,
				Operation:     "+",
				OperationTime: 10,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.TaskResponse{Task: &Task})
		}
	}))
	defer server.Close()

	agent := NewAgent(server.URL)
	task, err := agent.getTask()

	if err != nil {
		t.Errorf("getTask() error = %v", err)
	}

	if task == nil {
		t.Errorf("getTask() returned nil task")
	} else {
		if task.ID != "1" || task.Arg1 != 2 || task.Arg2 != 3 || task.Operation != "+" {
			t.Errorf("getTask() returned unexpected task: %+v", task)
		}
	}
}

func TestAgent_SubmitResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" && r.Method == http.MethodPost {
			var result models.TaskResult
			json.NewDecoder(r.Body).Decode(&result)

			if result.ID != "1" || result.Result != 5 {
				t.Errorf("Unexpected result: %+v", result)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		}
	}))
	defer server.Close()

	agent := NewAgent(server.URL)
	err := agent.submitResult(models.TaskResult{
		ID:     "1",
		Result: 5,
	})

	if err != nil {
		t.Errorf("submitResult() error = %v", err)
	}
}
