// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
  "fmt"
  "log"
  "os"
	"crypto/tls"
	"crypto/x509"
  "net/http"

  cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CloudEventClient struct {
  Address string
  Certificate string
  CertificateKey string
  Insecure bool
  Port    int
  Scheme string
  SkipVerify bool
  Config *tls.Config
}

func (cc CloudEventClient) Client() (cloudevents.Client, error) {
  var client cloudevents.Client

  if cc.Insecure {
    cc.Scheme = "http"
    log.Printf("Insecure mode enabled, skipping TLS verification")

    // Create protocol and client for insecure mode
    protocol, e := cloudevents.NewHTTP(cloudevents.WithTarget(cc.Url()))
    if e != nil {
      err.Code = ErrUnknown
      err.Message = e.Error()
      return nil, err.Error()
    }
    client, e = cloudevents.NewClient(protocol, cloudevents.WithTimeNow())
    if e != nil {
      err.Code = ErrUnknown
      err.Message = e.Error()
      return nil, err.Error()
    }
  } else {
    cc.Scheme = "https"
    pool := x509.NewCertPool()

    // Configure a new http.Transport with TLS
    cert, e := tls.LoadX509KeyPair(cc.Certificate, cc.CertificateKey)
    if e != nil {
      err.Code = ErrTlsConfig
      err.Message = e.Error()
      return nil, err.Error()
    }
    ca, e := os.ReadFile(cc.Certificate)
    if e != nil {
      err.Code = ErrTlsConfig
      err.Message = e.Error()
      return nil, err.Error()
    }

    pool.AppendCertsFromPEM(ca)

    cc.Config = &tls.Config{
      Certificates:       []tls.Certificate{cert},
      RootCAs:            pool,
      InsecureSkipVerify: cc.SkipVerify,
    }

    // Create protocol and client
    protocol, e := cloudevents.NewHTTP(cloudevents.WithTarget(cc.Url()), cloudevents.WithRoundTripper(cc.Transport()))
    if e != nil {
      err.Code = ErrUnknown
      err.Message = e.Error()
      return nil, err.Error()
    }
    client, e = cloudevents.NewClient(protocol, cloudevents.WithTimeNow())
    if e != nil {
      err.Code = ErrUnknown
      err.Message = e.Error()
      return nil, err.Error()
    }
  }

  return client, nil
}

func (cc CloudEventClient) Url() string {
  return fmt.Sprintf("%s://%s:%d", cc.Scheme, cc.Address, cc.Port)
}

func (cc CloudEventClient) Transport() *http.Transport {
  return &http.Transport{TLSClientConfig: cc.Config}
}
