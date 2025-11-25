package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/rangira25/notification/internal/services"
	"github.com/rangira25/notification/internal/tasks"
)

type Processor struct {
	Redis    *redis.Client
	EmailSvc *services.EmailService
}

func NewProcessor(r *redis.Client, es *services.EmailService) *Processor {
	return &Processor{Redis: r, EmailSvc: es}
}

func (p *Processor) HandleWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var payload tasks.WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// idempotency at worker: ensure task not executed twice
	taskID := t.ResultWriter().TaskID()
	if taskID == "" {
		taskID = "task:" + payload.ReqID
	}
	processedKey := "processed:" + taskID
	ok, err := p.Redis.SetNX(ctx, processedKey, "1", 48*time.Hour).Result() // 48h TTL
	if err != nil {
		// Redis error -> fail to trigger retry
		return err
	}
	if !ok {
		// already processed
		log.Printf("task %s already processed, skip\n", taskID)
		return nil
	}

	// Do the actual work (send email)
	if err := p.EmailSvc.SendWelcome(ctx, payload.Email, "Welcome!", "Welcome to our service"); err != nil {
		// remove processed key so retry can happen
		p.Redis.Del(ctx, processedKey)
		return err
	}

	log.Printf("sent welcome email to %s for user %s\n", payload.Email, payload.UserID)
	return nil
}
