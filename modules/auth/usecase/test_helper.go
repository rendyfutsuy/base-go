package usecase

import (
	"time"

	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewTestAuthUsecase creates a new authUsecase instance for testing purposes
// This allows test packages to create authUsecase instances with custom configurations
func NewTestAuthUsecase(r auth.Repository, timeout time.Duration, hashSalt string, signingKey []byte, expireDuration time.Duration) *authUsecase {
	return &authUsecase{
		authRepo:       r,
		contextTimeout: timeout,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

// SetupTestLogger initializes a no-op logger for testing
// This prevents nil pointer panics when Logger is used in usecase code
func SetupTestLogger() {
	if utils.Logger == nil {
		core := zapcore.NewNopCore()
		utils.Logger = zap.New(core)
	}
}
