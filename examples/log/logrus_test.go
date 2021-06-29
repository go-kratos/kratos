// Copyright (c) 2012-2021 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
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

package log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerLog(t *testing.T) {
	tests := map[string]struct {
		level     logrus.Level
		formatter logrus.Formatter
		logLevel  log.Level
		kvs       []interface{}
		want      string
	}{
		"json format":   {level: logrus.InfoLevel, formatter: &logrus.JSONFormatter{}, logLevel: log.LevelInfo, kvs: []interface{}{"case", "json format", "msg", "1"}, want: `{"case":"json format","level":"info","msg":"1"`},
		"level unmatch": {level: logrus.InfoLevel, formatter: &logrus.JSONFormatter{}, logLevel: log.LevelDebug, kvs: []interface{}{"case", "level unmatch", "msg", "1"}, want: ""},
		"no tags":       {level: logrus.InfoLevel, formatter: &logrus.JSONFormatter{}, logLevel: log.LevelInfo, kvs: []interface{}{"msg", "1"}, want: `{"level":"info","msg":"1"`},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := new(bytes.Buffer)
			logger := NewLogrusLogger(Level(test.level), Formatter(test.formatter), Output(output))
			_ = logger.Log(test.logLevel, test.kvs...)

			assert.True(t, strings.HasPrefix(output.String(), test.want))
		})
	}
}
