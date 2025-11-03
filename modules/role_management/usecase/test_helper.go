package usecase

import (
	"time"

	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/role_management"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewTestRoleUsecase creates a new roleUsecase instance for testing purposes
// This allows test packages to create roleUsecase instances with custom configurations
func NewTestRoleUsecase(r role_management.Repository, a auth.Repository, timeout time.Duration) *roleUsecase {
	return &roleUsecase{
		roleRepo:       r,
		authRepo:       a,
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
