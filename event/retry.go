// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

type Retry struct {
  Enable bool
  Attempts int
  Timeout int
}

func DefaultRetry() Retry {
  return Retry{
    Enable:   false,
    Attempts: 3,
    Timeout:  1000,
  }
}
