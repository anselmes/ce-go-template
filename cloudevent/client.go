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

type CloudEventClient struct {
  Address string
  Certificate string
  CertificateKey string
  Insecure bool
  Port    int
  SkipVerify bool
  Config *tls.Config
}

func (cc CloudEventClient) Client() (cloudevents.Client, error) {
  var client cloudevents.Client

  if cc.Insecure {
    log.Printf("Insecure mode enabled, skipping TLS verification")

    // Create protocol and client for insecure mode
    protocol, err := cloudevents.NewHTTP(cloudevents.WithTarget(cc.Url().String()))
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
    cert, err := tls.LoadX509KeyPair(cc.Certificate, cc.CertificateKey)
    if err != nil {
      return nil, Error(ErrTlsConfig, err.Error())
    }
    ca, err := os.ReadFile(cc.Certificate)
    if err != nil {
      return nil, Error(ErrTlsConfig, err.Error())
    }

    pool.AppendCertsFromPEM(ca)

    cc.Config = &tls.Config{
      Certificates:       []tls.Certificate{cert},
      RootCAs:            pool,
      InsecureSkipVerify: cc.SkipVerify,
    }

    // Create protocol and client
    protocol, err := cloudevents.NewHTTP(cloudevents.WithTarget(cc.Url().String()), cloudevents.WithRoundTripper(cc.Transport()))
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

func (cc CloudEventClient) Url() *url.URL {
  scheme := "https"
  if cc.Insecure {
    scheme = "http"
  }
  return &url.URL{
    Scheme: scheme,
    Host:   fmt.Sprintf("%s:%d", cc.Address, cc.Port),
  }
}

func (cc CloudEventClient) Transport() *http.Transport {
  return &http.Transport{TLSClientConfig: cc.Config}
}
