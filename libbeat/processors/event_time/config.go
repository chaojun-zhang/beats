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

type config struct {
	ID           string `config:"id"` // An identifier for this processor. Useful for debugging.
	TargetField  string `config:"target_field"`
	MessageFormat    string `config:"log_format"`                      // log format e.g. (json,csv,text)
	Timezone     string `config:"timezone"`                        // Timezone (e.g. America/New_York) to use when parsing a timestamp not containing a timezone.
	TimeLayout   string `config:"time_layout" validate:"required"` // Timestamp layouts that define the expected time value format.
	TimeFormat   string `config:"time_format" validate:"required"` // time format .e.g. (date,unix,unix_ms,local)
	TimeStartPos int    `config:"time_start_pos"`                  // start position of time
	CsvDelimiter string `config:"csv_delimiter"`                   //csv delimiter
	CsvFieldPos  int    `config:"csv_field_pos"`
	JsonField    string `config:"json_field"`
	IgnoreFailure bool   `config:"json_field"`                     // Ignore errors when parsing the timestamp.
	IgnoreMissing bool   `config:"json_field"`                     // Ignore errors when the source field is missing.
}

func defaultConfig() config {
	return config{
		TargetField:   "@event_time",
		TimeFormat:    "DATE",
		MessageFormat: "TEXT",
		CsvDelimiter:   ",",
	}
}
