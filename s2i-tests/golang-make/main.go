// Command golang-make is a minimal HTTP/HTTPS server used to verify the Go
// S2I builder/runtime images. It is built as a static binary through its Makefile (`make build`).
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

const appName = "golang-make"

// newMux registers the application routes.
// /healthz follows the Kubernetes health endpoint convention.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Hello from " + appName + "\n"))
	})
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

// loadPKCS12TLSConfig loads a TLS certificate chain and private key from a
// PKCS#12 (.p12 / .pfx) keystore file.
func loadPKCS12TLSConfig(keystorePath, password string) (*tls.Config, error) {
	p12Data, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, fmt.Errorf("read keystore %q: %w", keystorePath, err)
	}

	privateKey, cert, caCerts, err := pkcs12.DecodeChain(p12Data, password)
	if err != nil {
		return nil, fmt.Errorf("decode pkcs12 keystore: %w", err)
	}

	certChain := make([][]byte, 0, 1+len(caCerts))
	certChain = append(certChain, cert.Raw)
	for _, ca := range caCerts {
		certChain = append(certChain, ca.Raw)
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		Certificates: []tls.Certificate{{
			Certificate: certChain,
			PrivateKey:  privateKey,
			Leaf:        cert,
		}},
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tlsPort := os.Getenv("TLS_PORT")
	if tlsPort == "" {
		tlsPort = "8443"
	}

	handler := newMux()
	servers := []*http.Server{
		{
			Addr:              ":" + port,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}

	if keystorePath := os.Getenv("KEYSTORE_P12"); keystorePath != "" {
		tlsCfg, err := loadPKCS12TLSConfig(keystorePath, os.Getenv("KEYSTORE_PASSWORD"))
		if err != nil {
			log.Fatalf("failed to configure TLS from KEYSTORE_P12: %v", err)
		}
		httpsSrv := &http.Server{
			Addr:              ":" + tlsPort,
			Handler:           handler,
			TLSConfig:         tlsCfg,
			ReadHeaderTimeout: 5 * time.Second,
		}
		servers = append(servers, httpsSrv)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, srv := range servers {
		go func(s *http.Server) {
			if s.TLSConfig != nil {
				log.Printf("%s listening on https port %s", appName, tlsPort)
				if err := s.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("https server error: %v", err)
				}
				return
			}
			log.Printf("%s listening on http port %s", appName, port)
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("http server error: %v", err)
			}
		}(srv)
	}

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, srv := range servers {
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error (%s): %v", srv.Addr, err)
		}
	}
}
