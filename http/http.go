package http

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type ServerParams struct {
	Handler http.Handler
}

type Server struct {
	websrv *http.Server
}

func New(params *ServerParams) *Server {
	ans := Server{
		websrv: &http.Server{
			Addr:              ":443",
			Handler:           params.Handler,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}

	return &ans
}

func (s *Server) Start(ctx context.Context) error {
	domain := os.Getenv("FIH_DOMAIN")

	var certFile, keyFile string

	if domain == "" {
		certFile = "/app/certs/local-cert.crt"
		keyFile = "/app/certs/local-cert.key"
	} else {
		s.setupAutoTLS([]string{domain})
	}

	return s.websrv.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) setupAutoTLS(domains []string) {
	const defaultCertCache = "/.cache/.certs"

	autoTLSManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(defaultCertCache),
		HostPolicy: autocert.HostWhitelist(domains...),
	}
	// https://ssl-config.mozilla.org/#server=go&version=1.22.0&config=intermediate&guideline=5.7
	s.websrv.TLSConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		GetCertificate: autoTLSManager.GetCertificate,
		NextProtos: []string{
			"h2", "http/1.1", // enable HTTP/2
			acme.ALPNProto, // enable tls-alpn ACME challenges
		},
	}
}
