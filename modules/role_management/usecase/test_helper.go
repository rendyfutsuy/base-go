package usecase

import (
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// setupTestLogger initializes a no-op logger for testing
// This prevents nil pointer panics when Logger is used in usecase code
func setupTestLogger() {
	if utils.Logger == nil {
		core := zapcore.NewNopCore()
		utils.Logger = zap.New(core)
	}
}
