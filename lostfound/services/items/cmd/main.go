package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/lost-and-found/items/pkg"
)

func main() {
	svc, err := pkg.NewItemService()
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create a context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start (blocks until shutdown or error)
	if err := svc.Start(ctx); err != nil {
		// start returned an error (listen error or shutdown error)
		log.Fatalf("service stopped with error: %v", err)
	}

	// optionally wait a moment for logs to flush
	time.Sleep(100 * time.Millisecond)
}
