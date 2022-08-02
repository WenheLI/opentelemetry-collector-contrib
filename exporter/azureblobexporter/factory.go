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
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "azureblob"
)

// NewFactory returns a factory for Azure Monitor exporter.
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter),
		component.WithMetricsExporter(createMetricsExporter),
		component.WithLogsExporter(createLogsExporter))
}

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
	}
}

func createTracesExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.TracesExporter, error) {
	fe := exporters.GetOrAdd(cfg, func() component.Component {
		return &blobExporter{
			storageAccountName: cfg.(*Config).StorageAccountName,
			containerName:      cfg.(*Config).ContainerName,
			endpoint:           cfg.(*Config).Endpoint,
		}
	})
	return exporterhelper.NewTracesExporter(
		cfg,
		set,
		fe.Unwrap().(*blobExporter).ConsumeTraces,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithShutdown(fe.Shutdown),
	)
}

func createMetricsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.MetricsExporter, error) {
	fe := exporters.GetOrAdd(cfg, func() component.Component {
		return &blobExporter{
			storageAccountName: cfg.(*Config).StorageAccountName,
			containerName:      cfg.(*Config).ContainerName,
			endpoint:           cfg.(*Config).Endpoint,
		}
	})
	return exporterhelper.NewMetricsExporter(
		cfg,
		set,
		fe.Unwrap().(*blobExporter).ConsumeMetrics,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithShutdown(fe.Shutdown),
	)
}

func createLogsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.LogsExporter, error) {
	fe := exporters.GetOrAdd(cfg, func() component.Component {
		return &blobExporter{
			storageAccountName: cfg.(*Config).StorageAccountName,
			containerName:      cfg.(*Config).ContainerName,
			endpoint:           cfg.(*Config).Endpoint,
		}
	})
	return exporterhelper.NewLogsExporter(
		cfg,
		set,
		fe.Unwrap().(*blobExporter).ConsumeLogs,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithShutdown(fe.Shutdown),
	)
}

// This is the map of already created File exporters for particular configurations.
// We maintain this map because the Factory is asked trace and metric receivers separately
// when it gets CreateTracesReceiver() and CreateMetricsReceiver() but they must not
// create separate objects, they must use one Receiver object per configuration.
var exporters = sharedcomponent.NewSharedComponents()
