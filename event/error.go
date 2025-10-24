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

func (err *CloudEventError) Error() error {
  return fmt.Errorf("CloudEventError - %d: %s", err.Code, err.Message)
}

func Error(code CloudEventErrorCodes, message ...string) error {
  err := CloudEventError{}
  err.Code = code

  if len(message) == 0 {
    err.Message = "An unknown error occurred"
  } else {
    err.Message = message[0]
  }

  return err.Error()
}
