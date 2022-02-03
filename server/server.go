package server

import (
	"crypto/tls"
	"net/http"
	"time"
)

// New returns a https secured server
func New(mux *http.ServeMux, address string) *http.Server {
	// generating certs and config from
	// https://gist.github.com/denji/12b3a568f092ab951456#perfect-ssl-labs-score-with-go
	// perfect ssl labs score config
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		// use Go's default preferences which are tuned to avoid attacks
		// Does nothing on clients.
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:      address,
		Handler:   mux,
		TLSConfig: tlsConfig,
		// ReadTimeout is the time from when the connection is accepted to request body fully read
		// it doesn't allow client more time to stream the body of a request
		ReadTimeout: 5 * time.Second,
		// WriteTimeout is time from when the end of request header is read
		// to the end of the response write
		WriteTimeout: 10 * time.Second,
		// IdleTimeout limits the amount of time the connection is alive on
		// the server-side
		IdleTimeout:  120 * time.Second,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return srv
}
