// SPDX-License-Identifier: GPL-3.0
// Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

package cmd

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/cobra"

	ev "github.com/anselmes/ce-go-template/cloudevent"
)

var (
  address string
  port int
  endpoint string

  cert string
  key string
  insecure bool
  verify bool

  client cloudevents.Client
  cc *ev.CloudEventClient
  cm *ev.CloudEventManager
  ctx context.Context

  data string
)

// MARK: - Command

var EventCmd = &cobra.Command{
  Use:   "event",
  Aliases: []string{"ev", "evt"},
  Short: "Send & Receive CloudEvent",
  Long:  `
  Send and Receive a CloudEvent to and from a specified target.
  `,
  Run: func(cmd *cobra.Command, args []string) {
    log.Printf("Hello from CE (%s)!", endpoint)

    // host, node, opts := configSink()
    // p, e := ceamqp.NewProtocol(host, node, []amqp.ConnOption{}, []amqp.SessionOption{}, opts...)
    // if e != nil {
    //   err.Code = ErrUnknown
    //   err.Message = e.Error()
    //   log.Fatalln(err.Error())
    // }

    // Close the connection when finished
    // defer p.Close(context.Background())
  },
}

func init() {
  EventCmd.PersistentFlags().StringVar(&address, "address", "localhost", "The address to listen on")
  EventCmd.PersistentFlags().IntVar(&port, "port", 8080, "The port to listen on")

  EventCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Disable TLS verification")
  EventCmd.PersistentFlags().BoolVar(&verify, "verify", true, "Enable TLS verification")
  EventCmd.PersistentFlags().StringVar(&cert, "cert", "tls-bundle.pem", "Path to TLS certificate file")
  EventCmd.PersistentFlags().StringVar(&key, "key", "tls-key.pem", "Path to TLS key file")

  EventCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "CloudEvent data payload to send")

  // MARK: - Sub Command

  EventCmd.AddCommand(SendEventCmd)
  EventCmd.AddCommand(ListenEventCmd)
  EventCmd.AddCommand(EventWebhookCmd)
}

func initializeClient() error {
  cm = ev.NewCloudEventManager(ev.Message{})
  cc = &ev.CloudEventClient{
    Address: address,
    Port: port,
    Certificate: cert,
    CertificateKey: key,
    Insecure: insecure,
    SkipVerify: !verify,
  }

  endpoint = cc.Url()
  ctx = cloudevents.ContextWithTarget(context.Background(), endpoint)

  var err error
  if client, err = cc.Client(); err != nil {
    return ev.Error(ev.ErrUnknown, err.Error())
  }

  return nil
}

// Parse AMQP_URL env variable. Return server URL, AMQP node (from path) and SASLPlain
// option if user/pass are present.
// func configSink() (server, node string, opts []ceamqp.Option) {
// 	env := os.Getenv("AMQP_URL")
// 	if env == "" {
// 		env = "/test"
// 	}
// 	u, err := url.Parse(env)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if u.User != nil {
// 		user := u.User.Username()
// 		pass, _ := u.User.Password()
// 		opts = append(opts, ceamqp.WithConnOpt(amqp.ConnSASLPlain(user, pass)))
// 	}
// 	return env, strings.TrimPrefix(u.Path, "/"), opts
// }
