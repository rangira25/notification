package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeSendWelcomeEmail = "email:welcome"
)

type WelcomeEmailPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	ReqID  string `json:"req_id"` // client-provided idempotency
}

func NewWelcomeEmailTask(p *WelcomeEmailPayload) (*asynq.Task, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendWelcomeEmail, b), nil
}

// Enqueue helper with idempotency via TaskID
func EnqueueWelcomeEmail(client *asynq.Client, p *WelcomeEmailPayload) (*asynq.TaskInfo, error) {
	task, err := NewWelcomeEmailTask(p)
	if err != nil {
		return nil, err
	}
	return client.Enqueue(task,
		asynq.Queue("default"),
		asynq.MaxRetry(5),
		asynq.Timeout(30*time.Second),
		asynq.TaskID(p.ReqID),
	)
}
