// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "context"
  "fmt"
  "log"

  "github.com/spf13/cobra"
  cloudevents "github.com/cloudevents/sdk-go/v2"

  cecli "github.com/anselmes/ce-go-template/cloudevent"
)

var (
  url string
  port int
  address string

  ca string
  cert string
  key string
  insecure bool

  ctx context.Context
  cm *cecli.CloudEventManager
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
  },
}

func init() {
  EventCmd.PersistentFlags().StringVar(&address, "address", "localhost", "The address to listen on")
  EventCmd.PersistentFlags().IntVar(&port, "port", 8080, "The port to listen on")

  EventCmd.PersistentFlags().StringVar(&ca, "ca", "ca.crt", "Path to CA certificate file")
  EventCmd.PersistentFlags().StringVar(&cert, "cert", "tls.crt", "Path to TLS certificate file")
  EventCmd.PersistentFlags().StringVar(&key, "key", "tls.key", "Path to TLS key file")
  EventCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Disable TLS verification")

  if insecure {
    url = fmt.Sprintf("http://%s:%d", address, port)
  } else {
    url = fmt.Sprintf("https://%s:%d", address, port)
  }

  cm = cecli.NewCloudEventManager(cecli.Message{})
  ctx = cloudevents.ContextWithTarget(context.Background(), url)

  // MARK: - Sub Command

  EventCmd.AddCommand(SendEventCmd)
  EventCmd.AddCommand(ListenEventCmd)
}
