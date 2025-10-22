// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "fmt"
  "log"

  "github.com/spf13/cobra"
)

var (
  port int
  address string

  ca string
  cert string
  key string
  insecure bool
)

// MARK: - Command

var EventCmd = &cobra.Command{
  Use:   "event",
  Short: "Send & Receive CloudEvent",
  Long:  `
  Send and Receive a CloudEvent to and from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    url := ""

    if insecure {
      url = fmt.Sprintf("http://%s:%d", address, port)
    } else {
      url = fmt.Sprintf("https://%s:%d", address, port)
    }

    log.Printf("Hello from CE (%s)!", url)
  },
}

func init() {
  EventCmd.PersistentFlags().StringVarP(&address, "address", "", "localhost", "The address to listen on")
  EventCmd.PersistentFlags().IntVarP(&port, "port", "", 8080, "The port to listen on")

  // MARK: - Certificate Flags

  EventCmd.PersistentFlags().StringVarP(&ca, "ca", "", "ca.crt", "Path to CA certificate file")
  EventCmd.PersistentFlags().StringVarP(&cert, "cert", "", "tls.crt", "Path to TLS certificate file")
  EventCmd.PersistentFlags().StringVarP(&key, "key", "", "tls.key", "Path to TLS key file")
  EventCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "", false, "Disable TLS verification")

  // MARK: - Sub Command

  EventCmd.AddCommand(SendEventCmd)
  EventCmd.AddCommand(ReceiveEventCmd)
}
