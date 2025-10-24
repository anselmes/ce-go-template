// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"context"

	"github.com/anselmes/ce-go-template/api"
	event "github.com/anselmes/ce-go-template/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/cobra"
)

var (
  address string
  port int
  endpoint string

  cert string
  key string
  insecure bool
  verify bool

  client cloudevents.Client
  config *event.CloudEventConfig
  manager *event.CloudEventManager
  ctx context.Context

  data string
)

// MARK: - Command

var EventCmd = &cobra.Command{
  Use:   "event",
  Aliases: []string{"ev", "evt"},
  Short: "Send & Receive CloudEvent",
  Long:  `
  Send and Receive a CloudEvent to and from a specified target.
  `,
}

func init() {
  EventCmd.PersistentFlags().StringVar(&address, "address", "localhost", "The address to listen on")
  EventCmd.PersistentFlags().IntVar(&port, "port", 8080, "The port to listen on")

  EventCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Disable TLS verification")
  EventCmd.PersistentFlags().BoolVar(&verify, "verify", true, "Enable TLS verification")
  EventCmd.PersistentFlags().StringVar(&cert, "cert", "tls-bundle.pem", "Path to TLS certificate file")
  EventCmd.PersistentFlags().StringVar(&key, "key", "tls-key.pem", "Path to TLS key file")

  EventCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "CloudEvent data payload to send")

  // MARK: - Sub Command

  EventCmd.AddCommand(EventWebhookCmd)
  EventCmd.AddCommand(ListenEventCmd)
  EventCmd.AddCommand(SendEventCmd)
}

func initializeClient() error {
  manager = event.NewCloudEventManager(&api.Data{}, nil)
  config = &event.CloudEventConfig{
    Address: address,
    Port: port,
    Certificate: cert,
    CertificateKey: key,
    Insecure: insecure,
    SkipVerify: !verify,
  }

  endpoint = config.Url().String()
  ctx = cloudevents.ContextWithTarget(context.Background(), endpoint)

  var err error
  if client, err = config.Client(); err != nil {
    return event.Error(event.ErrUnknown, err.Error())
  }

  return nil
}
