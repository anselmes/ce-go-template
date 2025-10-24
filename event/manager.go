// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cloudevent

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	api "github.com/anselmes/ce-go-template/api/v1"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type callback func(event cloudevents.Event)

type CloudEventManager struct {
  Data *api.Data
  Event cloudevents.Event
  retry int
  timeout int
  uri string
  cetype string
  callback callback
}

func (manager *CloudEventManager) RetryCount() int { return manager.retry }
func (manager *CloudEventManager) Timeout() time.Duration { return time.Duration(manager.timeout) }

func (manager *CloudEventManager) SetRetry(count int) { manager.retry = count }
func (manager *CloudEventManager) SetTimeout(time time.Duration) { manager.timeout = int(time) }
func (manager *CloudEventManager) SetCallback(cb callback) { manager.callback = cb }

func (manager *CloudEventManager) Send(ctx context.Context, client cloudevents.Client) {
  count := manager.retry
  timeout := manager.timeout

  for i := 0; i < count; i++ {
    result := client.Send(ctx, manager.Event)

    if cloudevents.IsACK(result) {
      log.Printf("Result: 200")
      break // Success - exit retry loop
    } else if cloudevents.IsNACK(result) {
      log.Printf("CloudEvent was rejected: %v", result)
      if i == count-1 {
        log.Fatalln(Error(ErrNotAccepted, result.Error()))
      }
    } else if cloudevents.IsUndelivered(result) {
      log.Printf("CloudEvent delivery failed: %v", result)
      if i == count-1 {
        log.Fatalln(Error(ErrSendFailed, result.Error()))
      }
    } else {
      log.Printf("Result: %v", result)
      if i == count-1 {
        log.Printf("Exhausted all retry attempts")
      }
    }

    // Only sleep and retry if this isn't the last attempt
    if i < count-1 {
      time.Sleep(time.Duration(timeout) * time.Millisecond)
      log.Printf("Retrying to send CloudEvent, attempt %d/%d", i+1, count)
    }
  }
}

func (manager *CloudEventManager) Listen(ctx context.Context, config *CloudEventConfig, callback callback) error {
	manager.SetCallback(callback)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Address, config.Port),
		Handler: manager.Handler(),
	}


  log.Printf("Listening for CloudEvent on %s...", config.Url())

	var err error
	if config.Insecure {
		err = server.ListenAndServe()
	} else {
		// Load TLS configuration
		cert, loadErr := tls.LoadX509KeyPair(config.Certificate, config.CertificateKey)
		if loadErr != nil {
			return Error(ErrTlsConfig, fmt.Sprintf("Failed to load TLS certificates: %v", loadErr))
		}
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		err = server.ListenAndServeTLS("", "")
	}

	if err != nil && err != http.ErrServerClosed {
		return Error(ErrReceiveFailed, fmt.Sprintf("Server failed: %v", err))
	}

	return nil
}

func (manager *CloudEventManager) Receive(ctx context.Context, client cloudevents.Client, callback callback) error {
	if err := client.StartReceiver(ctx, callback); err != nil {
		return Error(ErrReceiveFailed, err.Error())
	}
	return nil
}

func (manager *CloudEventManager) Display(event cloudevents.Event) {
  log.Printf("Context Attributes,")
  log.Printf("  specversion: %s", event.SpecVersion())
  log.Printf("  type: %s", event.Type())
  log.Printf("  source: %s", event.Source())
  log.Printf("  id: %s", event.ID())
  log.Printf("  datacontenttype: %s", event.DataContentType())

  log.Printf("Data,")
  log.Printf("  %s", string(event.Data()))
}

func (manager *CloudEventManager) Handler() http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    log.Println("Received HTTP request for CloudEvent")

    event, err := cloudevents.NewEventFromHTTPRequest(req)
    if err != nil {
      http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
      return
    }

    // Use the callback if set, otherwise use Display as default
    if manager.callback != nil {
      manager.callback(*event)
    } else {
      manager.Display(*event)
    }

    w.WriteHeader(http.StatusOK)
  })
}

func (manager *CloudEventManager) Json() ([]byte, error) {
  result, err := json.Marshal(manager.Event)
  if err != nil {
    return nil, Error(ErrInvalidFormat, err.Error())
  }
  return result, nil
}

func (manager *CloudEventManager) FromJson(bytes []byte) {
  err := json.Unmarshal(bytes, manager.Data)
  if err != nil {
    log.Fatalln(Error(ErrInvalidFormat, err.Error()))
    return
  }
  manager.Event.SetData(cloudevents.ApplicationJSON, manager.Data)
}

func NewCloudEventManager(data *api.Data, opts *CloudEventOptions) *CloudEventManager {
  manager := &CloudEventManager{Data: data}

  event := cloudevents.NewEvent()
  event.SetID(uuid.New().String())
  source := "ce/uri"
  cetype := "ce.type"

  if opts != nil {
    if opts.Source != "" { source = opts.Source }
    if opts.Type != "" { cetype = opts.Type }
  }

  event.SetSource(source)
  event.SetType(cetype)
  event.SetData(cloudevents.ApplicationJSON, data.Message)

  manager.uri = source
  manager.cetype = cetype
  manager.Event = event

  return manager
}
