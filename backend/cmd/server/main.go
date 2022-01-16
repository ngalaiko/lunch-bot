package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lunch/pkg/http"
	"lunch/pkg/jwt"
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/events"
	service_users "lunch/pkg/users/service"
)

var (
	roller       = lunch.New(placesStore, boostsStore, rollsStore, events.NewRegistry())
	jwtService   = jwt.NewService(jwtKeysStore)
	usersService = service_users.New(usersStore)
)

var (
	addr = flag.String("addr", ":8000", "http listen address")

	enableTLS = flag.Bool("tls", false, "enable TLS")
	tlsCert   = flag.String("tls-cert", ".cert/cert.pem", "path to TLS certificate")
	tlsKey    = flag.String("tls-key", ".cert/key.pem", "path to TLS key")
)

func main() {
	flag.Parse()

	cfg := &http.Configuration{}
	if err := cfg.Parse(); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	srv := http.NewServer(cfg, roller, jwtService, usersService)

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
	if *enableTLS {
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
