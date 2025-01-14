// Copyright OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azureblobexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azureblobexporter"

import (
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for Azure Blob
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	Endpoint                string `mapstructure:"endpoint"`
	ContainerName           string `mapstructure:"container_name"`
	StorageAccountName      string `mapstructure:"storage_account_name"`
}
