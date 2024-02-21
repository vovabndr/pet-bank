package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "pet-bank/db/sqlc"
	"pet-bank/utils"
)
import "context"

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskPayloadSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %s", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %s", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retries", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: utils.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	to := []string{user.Email}
	subject := "Welcome to pet-bank"
	verifyUrl := fmt.Sprintf(
		"http://%s/v1/verify_email?email_id=%d&secret_code=%s",
		processor.config.VerifyUrl,
		verifyEmail.ID,
		verifyEmail.SecretCode,
	)
	content := fmt.Sprintf(`Hello %s,<br/>
	Please <a href="%s">click here</a> to verify your email address. 
	`, user.Username, verifyUrl)
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)

	if err != nil {
		return fmt.Errorf("faild to sent verify email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")

	return nil
}
