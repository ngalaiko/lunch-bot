package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"lunch/pkg/jwt"
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/events"
	service_users "lunch/pkg/users/service"
)

type Server struct {
	handler    http.Handler
	httpServer *http.Server
}

func NewServer(
	cfg *Configuration,
	roller *lunch.Roller,
	jwtService *jwt.Service,
	usersService *service_users.Service,
	eventsRegistry *events.Registry,
) *Server {
	return &Server{
		handler: NewHandler(cfg, roller, jwtService, usersService, eventsRegistry),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string, certs ...tls.Certificate) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	if len(certs) > 0 {
		tlsConfig := &tls.Config{
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
		}
		s.httpServer.TLSConfig = tlsConfig
		ln = tls.NewListener(ln, tlsConfig)
		log.Printf("[INFO] listening https on %s", addr)
	} else {
		log.Printf("[INFO] listening http on %s", addr)
	}
	if err := s.httpServer.Serve(ln); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("[INFO] stopping http server")
	defer log.Printf("[INFO] http server stopped")
	return s.httpServer.Shutdown(ctx)
}
