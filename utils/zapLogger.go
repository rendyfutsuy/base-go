package utils

import (
	"fmt"
	"os"

	// "github.com/getsentry/sentry-go"
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

func InitializedLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel

	fmt.Println("config vars")
	fmt.Println("local")
	// newUUID := uuid.New()

	//sentry
	// err := sentry.Init(sentry.ClientOptions{
	// 	Dsn:         ConfigVars.String("sentry_dsn"),
	// 	Environment: ConfigVars.String("sentry_env") + "-" + newUUID.String(),
	// 	SampleRate:  1,
	// })
	// if err != nil {
	// 	Logger.Error("Error initializing Sentry", zap.Error(err))
	// }

	// sentryWriter := SentryWriter{}

	// sentryCore := zapcore.NewCore(
	// 	zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	// 	zapcore.AddSync(sentryWriter),
	// 	zapcore.InfoLevel,
	// )

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		// sentryCore,
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
