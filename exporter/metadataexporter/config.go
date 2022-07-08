// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadataexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/metadataexporter"

import (
	"errors"

	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for file exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	// Path of the file to write to. Path is relative to current directory.
	Path string `mapstructure:"path"`

	// Destinations. A list of endpoints to which the exporter will send data.
	Destinations []string `mapstructure:"destinations"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Path == "" {
		return errors.New("path must be non-empty")
	}

	if cfg.Destinations == nil || len(cfg.Destinations) == 0 {
		return errors.New("destinations must be non-empty")
	}

	return nil
}
