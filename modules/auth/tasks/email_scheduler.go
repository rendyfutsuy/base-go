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

	redisSetting := asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	}

	config := asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}

	emailService, _ := services.NewEmailService()

	srv := asynq.NewServer(redisSetting, config)
	mux := asynq.NewServeMux()

	// Register all handlers
	mux.HandleFunc(TypeEmailDelivery, func(ctx context.Context, t *asynq.Task) error {
		return HandleEmailResetPasswordRequestTask(ctx, t, emailService)
	})
	RegisterVerificationEmailHandler(mux, emailService)

	scheduler := asynq.NewScheduler(redisSetting, nil)

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
