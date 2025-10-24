// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"fmt"
	"log"
	"net/http"

	event "github.com/anselmes/ce-go-template/cloudevent"
	"github.com/spf13/cobra"
)

var EventWebhookCmd = &cobra.Command {
  Use: "webhook",
  Aliases: []string{"wh"},
  Short: "Handle CloudEvent via Webhook",
  Long: `
  Handle CloudEvent via Webhook.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Println("Starting webhook server to handle CloudEvents...")

    if err := initializeClient(); err != nil {
      log.Fatalln(event.Error(event.ErrReceiveFailed, err.Error()))
    }

    // Set up HTTP handler for CloudEvents
    http.Handle("/", manager.Handler())

    // Start HTTP server
    log.Println("Listening on:", config.Url())
    err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Address, config.Port), nil)
    if err != nil {
      log.Fatalln(event.Error(event.ErrUnknown, err.Error()))
    }
  },
}
