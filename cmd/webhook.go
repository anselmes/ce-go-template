// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var sink string
// var wr http.ResponseWriter
// var req *http.Request

var EventWebhookCmd = &cobra.Command {
  Use: "webhook",
  Aliases: []string{"wh"},
  Short: "Handle CloudEvent via Webhook",
  Long: `
  Handle CloudEvent via Webhook.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Println("Webhook handler not yet implemented.")
    log.Println("Sink:", sink)
    // TODO:
    // http listener
    // on receive, send
  },
}

func init() {
  EventWebhookCmd.Flags().StringVarP(&sink, "sink", "K", "", "The target sink URL to send the CloudEvent to")
}

// func send(event cloudevents.Event) {
//   log.Println("Sending CloudEvent via Webhook...")
//   cm.Handler(writer, request)
// }
