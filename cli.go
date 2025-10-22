// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package main

import (
  "log"
  "os"

  "github.com/spf13/cobra"
  "github.com/anselmes/ce-go-template/cmd"
)

var (
  rootCmd = &cobra.Command {
    Use: "cecli",
    Short: "A CloudEvents CLI tool",
    Long: `
    A CloudEvents CLI tool to send and receive CloudEvents over HTTP/S.
    `,
  }
)

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.AddCommand(cmd.VersionCmd)
  rootCmd.AddCommand(cmd.EventCmd)
}
