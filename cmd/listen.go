// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "log"

  "github.com/spf13/cobra"
  . "github.com/anselmes/ce-go-template/cloudevent"
)

var ListenEventCmd = &cobra.Command{
  Use:   "listen",
  Short: "Listen for CloudEvent",
  Long:  `
  Listen CloudEvent from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    if e := initializeClient(); e != nil {
      err.Code = ErrReceiveFailed
      err.Message = e.Error()
      log.Fatalln(err.Error())
    }

    log.Printf("Listening for CloudEvent...")
    log.Fatal(cm.Receive(ctx, client, cm.Display))
  },
}
