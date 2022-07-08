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

/*
 * Main entrance for consuming incoming OTEL data flow and converting it to metadata format
 */

package metadataexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/metadataexporter"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type metadataExporter struct {
	destinations []string
	endpoint     string
	accountName  string
	client       IPurviewClient
}

func (e *metadataExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *metadataExporter) ConsumeTraces(_ context.Context, td ptrace.Traces) error {
	return errors.New("not implemented")
}

func (e *metadataExporter) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	var MetricMetadataPoints []MetricMetadataPoint
	var metricMetadataList = make([]MetricMetadata, 0)

	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resouceMetric := resourceMetrics.At(i)
		// extract resources for the metric
		resources := extractResource(resouceMetric.Resource())

		// extract metrics metadata
		MetricMetadataPoints = extractMetric(resouceMetric.ScopeMetrics())
		metricMetadata := MetricMetadata{
			Resources:            resources,
			MetricMetadataPoints: MetricMetadataPoints,
			Destinations:         e.destinations,
		}

		metricMetadataList = append(metricMetadataList, metricMetadata)
	}

	purviewEntities := make([]PurviewMetadataEntity, 0)

	for _, metricMetadata := range metricMetadataList {
		for _, metricMetadataPoint := range metricMetadata.MetricMetadataPoints {
			purview := NewPurviewEntity(metricMetadataPoint, metricMetadata.Resources, metricMetadata.Destinations)
			purviewEntities = append(purviewEntities, purview)
		}
	}
	e.client.CreateMetadataEntity(PurviewEntityBulkType{
		Entities: purviewEntities,
	})
	return nil
}

func (e *metadataExporter) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	return errors.New("not implemented")
}

func (e *metadataExporter) Start(context.Context, component.Host) error {
	var err error
	if e.client == nil {
		e.client, err = NewPurviewClient(e.endpoint, e.accountName)
	}
	if err != nil {
		return err
	}
	_, err = e.client.CreateMetadataType()
	return err
}

// Shutdown stops the exporter and is invoked during shutdown.
func (e *metadataExporter) Shutdown(context.Context) error {
	e.client = nil
	return nil
}
