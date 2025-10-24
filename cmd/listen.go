// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"log"

	"github.com/spf13/cobra"

	event "github.com/anselmes/ce-go-template/cloudevent"
)

var ListenEventCmd = &cobra.Command{
  Use:   "listen",
  Aliases: []string{"lis"},
  Short: "Listen for CloudEvent",
  Long:  `
  Listen CloudEvent from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    if err := initializeClient(); err != nil {
      log.Fatalln(event.Error(event.ErrReceiveFailed, err.Error()))
    }
    log.Fatal(manager.Listen(ctx, config, manager.Display))
  },
}
