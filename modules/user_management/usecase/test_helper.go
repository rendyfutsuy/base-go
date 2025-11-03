package usecase

import (
	"time"

	auth "github.com/rendyfutsuy/base-go/modules/auth"
	roleManagement "github.com/rendyfutsuy/base-go/modules/role_management"
	"github.com/rendyfutsuy/base-go/modules/user_management"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewTestUserUsecase creates a new userUsecase instance for testing purposes
// This allows test packages to create userUsecase instances with custom configurations
func NewTestUserUsecase(r user_management.Repository, rm roleManagement.Repository, authRepo auth.Repository, timeout time.Duration) *userUsecase {
	return &userUsecase{
		userRepo:       r,
		roleManagement: rm,
		auth:           authRepo,
		contextTimeout: timeout,
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
