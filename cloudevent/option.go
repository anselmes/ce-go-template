// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

type CloudEventOptions struct {
  Source string
  Type   string
  Data []byte
}
