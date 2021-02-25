package flume_event

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
	jsprocessor "github.com/elastic/beats/v7/libbeat/processors/script/javascript/module/processor"
	"github.com/pkg/errors"
	"time"
)

func init() {
	processors.RegisterPlugin("flume_event", New)
	jsprocessor.RegisterPlugin("Flume_event", New)
}

type processor struct {
	config
	log *logp.Logger
}

// New constructs a new timestamp processor for parsing time strings into
// time.Time values.
func New(cfg *common.Config) (processors.Processor, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, errors.Wrap(err, "failed to unpack the flume event configuration")
	}
	return newFromConfig(c)
}

const logName = "processor.flume_event"

func newFromConfig(c config) (processors.Processor, error) {
	p := &processor{
		config: c,
		log:    logp.NewLogger(logName),
	}
	if c.ID != "" {
		p.log = p.log.With("instance_id", c.ID)
	}

	return p, nil
}

func (p *processor) String() string {
	return fmt.Sprintf("flume_event=[ID=%s,EventTimeField=%s]", p.ID, p.EventTimeField)
}

func (p *processor) Run(event *beat.Event) (*beat.Event, error) {

	fields, err := event.GetValue("fields")
	if err != nil {
		return event, errors.Wrap(err, "field 'fields' not exist")
	}
	event.PutValue("headers", fields)

	body, err := event.GetValue("message")
	if err != nil {
		return event, errors.Wrap(err, "field 'message' not exist")
	}

	event.PutValue("body", body)

	ts := event.Timestamp
	eventTime, err := event.GetValue(p.EventTimeField)
	if err == nil {
		ts = eventTime.(time.Time)
	}
	event.PutValue("headers.timestamp", ts.UnixNano()/1e6)
	return event, nil
}
