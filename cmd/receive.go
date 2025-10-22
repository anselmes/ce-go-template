// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "log"
  "github.com/spf13/cobra"
)

var ReceiveEventCmd = &cobra.Command{
  Use:   "receive",
  Short: "Receive CloudEvent",
  Long:  `
  Receive a CloudEvent from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Printf("Hello from ReceiveEventCmd!")
  },
}
