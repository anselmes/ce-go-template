// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CloudEventConfig struct {
  Address string `envconfig:"CE_ADDRESS" default:"localhost"`
  Certificate string `envconfig:"CE_CERT" default:"tls-bundle.pem"`
  CertificateKey string `envconfig:"CE_KEY" default:"tls-key.pem"`
  Insecure bool `envconfig:"CE_INSECURE" default:"false"`
  Port    int `envconfig:"CE_PORT" default:"8080"`
  SkipVerify bool `envconfig:"CE_SKIP_VERIFY" default:"false"`
  Config *tls.Config
}

func (config CloudEventConfig) Client() (cloudevents.Client, error) {
  var client cloudevents.Client

  if config.Insecure {
    log.Printf("Insecure mode enabled, skipping TLS verification")

    // Create protocol and client for insecure mode
    protocol, err := cloudevents.NewHTTP(cloudevents.WithTarget(config.Url().String()))
    if err != nil {
      return nil, Error(ErrUnknown, err.Error())
    }
    client, err = cloudevents.NewClient(protocol, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
    if err != nil {
      return nil, Error(ErrUnknown, err.Error())
    }
  } else {
    pool := x509.NewCertPool()

    // Configure a new http.Transport with TLS
    cert, err := tls.LoadX509KeyPair(config.Certificate, config.CertificateKey)
    if err != nil {
      return nil, Error(ErrTlsConfig, err.Error())
    }

    ca, err := os.ReadFile(config.Certificate)
    if err != nil {
      return nil, Error(ErrTlsConfig, err.Error())
    }

    pool.AppendCertsFromPEM(ca)

    config.Config = &tls.Config{
      Certificates:       []tls.Certificate{cert},
      RootCAs:            pool,
      InsecureSkipVerify: config.SkipVerify,
    }

    // Create protocol and client
    protocol, err := cloudevents.NewHTTP(cloudevents.WithTarget(config.Url().String()), cloudevents.WithRoundTripper(config.Transport()))
    if err != nil {
      return nil, Error(ErrUnknown, err.Error())
    }
    client, err = cloudevents.NewClient(protocol, cloudevents.WithTimeNow())
    if err != nil {
      return nil, Error(ErrUnknown, err.Error())
    }
  }

  return client, nil
}

func (config CloudEventConfig) Url() *url.URL {
  scheme := "https"
  if config.Insecure {
    scheme = "http"
  }
  return &url.URL{
    Scheme: scheme,
    Host:   fmt.Sprintf("%s:%d", config.Address, config.Port),
  }
}

func (config CloudEventConfig) Transport() *http.Transport {
  return &http.Transport{TLSClientConfig: config.Config}
}
