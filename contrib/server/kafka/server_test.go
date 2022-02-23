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
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/go-kratos/kratos/v2/log"
)

func TestNewServer(t *testing.T) {
	topic := "test"
	mockConsumer := &MockConsumer{}
	mockConsumer.On("Topic").Return(topic)
	mockConsumer.On("RegisterHandler", mock.Anything).Return()
	mockHandler := &MockHandler{}
	mockHandler.On("Topic").Return(topic)
	mockErrHandler := &MockHandler{}
	mockErrHandler.On("Topic").Return("")

	type args struct {
		consumers []Consumer
		handlers  []Handler
		opts      []ServerOption
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "happy path",
			args: struct {
				consumers []Consumer
				handlers  []Handler
				opts      []ServerOption
			}{
				consumers: []Consumer{mockConsumer},
				handlers:  []Handler{mockHandler},
				opts:      []ServerOption{Logger(log.GetLogger())},
			},
			wantErr: false,
		},
		{
			name: "no consumers",
			args: struct {
				consumers []Consumer
				handlers  []Handler
				opts      []ServerOption
			}{
				consumers: []Consumer{},
				handlers:  []Handler{mockHandler},
				opts:      []ServerOption{Logger(log.GetLogger())},
			},
			wantErr: true,
		},
		{
			name: "no handlers",
			args: struct {
				consumers []Consumer
				handlers  []Handler
				opts      []ServerOption
			}{
				consumers: []Consumer{mockConsumer},
				handlers:  []Handler{},
				opts:      []ServerOption{Logger(log.GetLogger())},
			},
			wantErr: true,
		},
		{
			name: "no matched handler",
			args: struct {
				consumers []Consumer
				handlers  []Handler
				opts      []ServerOption
			}{
				consumers: []Consumer{mockConsumer},
				handlers:  []Handler{mockErrHandler},
				opts:      []ServerOption{Logger(log.GetLogger())},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServer(tt.args.consumers, tt.args.handlers, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	topic := "test"
	mockConsumer := &MockConsumer{}
	mockConsumer.On("Topic").Return(topic)
	mockConsumer.On("Consume", mock.Anything).Return(nil)

	type fields struct {
		consumers []Consumer
		handlers  map[string]Handler
		logger    *log.Helper
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				consumers: []Consumer{mockConsumer},
				handlers:  map[string]Handler{topic: &MockHandler{}},
				logger:    log.NewHelper(log.GetLogger()),
			},
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				consumers: tt.fields.consumers,
				handlers:  tt.fields.handlers,
				logger:    tt.fields.logger,
			}
			if err := s.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Stop(t *testing.T) {
	mockConsumer := &MockConsumer{}
	mockConsumer.On("Close", mock.Anything).Return(nil)
	mockErrConsumer := &MockConsumer{}
	mockErrConsumer.On("Close", mock.Anything).Return(fmt.Errorf("error"))

	type fields struct {
		consumers []Consumer
		handlers  map[string]Handler
		logger    *log.Helper
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				consumers: []Consumer{mockConsumer},
				handlers:  map[string]Handler{"handler": &MockHandler{}},
				logger:    log.NewHelper(log.GetLogger()),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name: "stop error",
			fields: fields{
				consumers: []Consumer{mockErrConsumer},
				handlers:  map[string]Handler{"handler": &MockHandler{}},
				logger:    log.NewHelper(log.GetLogger()),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				consumers: tt.fields.consumers,
				handlers:  tt.fields.handlers,
				logger:    tt.fields.logger,
			}
			if err := s.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
