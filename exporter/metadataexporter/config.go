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

	// Destinations. A list of endpoints to which the exporter will send data.
	Destinations []string `mapstructure:"destinations"`

	// Endpoint for the Purview
	Endpoint string `mapstructure:"endpoint"`

	// Account name for the Purview
	AccountName string `mapstructure:"account_name"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Destinations == nil || len(cfg.Destinations) == 0 {
		return errors.New("destinations must be non-empty")
	}

	if cfg.Endpoint == "" {
		return errors.New("endpoint must be non-empty")
	}

	if cfg.AccountName == "" {
		return errors.New("account_name must be non-empty")
	}

	return nil
}
