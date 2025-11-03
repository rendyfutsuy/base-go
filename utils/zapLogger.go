package utils

import (
	"fmt"
	"os"
	"time"

	// "github.com/getsentry/sentry-go"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// SentryWriter is a custom io.Writer implementation that forwards logs to Sentry.
// type SentryWriter struct{}

// Write forwards the log message to Sentry.
// func (sw SentryWriter) Write(p []byte) (n int, err error) {
// 	sentry.CaptureMessage(string(p))
// 	return len(p), nil
// }

func InitializedLogger(newRelicApp *newrelic.Application) {
	serverID := uuid.New().String()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, os.ModePerm); err != nil {
			panic(fmt.Sprintf("Failed to create logs directory: %s", err))
		}
	}
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("%s/%s_%s.log", logDir, timestamp, serverID)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %s", err))
	}

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.ErrorLevel),
		// sentryCore,
	)

	if ConfigVars.Bool("newrelic.enable_new_relic_logging") {
		backgroundCore, err := nrzap.WrapBackgroundCore(core, newRelicApp)
		if err != nil {
			if err == nrzap.ErrNilApp {
				panic("New Relic Application is not initialized")
			}
			panic("Failed to wrap Zap core with New Relic integration: " + err.Error())
		}
		Logger = zap.New(backgroundCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		Logger = Logger.With(zap.String("identifier", ConfigVars.String("newrelic.identifier")))
		Logger.Info("Logger connected to new relic")

	} else {
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		Logger.Info("Logger not connected to new relic")
	}
	Logger = Logger.With(zap.String("server_id", serverID))
	Logger.Info("Logger initialized with server ID", zap.String("server_id", serverID))
}
