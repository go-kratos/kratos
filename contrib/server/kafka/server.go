// Copyright (c) 2012-2022 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package kafka

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"

	"golang.org/x/sync/errgroup"
)

var _ transport.Server = (*Server)(nil)

// Server is a Kafka server wrapper
type Server struct {
	consumers []Consumer
	handlers  map[string]Handler
	logger    log.Helper
}

// ServerOption is a Kafka server option.
type ServerOption func(server *Server)

// Consumers registers a set of consumers to the Server.
func Consumers(consumers []Consumer) ServerOption {
	return func(server *Server) {
		server.consumers = consumers
	}
}

// Handlers registers a set of handlers to the Server.
func Handlers(handlers []Handler) ServerOption {
	return func(server *Server) {
		for _, handler := range handlers {
			server.handlers[handler.Topic()] = handler
		}
	}
}

// NewServer creates a Kafka server by options.
func NewServer(opts ...ServerOption) (*Server, error) {
	server := &Server{handlers: make(map[string]Handler)}

	for _, o := range opts {
		o(server)
	}

	if len(server.consumers) == 0 {
		return nil, fmt.Errorf("no consumers")
	}
	if len(server.handlers) == 0 {
		return nil, fmt.Errorf("no handlers")
	}

	for _, srvConsumer := range server.consumers {
		for _, topic := range srvConsumer.Topics() {
			if srvConsumer.HasHandler(topic) {
				return nil, fmt.Errorf("duplicated handler for topic %s", topic)
			}
			handler, ok := server.handlers[topic]
			if !ok {
				return nil, fmt.Errorf("no available handler for topic %s", topic)
			}
			srvConsumer.RegisterHandler(handler)
		}
	}
	return server, nil
}

// Start starts the Kafka server
func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	for _, consumer := range s.consumers {
		consumer := consumer
		eg.Go(func() error {
			return consumer.Consume(ctx)
		})
	}

	return eg.Wait()
}

// Stop stops the Kafka server
func (s *Server) Stop(ctx context.Context) error {
	var result error
	for _, consumer := range s.consumers {
		if err := consumer.Close(); err != nil {
			s.logger.Errorf("close consumer error: %v", err)
			result = err
		}
	}

	return result
}
