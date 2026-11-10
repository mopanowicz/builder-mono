package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if body := strings.TrimSpace(rec.Body.String()); body != `{"status":"UP"}` {
		t.Errorf("body = %s, want {\"status\":\"UP\"}", body)
	}
}

func TestUnknownPath(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestPKCS12TLSServer(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &privKey.PublicKey, privKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
	}

	const password = "test-password"
	p12Bytes, err := pkcs12.Modern2023.Encode(privKey, cert, nil, password)
	if err != nil {
		t.Fatalf("encode pkcs12: %v", err)
	}

	keystorePath := filepath.Join(t.TempDir(), "keystore.p12")
	if err := os.WriteFile(keystorePath, p12Bytes, 0o600); err != nil {
		t.Fatalf("write keystore: %v", err)
	}

	tlsCfg, err := loadPKCS12TLSConfig(keystorePath, password)
	if err != nil {
		t.Fatalf("loadPKCS12TLSConfig: %v", err)
	}

	if _, err := loadPKCS12TLSConfig(keystorePath, "wrong-password"); err == nil {
		t.Fatal("expected error with wrong keystore password, got nil")
	}

	srv := httptest.NewUnstartedServer(newMux())
	srv.TLS = tlsCfg
	srv.StartTLS()
	defer srv.Close()

	rootPool := x509.NewCertPool()
	rootPool.AddCert(cert)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: rootPool},
		},
	}

	resp, err := client.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz over HTTPS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
