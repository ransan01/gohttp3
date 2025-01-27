package client

import (
	"crypto/tls"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func CreateHTTP3ConfigClient(tlsConfig *tls.Config, quicConfig *quic.Config) *http.Client {
	roundTripper := &http3.Transport{
		TLSClientConfig: http3.ConfigureTLSConfig(tlsConfig),
		QUICConfig:      quicConfig,
	}
	defer roundTripper.Close()
	client := &http.Client{
		Transport: roundTripper,
	}
	return client
}

func CreateHTTP3Client() *http.Client {
	tlsConfig := tls.Config{}
	quicConfig := quic.Config{}
	return CreateHTTP3ConfigClient(&tlsConfig, &quicConfig)
}
