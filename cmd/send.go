// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "log"
  "github.com/spf13/cobra"
)

var SendEventCmd = &cobra.Command{
  Use:   "send",
  Short: "Send CloudEvent",
  Long:  `
  Send a CloudEvent to a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Printf("Hello from SendEventCmd!")
  },
}
