// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
  "context"
  "fmt"
  "log"
  "encoding/json"
  "net/http"

  cloudevents "github.com/cloudevents/sdk-go/v2"
  "github.com/google/uuid"
)

type callback func(event cloudevents.Event)

type CloudEventManager struct {
  Message Message
  Event cloudevents.Event
  uri string
  cetype string
}

func (cm *CloudEventManager) Send(ctx context.Context, client cloudevents.Client) {
  if result := client.Send(ctx, cm.Event); cloudevents.IsUndelivered(result) {
    err.Code = ErrSendFailed
    err.Message = fmt.Sprintf("failed to send, %v", result)
    log.Fatalln(err.Error())
  } else {
    log.Printf("sent: %v", cm.Event)
    log.Printf("result: %v", result)
  }
}

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

  return &CloudEventManager{
    uri:     source,
    cetype:  cetype,
    Event:   event,
    Message: msg,
  }
}
