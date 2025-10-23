// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "context"
  "fmt"
  "log"

  "github.com/spf13/cobra"
  cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
  data string
  print bool
  verbose bool
)

var SendEventCmd = &cobra.Command{
  Use:   "send",
  Short: "Send CloudEvent",
  Long:  `
  Send a CloudEvent to a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    // Initialize client with current flag values
    if err := initializeClient(); err != nil {
      log.Fatalf("Failed to initialize client: %v", err)
    }

    // MARK: - Dry Run

    if print {
      json, err := cm.Json()
      if err != nil {
        log.Fatalf("failed to marshal CloudEvent, %v", err)
      }
      fmt.Printf("%s\n", json)
      return
    }

    if data != "" { cm.FromJson([]byte(data)) }
    if verbose {
      log.Printf("Sending CloudEvent...")
      log.Println(cm.Event)
    }

    // MARK: - Submit

    ctx := cloudevents.ContextWithTarget(context.Background(), url)
    if result := client.Send(ctx, cm.Event); !cloudevents.IsACK(result) {
      log.Fatalf("failed to send, %v", result)
    }
  },
}

func init() {
  // MARK: - Flags

  SendEventCmd.Flags().StringVarP(&data, "data", "d", "", "The data payload to send in the CloudEvent")
  SendEventCmd.Flags().BoolVar(&print, "dry-run", false, "")
  SendEventCmd.Flags().BoolVar(&verbose, "verbose", false, "")
}
