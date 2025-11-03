package utils

import (
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func InitializeNewRelic() *newrelic.Application {
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(ConfigVars.String("newrelic.app_name")),
		newrelic.ConfigLicense(ConfigVars.String("newrelic.license")),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize New Relic application: %v", err))
	}

	// Logger.Info("New Relic application initialized.")
	return newRelicApp
}
