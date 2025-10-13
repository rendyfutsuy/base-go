package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go.git/utils"
	"github.com/rendyfutsuy/base-go.git/utils/services"
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
// No parameters.
// Returns an error if the scheduler encounters any issues.
func RunResetPasswordEmailScheduler() error {
	utils.InitConfig("config.json")
	utils.InitializedLogger()
	log.Println("Starting scheduler")

	// Initialize the Redis client
	redisSetting := asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	}

	// Initialize the Asynq config
	config := asynq.Config{
		// Specify how many concurrent workers to use
		Concurrency: 10,
		// Optionally specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		// See the godoc for other configuration options
	}

	// Initialize the email service
	emailService, _ := services.NewEmailService()

	// Initialize the Asynq server
	srv := asynq.NewServer(
		redisSetting,
		config,
	)

	// Create a mux to register task handlers
	mux := asynq.NewServeMux()

	mux.HandleFunc(TypeEmailDelivery, func(ctx context.Context, t *asynq.Task) error {
		return HandleEmailResetPasswordRequestTask(ctx, t, emailService)
	})

	// Initialize the scheduler
	scheduler := asynq.NewScheduler(redisSetting, nil)

	// Run the Asynq server in a separate goroutine
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run asynq server: %v", err)
		}
	}()

	// Start the scheduler
	if err := scheduler.Run(); err != nil {
		log.Fatalf("could not run scheduler: %v", err)
		return err
	}

	defer srv.Stop()

	return nil
}

// creates a new task for resetting the user's password via email.
//
// Parameters:
// - userID: The unique identifier of the user.
// - email: The user's email address.
// - session: The session for password reset.
//
// Returns:
// - *asynq.Task: The task to reset the password.
// - error: An error if task creation fails.
func NewEmailResetPasswordRequestTask(userID uuid.UUID, email, session string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{
		UserID:  userID,
		Email:   email,
		Session: session,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
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
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%d, email=%s", p.UserID, p.Email)

	if err := emailService.SendPasswordResetEmail(p.Email, p.Session); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
