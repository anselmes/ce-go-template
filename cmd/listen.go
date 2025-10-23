// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "log"
  "github.com/spf13/cobra"
)

var ListenEventCmd = &cobra.Command{
  Use:   "listen",
  Short: "Listen for CloudEvent",
  Long:  `
  Listen CloudEvent from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    if err := initializeClient(); err != nil {
      log.Fatalf("Failed to initialize client: %v", err)
    }
    log.Printf("Listening for CloudEvent...")
    log.Fatal(cm.Receive(ctx, client, cm.Display))
  },
}
