// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type callback func(event cloudevents.Event)

type CloudEventManager struct {
  Message Message
  Event cloudevents.Event
  retry int
  timeout int
  uri string
  cetype string
}

func (cm *CloudEventManager) RetryCount() int { return cm.retry }
func (cm *CloudEventManager) Timeout() time.Duration { return time.Duration(cm.timeout) }

func (cm *CloudEventManager) SetRetry(count int) { cm.retry = count }
func (cm *CloudEventManager) SetTimeout(time time.Duration) { cm.timeout = int(time) }

func (cm *CloudEventManager) Send(ctx context.Context, cc CloudEventClient) {
  count := cm.retry
  timeout := cm.timeout

  client, e := cc.Client(); if e != nil {
    err.Code = ErrUnknown
    err.Message = e.Error()
    log.Fatalln(err.Error())
  }

  for i := 0; i < count; i++ {
    if result := client.Send(ctx, cm.Event); cloudevents.IsUndelivered(result) {
      err.Code = ErrSendFailed
      err.Message = result.Error()
      log.Fatalln(err.Error())
      continue
    } else if cloudevents.IsACK(result) {
      log.Printf("Result: 200")
      break
    } else if cloudevents.IsNACK(result) {
      err.Code = ErrNotAccepted
      err.Message = result.Error()
      log.Fatalln(err.Error())
      continue
    } else {
      log.Printf("Result: %v", result)
    }

    time.Sleep(time.Duration(timeout) * time.Millisecond)
    log.Printf("Retrying to send CloudEvent, attempt %d/%d", i+1, count)
  }
}

// FIXME: TLS & rename to Listen
func (cm *CloudEventManager) Receive(ctx context.Context, client cloudevents.Client, callback callback) error {
  if e := client.StartReceiver(ctx, callback); e != nil {
    err.Code = ErrReceiveFailed
    err.Message = e.Error()
    return err.Error()
  }
  return nil
}

func (cm *CloudEventManager) Display(event cloudevents.Event) {
  log.Printf("Context Attributes,")
  log.Printf("  specversion: %s", event.SpecVersion())
  log.Printf("  type: %s", event.Type())
  log.Printf("  source: %s", event.Source())
  log.Printf("  id: %s", event.ID())
  log.Printf("  datacontenttype: %s", event.DataContentType())

  log.Printf("Data,")
  log.Printf("  %s", string(event.Data()))
}

func (cm *CloudEventManager) Handler(w http.ResponseWriter, req *http.Request) error {
  event, e := cloudevents.NewEventFromHTTPRequest(req)
  if e != nil {
    err.Code = ErrUnknown
    err.Message = e.Error()

    log.Fatalln(err.Error())
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

    return err.Error()
  }

  w.Write([]byte(event.String()))
  log.Println(event.String())

  return nil
}

func (cm *CloudEventManager) Json() ([]byte, error) {
  result, e := json.Marshal(cm.Event)
  if e != nil {
    err.Code = ErrInvalidFormat
    err.Message = e.Error()
    return nil, err.Error()
  }
  return result, nil
}

func (cm *CloudEventManager) FromJson(bytes []byte) {
  e := json.Unmarshal(bytes, &cm.Message)
  if e != nil {
    err.Code = ErrInvalidFormat
    err.Message = e.Error()
    log.Fatalln(err.Error())
    return
  }
  cm.Event.SetData(cloudevents.ApplicationJSON, cm.Message)
}

func NewCloudEventManager(msg Message, opts ...CloudEventOptions) *CloudEventManager {
  cm := &CloudEventManager{Message: msg }

  event := cloudevents.NewEvent()
  event.SetID(uuid.New().String())
  source := "ce/uri"
  cetype := "ce.type"

  if len(opts) > 0 {
    if opts[0].Source != "" {
      source = opts[0].Source
    }
    if opts[0].Type != "" {
      cetype = opts[0].Type
    }
  }

  event.SetSource(source)
  event.SetType(cetype)
  event.SetData(cloudevents.ApplicationJSON, msg)

  cm.uri = source
  cm.cetype = cetype
  cm.Event = event

  return cm
}
