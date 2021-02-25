// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package eventtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
)

var expected = time.Date(2021, 1, 25, 7, 44, 52, 0, time.UTC)

func TestParsePatternsForCsv(t *testing.T) {
	logp.TestingSetup()

	c := defaultConfig()

	p, err := newFromConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	evt := &beat.Event{Fields: common.MapStr{}}

	t.Run("date", func(t *testing.T) {
		p.LogFormat = "csv"
		p.TimeFormat = "date"
		p.TimeStartPos = 1
		p.TimeLayout = "02/Jan/2006:15:04:05 -0700"
		p.CsvDelimiter = "|#|"
		p.CsvFieldPos = 2
		evt.Timestamp = time.Time{}
		evt.PutValue("message", "223.104.191.18|#|-|#|[25/Jan/2021:15:44:52 +0800]|#|1614235492.382|bbbbbbbbb")
		evt, err = p.Run(evt)
		if err != nil {
			t.Fatal(err)
		}
		event_time, err := evt.GetValue("@event_time")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, event_time)
	})

	t.Run("UNIX", func(t *testing.T) {
		p.LogFormat = "csv"
		p.TimeFormat = "UNIX"
		p.TimeStartPos = 1
		p.TimeLayout = "1611560692"
		p.CsvDelimiter = "|#|"
		p.CsvFieldPos = 2
		evt.Timestamp = time.Time{}
		evt.PutValue("message", "223.104.191.18|#|-|#|[1611560692]|#|1614235492.382|bbbbbbbbb")
		evt, err = p.Run(evt)
		if err != nil {
			t.Fatal(err)
		}
		event_time, err := evt.GetValue("@event_time")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, event_time)
	})

}
