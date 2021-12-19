package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lunch/pkg/http"
	"lunch/pkg/jwt"
	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	"lunch/pkg/lunch"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

var (
	placesStore  = storage_places.NewMemory()
	boostsStore  = storage_boosts.NewMemory()
	rollsStore   = storage_rolls.NewMemory()
	jwtKeysStore = storage_jwt_keys.NewMemory()
)

var (
	roller     = lunch.New(placesStore, boostsStore, rollsStore)
	jwtService = jwt.NewService(jwtKeysStore)
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
	addr    = flag.String("addr", ":8000", "http listen address")
	tlsCert = flag.String("tls-cert", ".cert/cert.pem", "path to TLS certificate")
	tlsKey  = flag.String("tls-key", ".cert/key.pem", "path to TLS key")
)

func main() {
	flag.Parse()

	cfg := &http.Configuration{}
	if err := cfg.Parse(); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	srv := http.NewServer(cfg, roller, jwtService)

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

	var certificates []tls.Certificate
	if *tlsCert != "" {
		cert, err := loadTLSCert(*tlsCert, *tlsKey)
		if err != nil {
			log.Fatalf("failed to load TLS certificate: %v", err)
		}
		certificates = append(certificates, cert)
	}

	if err := srv.ListenAndServe(*addr, certificates...); err != nil {
		log.Printf("[ERROR] http server: %s", err)
	}

	// Handle shutdown errors.
	if err := <-errCh; err != nil {
		log.Printf("[ERROR] error during shutdown: %s", err)
	}

	log.Printf("[INFO] application stopped")
}

func loadTLSCert(certPath, keyPath string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load TLS certificate: %v", err)
	}

	return cert, nil
}
