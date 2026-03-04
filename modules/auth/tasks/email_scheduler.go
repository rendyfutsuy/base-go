package tasks

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
)

// RunEmailScheduler initializes Asynq server and registers all email-related handlers.
//
// It sets up Redis client, configures queues, initializes EmailService,
// registers Reset Password and Verification email handlers, and runs the server & scheduler.
func RunEmailScheduler() error {
	utils.InitConfig("config.json")
	var newRelicApp *newrelic.Application
	if utils.ConfigVars.Exists("newrelic.enable_new_relic_logging") && utils.ConfigVars.Bool("newrelic.enable_new_relic_logging") {
		newRelicApp = utils.InitializeNewRelic()
	}
	utils.InitializedLogger(newRelicApp)
	log.Println("Starting email scheduler")

	emailService, _ := services.NewEmailService()

	q := services.NewQueueService()
	srv, err := q.NewAsynqServer()
	if err != nil {
		return err
	}
	mux := asynq.NewServeMux()

	// Register all handlers
	mux.HandleFunc(TypeEmailDelivery, func(ctx context.Context, t *asynq.Task) error {
		return HandleEmailResetPasswordRequestTask(ctx, t, emailService)
	})
	RegisterVerificationEmailHandler(mux, emailService)

	scheduler, err := q.NewAsynqScheduler()
	if err != nil {
		return err
	}

	// Run server in goroutine
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run asynq server: %v", err)
		}
	}()

	// Start scheduler
	if err := scheduler.Run(); err != nil {
		log.Fatalf("could not run scheduler: %v", err)
		return err
	}

	defer srv.Stop()

	return nil
}
