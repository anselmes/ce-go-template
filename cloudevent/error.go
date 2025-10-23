// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import "fmt"

const (
  ErrUnknown CloudEventErrorCodes = iota + 1
  ErrInvalidFormat
  ErrInvalidURL
  ErrTlsConfig
  ErrSendFailed
  ErrReceiveFailed
  ErrNotAccepted
)

var err = CloudEventError{}
type CloudEventErrorCodes int

type CloudEventError struct {
  Code    CloudEventErrorCodes `json:"code"`
  Message string `json:"message"`
}

func (e *CloudEventError) Error() error {
  return fmt.Errorf("CloudEventError - %d: %s", e.Code, e.Message)
}
