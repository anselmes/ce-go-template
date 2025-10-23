// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/cobra"

	. "github.com/anselmes/ce-go-template/cloudevent"
)

var (
  address string
  port int
  url string

  // ca string
  cert string
  key string
  insecure bool
  verify bool

  client cloudevents.Client
  cc *CloudEventClient
  cm *CloudEventManager
  ctx context.Context

  err = CloudEventError{}
)

// MARK: - Command

var EventCmd = &cobra.Command{
  Use:   "event",
  Short: "Send & Receive CloudEvent",
  Long:  `
  Send and Receive a CloudEvent to and from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Printf("Hello from CE (%s)!", url)
    // TODO: to sink using cm.Handle()
  },
}

func init() {
  EventCmd.PersistentFlags().StringVar(&address, "address", "localhost", "The address to listen on")
  EventCmd.PersistentFlags().IntVar(&port, "port", 8080, "The port to listen on")

  EventCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Disable TLS verification")
  EventCmd.PersistentFlags().BoolVar(&verify, "verify", true, "Enable TLS verification")
  // EventCmd.PersistentFlags().StringVar(&ca, "ca", "ca.crt", "Path to CA certificate file")
  EventCmd.PersistentFlags().StringVar(&cert, "cert", "tls-bundle.pem", "Path to TLS certificate file")
  EventCmd.PersistentFlags().StringVar(&key, "key", "tls-key.pem", "Path to TLS key file")

  // MARK: - Sub Command

  EventCmd.AddCommand(SendEventCmd)
  EventCmd.AddCommand(ListenEventCmd)
}

func initializeClient() error {
  cm = NewCloudEventManager(Message{})
  cc = &CloudEventClient{
    Address: address,
    Port: port,
    Certificate: cert,
    CertificateKey: key,
    Insecure: insecure,
    SkipVerify: !verify,
  }

  url = cc.Url()
  ctx = cloudevents.ContextWithTarget(context.Background(), url)

  var e error
  if client, e = cc.Client(); e != nil {
    err.Code = ErrUnknown
    err.Message = e.Error()
    return err.Error()
  }

  return nil
}
