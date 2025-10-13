package main

import (
	"log"

	auth "github.com/rendyfutsuy/base-go.git/modules/auth/tasks"
)

// workers.go
// This file is for running the scheduler, that will executes tasks on queue
func main() {
	// module auth
	// Start The Scheduler for sending reset password email
	err := auth.RunResetPasswordEmailScheduler()

	if err != nil {
		log.Fatal(err)
	}
	// End The Scheduler for sending reset password email
}
