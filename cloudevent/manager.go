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

func (cm *CloudEventManager) Send(ctx context.Context, client cloudevents.Client) {
  count := cm.retry
  timeout := cm.timeout

  for i := 0; i < count; i++ {
    result := client.Send(ctx, cm.Event)

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

// Receive starts an HTTP(S) server to listen for incoming CloudEvents
func (cm *CloudEventManager) Receive(ctx context.Context, client cloudevents.Client, callback callback, cc *CloudEventClient) error {
	// Create HTTP server
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cc.Address, cc.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			event, err := cloudevents.NewEventFromHTTPRequest(r)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			callback(*event)
			w.WriteHeader(http.StatusOK)
		}),
	}

	// Start server with or without TLS
	var err error
	if cc.Insecure {
		log.Printf("Starting HTTP server on %s:%d", cc.Address, cc.Port)
		err = server.ListenAndServe()
	} else {
		// Load TLS configuration
		cert, loadErr := tls.LoadX509KeyPair(cc.Certificate, cc.CertificateKey)
		if loadErr != nil {
			return Error(ErrTlsConfig, fmt.Sprintf("Failed to load TLS certificates: %v", loadErr))
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		log.Printf("Starting HTTPS server on %s:%d", cc.Address, cc.Port)
		err = server.ListenAndServeTLS("", "")
	}

	if err != nil && err != http.ErrServerClosed {
		return Error(ErrReceiveFailed, fmt.Sprintf("Server failed: %v", err))
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

func (cm *CloudEventManager) Handler() http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    log.Println("Received HTTP request for CloudEvent")

    event, err := cloudevents.NewEventFromHTTPRequest(req)
    if err != nil {
      result := Error(ErrUnknown, err.Error())
      log.Fatalln(result)
      http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
      log.Fatalln(result)
    }

    w.Write([]byte(event.String()))
    log.Println(event.String())
  })
}

func (cm *CloudEventManager) Json() ([]byte, error) {
  result, err := json.Marshal(cm.Event)
  if err != nil {
    return nil, Error(ErrInvalidFormat, err.Error())
  }
  return result, nil
}

func (cm *CloudEventManager) FromJson(bytes []byte) {
  err := json.Unmarshal(bytes, &cm.Message)
  if err != nil {
    log.Fatalln(Error(ErrInvalidFormat, err.Error()))
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
