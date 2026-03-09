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
	TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
	UserID  uuid.UUID
	Email   string
	Session string
}

// runs the email scheduler for resetting passwords.
//
// ctx: the context.Context for the task
// t: the asynq.Task containing the payload
// emailService: the services.EmailService for sending the password reset email
// Returns an error if there is an issue with the scheduler.
func HandleEmailResetPasswordRequestTask(ctx context.Context, t *asynq.Task, emailService *services.EmailService) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		utils.Logger.Error(err.Error())
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%d, email=%s", p.UserID, p.Email)

	if err := emailService.SendPasswordResetEmail(p.Email, p.Session); err != nil {
		utils.Logger.Error(err.Error())
		return fmt.Errorf("failed to send email: %v", err)
	}
	utils.Logger.Info(fmt.Sprintf("Password reset email sent successfully: user_id=%s, email=%s", p.UserID.String(), p.Email))

	return nil
}
