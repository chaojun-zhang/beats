// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by x-pack/dev-tools/cmd/buildspec/buildspec.go - DO NOT EDIT.

package program

import (
	"strings"

	"github.com/elastic/beats/v7/x-pack/elastic-agent/pkg/packer"
)

var Supported []Spec
var SupportedMap map[string]bool

func init() {
	// Packed Files
	// spec/filebeat.yml
	// spec/metricbeat.yml
	unpacked := packer.MustUnpack("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")
	SupportedMap = make(map[string]bool)

	for f, v := range unpacked {
		s, err := NewSpecFromBytes(v)
		if err != nil {
			panic("Cannot read spec from " + f)
		}
		Supported = append(Supported, s)
		SupportedMap[strings.ToLower(s.Name)] = true
	}
}