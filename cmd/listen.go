// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "log"

  "github.com/spf13/cobra"
  cloudevents "github.com/cloudevents/sdk-go/v2"
)

var ListenEventCmd = &cobra.Command{
  Use:   "listen",
  Short: "Listen for CloudEvent",
  Long:  `
  Listen CloudEvent from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Printf("Listening for CloudEvent...")

    // FIXME: The default client is HTTP.
    client, err := cloudevents.NewClientHTTP()
    if err != nil {
      log.Fatalf("failed to create client, %v", err)
    }

    log.Fatal(cm.Receive(ctx, client))
  },
}
