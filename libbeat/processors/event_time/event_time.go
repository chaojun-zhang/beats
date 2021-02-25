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

package event_time

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"4d63.com/tz"
	"github.com/pkg/errors"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
	jsprocessor "github.com/elastic/beats/v7/libbeat/processors/script/javascript/module/processor"
)

const logName = "processor.eventtime"

func init() {
	processors.RegisterPlugin("event_time", New)
	jsprocessor.RegisterPlugin("event_time", New)
}

type processor struct {
	config
	log     *logp.Logger
	isDebug bool
	tz      *time.Location
}

// New constructs a new timestamp processor for parsing time strings into
// time.Time values.
func New(cfg *common.Config) (processors.Processor, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, errors.Wrap(err, "failed to unpack the timestamp configuration")
	}

	return newFromConfig(c)
}

func newFromConfig(c config) (*processor, error) {
	loc, err := loadLocation(c.Timezone)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load timezone")
	}

	p := &processor{
		config:  c,
		log:     logp.NewLogger(logName),
		isDebug: logp.IsDebug(logName),
		tz:      loc,
	}
	if c.ID != "" {
		p.log = p.log.With("instance_id", c.ID)
	}

	return p, nil
}

var timezoneFormats = []string{"-07", "-0700", "-07:00"}

func loadLocation(timezone string) (*time.Location, error) {
	for _, format := range timezoneFormats {
		t, err := time.Parse(format, timezone)
		if err == nil {
			name, offset := t.Zone()
			return time.FixedZone(name, offset), nil
		}
	}

	// Rest of location formats
	return tz.LoadLocation(timezone)
}

func (p *processor) String() string {
	return fmt.Sprintf("timestamp=[timezone=%v, layout=%v]", p.tz, p.TimeLayout)
}

func (p *processor) Run(event *beat.Event) (*beat.Event, error) {
	// Get the source field value.
	val, err := event.GetValue("message")
	if err != nil {
		return event, errors.Wrap(err, "failed to get time from message field")
	}

	// Try to convert the value to a time.Time.
	ts, err := p.tryToTime(val, event)
	if err != nil {
		return event, err
	}

	// Put the timestamp as UTC into the target field.
	_, err = event.PutValue(p.TargetField, ts.UTC())
	if err != nil {
		return event, err
	}

	return event, nil
}

func (p *processor) tryToTime(value interface{}, event *beat.Event) (time.Time, error) {
	if p.TimeFormat == "LOCAL" {
		return event.Timestamp, nil
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case common.Time:
		return time.Time(v), nil
	default:
		return p.parseValue(v)
	}
}

func (p *processor) parseValue(v interface{}) (time.Time, error) {
	timeStr, err := p.parseTime(v)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "fail to parse time value")
	}

	switch p.TimeFormat {
	case "UNIX":
		if sec, ok := common.TryToInt(timeStr); ok {
			return time.Unix(int64(sec), 0), nil
		} else if sec, ok := common.TryToFloat64(v); ok {
			return time.Unix(0, int64(sec*float64(time.Second))), nil
		}
		return time.Time{}, errors.New("could not parse time field as int or float")
	case "UNIX_MS":
		if ms, ok := common.TryToInt(timeStr); ok {
			return time.Unix(0, int64(ms)*int64(time.Millisecond)), nil
		} else if ms, ok := common.TryToFloat64(timeStr); ok {
			return time.Unix(0, int64(ms*float64(time.Millisecond))), nil
		}
		return time.Time{}, errors.New("could not parse time field as int or float")
	default:
		ts, err := time.ParseInLocation(p.TimeLayout, timeStr, p.tz)
		if err == nil {
			// Use current year if no year is zero.
			if ts.Year() == 0 {
				currentYear := time.Now().In(ts.Location()).Year()
				ts = ts.AddDate(currentYear, 0, 0)
			}
		}
		return ts, err
	}
}

func (p *processor) lengthOfLayout() int {
	return len(p.TimeLayout)
}

func (p *processor) parseTime(v interface{}) (string, error) {
	timeStr, ok := v.(string)
	if !ok {
		return "", errors.Errorf("unexpected type %T for time field", v)
	}
	switch p.LogFormat {
	case "json":
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(timeStr), &dat); err != nil {
			return timeStr, errors.Wrapf(err, "could not parse time field as int or float, message: %T", v)
		}
		timeStr = dat[p.JsonField].(string)
	case "csv":
		csvFields := strings.Split(timeStr, p.CsvDelimiter)
		timeStr = csvFields[p.CsvFieldPos]
	}

	if p.TimeStartPos+p.lengthOfLayout() > len(timeStr) {
		timeStr = timeStr[p.TimeStartPos:len(timeStr)]
	} else {
		timeStr = timeStr[p.TimeStartPos : p.TimeStartPos+p.lengthOfLayout()]
	}

	return timeStr, nil
}
