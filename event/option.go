// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import "github.com/anselmes/ce-go-template/api"

type CloudEventOptions struct {
  Source string
  Type   string
  Data api.Data
}
