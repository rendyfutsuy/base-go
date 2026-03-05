package tasks

import (
	"encoding/json"
	"log"

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
	workers := map[string]func([]byte) error{
		TypeEmailDelivery: func(body []byte) error {
			var p EmailDeliveryPayload
			if err := json.Unmarshal(body, &p); err != nil {
				return err
			}
			return emailService.SendPasswordResetEmail(p.Email, p.Session)
		},
		TypeEmailVerification: func(body []byte) error {
			var p VerificationEmailPayload
			if err := json.Unmarshal(body, &p); err != nil {
				return err
			}
			return emailService.SendVerificationEmail(p.Email, p.Code)
		},
	}
	if err := q.Run(workers); err != nil {
		return err
	}

	return nil
}
