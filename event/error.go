// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
	"fmt"
)

const (
  ErrUnknown CloudEventErrorCodes = iota + 1
  ErrInvalidFormat
  ErrInvalidURL
  ErrTlsConfig
  ErrSendFailed
  ErrReceiveFailed
  ErrNotAccepted
)

type CloudEventErrorCodes int

type CloudEventError struct {
  Code    CloudEventErrorCodes `json:"code"`
  Message string `json:"message"`
}

func (e *CloudEventError) Error() error {
  return fmt.Errorf("CloudEventError - %d: %s", e.Code, e.Message)
}

func Error(code CloudEventErrorCodes, msg ...string) error {
  err := CloudEventError{}
  err.Code = code

  if len(msg) == 0 {
    err.Message = "An unknown error occurred"
  } else {
    err.Message = msg[0]
  }

  return err.Error()
}
