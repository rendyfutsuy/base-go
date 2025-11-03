package test

import (
	"github.com/rendyfutsuy/base-go/modules/role_management/usecase"
)

// setupTestLogger initializes a no-op logger for testing
// This prevents nil pointer panics when Logger is used in usecase code
func setupTestLogger() {
	usecase.SetupTestLogger()
}
