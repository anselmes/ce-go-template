// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
)

var (
  Name    = "cecli"
  Version = "dev"

  VersionCmd = &cobra.Command {
    Use: "version",
    Short: "Show the CLI version",
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("%s version %s\n", Name, Version)
    },
  }
)
