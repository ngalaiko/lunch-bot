package http

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"lunch/pkg/http/slack"
	"lunch/pkg/lunch"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

type Server struct {
	handler    *handler
	roller     *lunch.Roller
	httpServer *http.Server
}

func NewServer(
	boostsStore storage_boosts.Storage,
	placesStore storage_places.Storage,
	rollsStore storage_rolls.Storage,
) *Server {
	s := &Server{
		handler: newHandler(),
		roller:  lunch.New(placesStore, boostsStore, rollsStore),
	}
	s.registerRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.handler.GET("/", slack.NewHandler(s.roller).ServeHTTP)
}

func (s *Server) ListenAndServe(addr string, certs ...tls.Certificate) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig: &tls.Config{
			Certificates:     certs,
			NextProtos:       []string{"h2", "http/1.1"},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
		},
	}
	if len(certs) > 0 {
		log.Printf("[INFO] listening https on %s", addr)
	} else {
		log.Printf("[INFO] listening http on %s", addr)
	}
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("[INFO] stopping http server")
	defer log.Printf("[INFO] http server stopped")
	return s.httpServer.Shutdown(ctx)
}
