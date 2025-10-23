// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cli

import (
	"log"
	"os"

	"github.com/anselmes/ce-go-template/cmd"
	"github.com/spf13/cobra"
)

var (
  RootCmd = &cobra.Command {
    Use: "cecli",
    Short: "A CloudEvents CLI tool",
    Long: `
    A CloudEvents CLI tool to send and receive CloudEvents over HTTP/S.
    `,
  }
)

func Execute() {
  if err := RootCmd.Execute(); err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
}

func init() {
  RootCmd.AddCommand(cmd.VersionCmd)
  RootCmd.AddCommand(cmd.EventCmd)
}
