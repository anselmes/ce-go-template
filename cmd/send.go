// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"fmt"
	"log"
	"time"

	event "github.com/anselmes/ce-go-template/event"
	"github.com/spf13/cobra"
)

var (
  attempt int
  print bool
  retry bool
  timeout int
  verbose bool
)

var SendEventCmd = &cobra.Command{
  Use:   "send",
  Short: "Send CloudEvent",
  Long:  `
  Send a CloudEvent to a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    if err := initializeClient(); err != nil {
      log.Fatalln(event.Error(event.ErrReceiveFailed, err.Error()))
    }

    if print {
      json, err := manager.Json()
      if err != nil {
        log.Fatalf("failed to marshal CloudEvent, %v", err)
      }
      fmt.Printf("%s\n", json)
      return
    }

    if data != "" { manager.FromJson([]byte(data)) }

    if retry {
      log.Printf("Retry enabled: %d attempts with %d ms timeout", attempt, timeout)
      manager.SetRetry(attempt)
      manager.SetTimeout(time.Duration(timeout))
    } else {
      // Default to single attempt when retry is disabled
      manager.SetRetry(1)
      manager.SetTimeout(time.Duration(1000))
    }

    log.Printf("Sending CloudEvent...")

    if verbose {
      log.Println(manager.Event)
    }

    manager.Send(ctx, client)
  },
}

func init() {
  SendEventCmd.Flags().BoolVar(&retry, "retry", false, "Enable retry mechanism")
  SendEventCmd.Flags().IntVar(&attempt, "attempts", 3, "Number of retry attempts")
  SendEventCmd.Flags().IntVar(&timeout, "timeout", 1000, "Timeout between retry attempts in milliseconds")
  SendEventCmd.Flags().BoolVar(&print, "dry-run", false, "Print the CloudEvent JSON without sending it")
  SendEventCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
}
