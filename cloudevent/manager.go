// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
  "context"
  "log"
  "encoding/json"
  "net/http"

  cloudevents "github.com/cloudevents/sdk-go/v2"
  "github.com/google/uuid"
)

type CloudEventManager struct {
  Message Message
  Event cloudevents.Event
  uri string
  cetype string
}

func (cm *CloudEventManager) Send(ctx context.Context, client cloudevents.Client) {
  if result := client.Send(ctx, cm.Event); cloudevents.IsUndelivered(result) {
    log.Fatalf("failed to send, %v", result)
  } else {
    log.Printf("sent: %v", cm.Event)
    log.Printf("result: %v", result)
  }
}

func (cm *CloudEventManager) Receive(ctx context.Context, client cloudevents.Client) error {
  err := client.StartReceiver(ctx, cm.onReceive)
  if err != nil { log.Fatal("failed to start receiver: %v", err) }
  return nil
}

func (cm *CloudEventManager) Handler(w http.ResponseWriter, req *http.Request) error {
  event, err := cloudevents.NewEventFromHTTPRequest(req)

  if err != nil {
    log.Printf("failed to parse CloudEvent from request: %v", err)
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return err
  }

  w.Write([]byte(event.String()))
  log.Println(event.String())

  return nil
}

func (cm *CloudEventManager) Json() ([]byte, error) {
  return json.Marshal(cm.Event)
}

func (cm *CloudEventManager) FromJson(bytes []byte) {
  err := json.Unmarshal(bytes, &cm.Message)
  if err != nil {
    log.Fatalf("failed to unmarshal CloudEvent, %v", err)
  }
  cm.Event.SetData(cloudevents.ApplicationJSON, cm.Message)
}

// TODO: Process the received event
func (cm *CloudEventManager) onReceive(event cloudevents.Event) {
  log.Printf("Context Attributes,")
  log.Printf("  specversion: %s", event.SpecVersion())
  log.Printf("  type: %s", event.Type())
  log.Printf("  source: %s", event.Source())
  log.Printf("  id: %s", event.ID())
  log.Printf("  datacontenttype: %s", event.DataContentType())

  log.Printf("Data,")
  log.Printf("  %s", string(event.Data()))
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
