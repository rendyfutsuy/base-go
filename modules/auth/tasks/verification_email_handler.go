package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
)

const (
	TypeEmailVerification = "email:verification"
)

type VerificationEmailPayload struct {
	UserID uuid.UUID
	Email  string
	Code   string
}

func NewEmailVerificationTask(userID uuid.UUID, email, code string) (*asynq.Task, error) {
	payload, err := json.Marshal(VerificationEmailPayload{
		UserID: userID,
		Email:  email,
		Code:   code,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailVerification, payload), nil
}

func HandleVerificationEmailTask(ctx context.Context, t *asynq.Task, emailService *services.EmailService) error {
	var p VerificationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		utils.Logger.Error(err.Error())
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Verification Email: user_id=%d, email=%s", p.UserID, p.Email)
	if err := emailService.SendVerificationEmail(p.Email, p.Code); err != nil {
		utils.Logger.Error(err.Error())
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	return nil
}

func RegisterVerificationEmailHandler(mux *asynq.ServeMux, emailService *services.EmailService) {
	mux.HandleFunc(TypeEmailVerification, func(ctx context.Context, t *asynq.Task) error {
		return HandleVerificationEmailTask(ctx, t, emailService)
	})
}
