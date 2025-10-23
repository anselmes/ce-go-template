// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	. "github.com/anselmes/ce-go-template/cloudevent"
)

var (
  data string
  print bool
  attempt int
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
    if e := initializeClient(); e != nil {
      err.Code = ErrReceiveFailed
      err.Message = e.Error()
      log.Fatalln(err.Error())
    }

    if print {
      json, err := cm.Json()
      if err != nil {
        log.Fatalf("failed to marshal CloudEvent, %v", err)
      }
      fmt.Printf("%s\n", json)
      return
    }

    if data != "" { cm.FromJson([]byte(data)) }
    cm.SetRetry(attempt)
    cm.SetTimeout(time.Duration(1000))

    if verbose {
      log.Printf("Sending CloudEvent...")
      log.Printf("Retry enabled: %d attempts with %d ms timeout", attempt, timeout)
      log.Println(cm.Event)
    }

    cm.Send(ctx, client)
  },
}

func init() {
  SendEventCmd.Flags().StringVarP(&data, "data", "d", "", "The data payload to send in the CloudEvent")
  SendEventCmd.Flags().IntVar(&attempt, "attempts", 3, "Number of retry attempts")
  SendEventCmd.Flags().IntVar(&timeout, "timeout", 1000, "Timeout between retry attempts in milliseconds")
  SendEventCmd.Flags().BoolVar(&print, "dry-run", false, "Print the CloudEvent JSON without sending it")
  SendEventCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
}
