package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

// Runs HTTP3 Server on UDP
func RunQuicHTTP3Server(certificateFilePath, keyFilePath string, mux *http.ServeMux) {
	fmt.Println("Running Quic HTTP3 Server on UDP")
	err := http3.ListenAndServeQUIC(":443", certificateFilePath, keyFilePath, mux)
	if err != nil {
		fmt.Println("Failed to run HTTP3 Server on UDP address with error :", err)
	}
}

// Runs HTTP3 Server on UDP with provided configuration
func RunQuicHTTP3ConfigServer(mux *http.ServeMux, tlsConfig *tls.Config, quicConfig *quic.Config) {
	server := http3.Server{
		Handler:    mux,
		Addr:       ":443",
		TLSConfig:  http3.ConfigureTLSConfig(tlsConfig),
		QUICConfig: quicConfig,
	}
	fmt.Println("Running Quic HTTP3 Server on UDP")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to run HTTP3 Server on UDP address with error :", err)
	}
}

// Runs HTTP3 Server on both TLS/TCP and QUIC connections in parallel.
// It returns if one of the two returns an error.
func RunHTTP3ServerForTlsTcpAndQuic(certificateFilePath, keyFilePath string, mux *http.ServeMux) {
	fmt.Println("Running Quic HTTP3 Server on both TLS/TCP and QUIC/UDP")
	err := http3.ListenAndServeTLS(":443", certificateFilePath, keyFilePath, mux)
	if err != nil {
		fmt.Println("Failed to run HTTP3 Server on TLS/TCP and QUIC connections in parallel with error :", err)
	}
}

// Graceful Shutdown of HTTP3 Server
func HTTP3ServerShutdown(server *http3.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	server.Shutdown(ctx)
}
