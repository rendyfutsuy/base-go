package main

import (
	"log"

	"github.com/rendyfutsuy/base-go/modules/auth/tasks"
)

func main() {
	if err := tasks.RunEmailScheduler(); err != nil {
		log.Fatalf("email worker failed: %v", err)
	}
}
