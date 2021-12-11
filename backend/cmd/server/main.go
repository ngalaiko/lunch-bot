package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lunch/pkg/http"
	"lunch/pkg/lunch"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

var (
	placesStore = storage_places.NewMemory()
	boostsStore = storage_boosts.NewMemory()
	rollsStore  = storage_rolls.NewMemory()
)

func init() {
	for i := 0; i < 10; i++ {
		placesStore.Store(context.Background(), &places.Place{
			ID:      places.ID(fmt.Sprint(i)),
			Name:    places.Name(fmt.Sprintf("Place %d", i)),
			AddedAt: time.Now().Add(-1 * time.Duration(i) * time.Hour),
		})
	}
}

var (
	addr = flag.String("addr", ":8000", "http listen address")
)

func main() {
	flag.Parse()

	cfg := &http.Configuration{}
	if err := cfg.Parse(); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	roller := lunch.New(placesStore, boostsStore, rollsStore)
	srv := http.NewServer(cfg, roller)

	// Wait for shut down in a separate goroutine.
	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal, 1)
		signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
		sig := <-shutdownCh

		log.Printf("[INFO] received %s, shutting down", sig)

		shutdownTimeout := 15 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		errCh <- srv.Shutdown(shutdownCtx)
	}()

	if err := srv.ListenAndServe(*addr); err != nil {
		log.Printf("[ERROR] http server: %s", err)
	}

	// Handle shutdown errors.
	if err := <-errCh; err != nil {
		log.Printf("[ERROR] error during shutdown: %s", err)
	}

	log.Printf("[INFO] application stopped")
}
